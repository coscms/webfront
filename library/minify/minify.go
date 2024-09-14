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
	"github.com/webx-top/echo"
)

var d = newMinify()

func newMinify() *myMinify {
	return &myMinify{}
}

type myMinify struct {
	relatedCSS *regexp.Regexp
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
	m.relatedCSS = regexp.MustCompile(`[\s]*<link[\s]+combine[\s]+(?:[^>]+\s)?href=["']\{\{(AssetsURL|AssetsXURL)\}\}([^'"]+)["'][^>]*>[\s]*`)
	m.relatedJS = regexp.MustCompile(`[\s]*<script[\s]+combine[\s]+(?:[^>]+\s)?src=["']\{\{(AssetsURL|AssetsXURL)\}\}([^'"]+)["'][^>]*>[\s]*</script>[\s]*`)
	m.minifyM = minify.New()
	m.minifyM.AddFunc("text/css", css.Minify)
	m.minifyM.AddFunc("application/javascript", js.Minify)
	m.saveDir = filepath.Join(echo.Wd(), backend.AssetsDir, `combined`)
	m.buildTime = config.Version.BuildTime
	if len(m.buildTime) == 0 {
		m.buildTime = `0`
	}
	os.RemoveAll(m.saveDir)
}

func MinifyJS(content string) (string, error) {
	return d.init().minifyM.String(`application/javascript`, content)
}

func MinifyCSS(content string) (string, error) {
	return d.init().minifyM.String(`text/css`, content)
}

func Minify() *minify.M {
	return d.init().minifyM
}
