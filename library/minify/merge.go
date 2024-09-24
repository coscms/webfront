package minify

import (
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

func Merge(b []byte, fileNoop ...bool) []byte {
	cdnCfg := config.Setting(`base`, `assetsCDN`)
	backendCDN := cdnCfg.String(`backend`)
	frontendCDN := cdnCfg.String(`frontend`)
	hasBackendCDN := len(backendCDN) > 0
	hasFrontendCDN := len(frontendCDN) > 0
	if hasBackendCDN && hasFrontendCDN {
		return b
	}
	m := d.init()
	s := engine.Bytes2str(b)
	var fnop bool
	if len(fileNoop) > 0 {
		fnop = fileNoop[0]
	}
	s = m.mergeBy(s, `css`, fnop, hasBackendCDN, hasFrontendCDN)
	s = m.mergeBy(s, `js`, fnop, hasBackendCDN, hasFrontendCDN)
	return engine.Str2bytes(s)
}

func (m *myMinify) mergeBy(s string, typ string, fileNoop bool, hasBackendCDN bool, hasFrontendCDN bool) string {
	var matches [][]int
	if typ == `css` {
		matches = m.relatedCSS.FindAllStringSubmatchIndex(s, -1)
	} else {
		matches = m.relatedJS.FindAllStringSubmatchIndex(s, -1)
	}
	if len(matches) == 0 {
		return s
	}
	var replaced string
	repl := com.ReplaceByMatchedIndex(s, matches, &replaced)
	end := len(matches) - 1
	var newContent string
	var combinedContent string
	buildTime := m.buildTime
	savDir := m.saveDir
	if !fileNoop {
		com.MkdirAll(savDir, os.ModePerm)
	}
	combinedPath := path.Join(httpserver.Backend.AssetsURLPath, `combined`)
	var groups []string
	files := map[string][]string{}
	eqNextGroup := func(k int, group string) bool {
		if k >= end {
			return false
		}
		var nextGroup string
		com.GetMatchedByIndex(s, matches[k+1], nil, &nextGroup)
		return group == nextGroup
	}
	for k, v := range matches {
		var group string
		var asset string
		var file string
		com.GetMatchedByIndex(s, v, nil, &group, &asset, &file)
		if hasBackendCDN {
			if asset == `AssetsURL` {
				repl(k, v)
				continue
			}
		} else if hasFrontendCDN {
			if asset != `AssetsURL` {
				repl(k, v)
				continue
			}
		}
		if _, ok := files[group]; !ok {
			files[group] = []string{}
			groups = append(groups, group)
		}
		if len(file) > 0 {
			file = strings.SplitN(file, `?`, 2)[0]
			files[group] = append(files[group], file)
			if !fileNoop {
				f, err := openfile(asset, file)
				if err != nil {
					log.Errorf(`[minify][merge]%s: %v`, file, err)
				} else {
					b, err := io.ReadAll(f)
					f.Close()
					if err != nil {
						log.Errorf(`[minify][merge]%s: %v`, file, err)
					} else {
						content := engine.Bytes2str(b)
						content = strings.TrimSpace(content)
						if typ == `css` {
							var pageURL string
							if asset == `AssetsURL` {
								pageURL = path.Join(httpserver.Backend.AssetsURLPath, file)
							} else {
								pageURL = path.Join(httpserver.Frontend.AssetsURLPath, file)
							}
							content = d.ReplaceCSSImportURL(content, pageURL, combinedPath)
						} else {
							if !strings.HasSuffix(content, `;`) {
								content += `;`
							}
						}
						combinedContent += content + "\n"
					}
				}
			}
		}
		newContent = ``
		if k == end || !eqNextGroup(k, group) {
			var err error
			var ext string
			switch typ {
			case `js`:
				combinedContent, err = d.minifyM.String(`application/javascript`, combinedContent)
				ext = `.min` + `.` + typ
			case `css`:
				combinedContent, err = d.minifyM.String(`text/css`, combinedContent)
				ext = `.min` + `.` + typ
			default:
				ext = `.` + typ
			}
			if err != nil {
				log.Errorf(`[minify][merge]%s: %v`, file, err)
			}
			//com.Dump(map[string]interface{}{`files`: files[group], `group`: group})
			newFile := com.Md5(strings.Join(files[group], `,`)) + ext
			if len(group) > 0 && com.StrIsAlphaNumeric(group) {
				newFile = group + `-` + newFile
			} else {
				newFile = strconv.Itoa(len(groups)-1) + `-` + newFile
			}
			savFile := savDir + echo.FilePathSeparator + newFile
			if !fileNoop {
				err = os.WriteFile(savFile, engine.Str2bytes(combinedContent), 0664)
				if err != nil {
					log.Errorf(`[minify][merge]%s: %v`, file, err)
				}
			}
			newContent = `{{AssetsURL}}/combined/` + newFile
			if len(buildTime) > 0 {
				newContent += `?t=` + buildTime
			}
			switch typ {
			case `js`:
				newContent = `<script src="` + newContent + `"></script>`
			case `css`:
				newContent = `<link href="` + newContent + `" rel="stylesheet" />`
			}
		}
		repl(k, v, newContent)
	}
	return replaced
}

func (m *myMinify) ReplaceCSSImportURL(content, pageURL, combinedPath string) string {
	content = d.importCSS.ReplaceAllStringFunc(content, func(s string) string {
		return replaceCSSImportURL(s, pageURL, combinedPath)
	})
	return content
}

func resolveURLPath(u string, targetPath string) string {
	if len(targetPath) == 0 {
		return u
	}
	if strings.HasPrefix(u, targetPath) {
		u = strings.TrimPrefix(u, targetPath)
		u = strings.TrimPrefix(u, `/`)
		return u
	}
	var prefix string
	tp := path.Dir(targetPath)
	for len(tp) > 0 {
		prefix += `../`
		if strings.HasPrefix(u, tp) {
			u = strings.TrimPrefix(u, tp)
			return prefix + strings.TrimPrefix(u, `/`)
		}
		tp = path.Dir(tp)
	}
	return u
}

func absURLPath(s string, pageURL string) string {
	if len(s) == 0 || strings.HasPrefix(s, `/`) || strings.Contains(s, `://`) {
		return s
	}
	for strings.HasPrefix(s, `./`) {
		s = strings.TrimPrefix(s, `./`)
	}
	urlPath := path.Dir(pageURL)
	for strings.HasPrefix(s, `../`) {
		urlPath = path.Dir(urlPath)
		s = strings.TrimPrefix(s, `../`)
	}
	s = path.Join(urlPath, s)
	return s
}

func replaceCSSImportURL(s string, pageURL string, combinedPath string) string {
	s = strings.TrimPrefix(s, `url(`)
	s = strings.TrimSuffix(s, `)`)
	s = strings.Trim(s, `"'`)
	if strings.HasPrefix(s, `data:`) { // data:application/x-font-woff2;charset=utf-8;base64,
		return s
	}
	s = absURLPath(s, pageURL)
	s = resolveURLPath(s, combinedPath)
	return `url(` + s + `)`
}
