package minify

import (
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/coscms/webcore/initialize/backend"
	"github.com/coscms/webcore/library/config"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

var d = newMinify()

func newMinify() *myMinify {
	return &myMinify{}
}

type myMinify struct {
	relatedCSS *regexp.Regexp
	importCSS  *regexp.Regexp
	relatedJS  *regexp.Regexp
	minifyM    *minify.M
	once       sync.Once
	saveDir    string
	buildTime  string
}

func (m *myMinify) init() *myMinify {
	m.once.Do(m.doinit)
	return m
}

func (m *myMinify) doinit() {
	m.importCSS = regexp.MustCompile(`url\(["']?([^"')]+)["']?\)`)
	m.relatedCSS = regexp.MustCompile(`[\s]*<link[\s]+combine(?:=["']([^"']+)["'])?[\s]+(?:[^>]+\s)?href=["']\{\{(AssetsURL|AssetsXURL)\}\}([^'"]+)["'][^>]*>[\s]*`)
	m.relatedJS = regexp.MustCompile(`[\s]*<script[\s]+combine(?:=["']([^"']+)["'])?[\s]+(?:[^>]+\s)?src=["']\{\{(AssetsURL|AssetsXURL)\}\}([^'"]+)["'][^>]*>[\s]*</script>[\s]*`)
	m.minifyM = minify.New()
	m.minifyM.AddFunc("text/css", css.Minify)
	m.minifyM.AddFunc("application/javascript", js.Minify)
	m.saveDir = filepath.Join(backend.AssetsDir, `backend`, `combined`)
	m.buildTime = config.Version.BuildTime
	if config.FromFile() == nil || config.FromFile().Extend == nil || !config.FromFile().Extend.GetStore(`minify`).Bool(`disableAutoClear`) {
		os.RemoveAll(m.saveDir)
	}
}

func MinifyJS(content string) (string, error) {
	return d.init().minifyM.String(`application/javascript`, content)
}

func MinifyCSS(content string) (string, error) {
	return d.init().minifyM.String(`text/css`, content)
}

func ReplaceCSSImportURL(content, pageURL, combinedPath string) string {
	return d.init().ReplaceCSSImportURL(content, pageURL, combinedPath)
}

func Minify() *minify.M {
	return d.init().minifyM
}
