package minify

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/coscms/webcore/initialize/backend"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

func Merge(b []byte, fs http.FileSystem) []byte {
	m := d.init()
	s := engine.Bytes2str(b)
	matches := m.relatedCSS.FindAllStringSubmatchIndex(s, -1)
	s = mergeBy(s, matches, fs, `css`)
	matches = m.relatedJS.FindAllStringSubmatchIndex(s, -1)
	s = mergeBy(s, matches, fs, `js`)
	return engine.Str2bytes(s)
}

func mergeBy(s string, matches [][]int, fs http.FileSystem, typ string) string {
	if len(matches) == 0 {
		return s
	}
	var replaced string
	repl := com.ReplaceByMatchedIndex(s, matches, &replaced)
	end := len(matches) - 1
	var newContent string
	var combinedContent string
	buildTime := d.buildTime
	savDir := d.saveDir + echo.FilePathSeparator + buildTime
	if fs != nil {
		com.MkdirAll(savDir, os.ModePerm)
	}
	var files []string
	for k, v := range matches {
		var asset string
		var file string
		com.GetMatchedByIndex(s, v, nil, &asset, &file)
		if len(file) > 0 && fs != nil {
			file = strings.SplitN(file, `?`, 2)[0]
			files = append(files, file)
			if asset == `AssetsURL` {
				file = filepath.Join(echo.Wd(), backend.AssetsDir, file)
			} else {
				file = filepath.Join(echo.Wd(), frontend.AssetsDir, file)
			}
			f, err := fs.Open(file)
			if err != nil {
				log.Errorf(`[minify][merge]%s: %v`, file, err)
			} else {
				b, err := io.ReadAll(f)
				f.Close()
				if err != nil {
					log.Errorf(`[minify][merge]%s: %v`, file, err)
				} else {
					combinedContent += engine.Bytes2str(b) + "\n"
				}
			}
		}
		if k == end {
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
			newFile := com.Md5(strings.Join(files, `,`)) + ext
			savFile := savDir + echo.FilePathSeparator + newFile
			if fs != nil {
				err = os.WriteFile(savFile, engine.Str2bytes(combinedContent), 0664)
				if err != nil {
					log.Errorf(`[minify][merge]%s: %v`, file, err)
				}
			}
			newContent = `{{AssetsURL}}/combined/` + buildTime + `/` + newFile
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
