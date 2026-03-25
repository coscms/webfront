package xcommon

import (
	"os"
	"testing"
)

func TestParseIconFromCSS(t *testing.T) {
	b, _ := os.ReadFile(`/home/swh/go/src/github.com/admpub/nging/public/assets/backend/css/fonts/fontawesome/all.css`)
	c, _ := os.ReadFile(`/home/swh/go/src/github.com/admpub/nging/public/assets/backend/css/fonts/fontawesome/v4-shims.css`)
	results := ParseIconFromCSS(string(b) + string(c))
	t.Logf(`%+v`, results)
}
