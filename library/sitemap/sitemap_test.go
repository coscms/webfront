package sitemap

import (
	"testing"
	"time"

	"github.com/admpub/sitemap-generator/smg"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

func TestGenerate(t *testing.T) {
	now := time.Now().UTC()
	cfg := Sitemap{
		Do: func(add Adder) error {
			return add(&smg.SitemapLoc{
				Loc:        "news/2021-01-05/a-news-page",
				LastMod:    &now,
				ChangeFreq: smg.Weekly,
				Priority:   1,
			})
		},
	}
	Registry.Add(`test`, `Test`, echo.KVxOptX[Sitemap, any](cfg))

	err := GenerateIndex(`https://www.webx.top`, false)
	assert.NoError(t, err)
	err = GenerateSingle(`https://www.webx.top`, cfg.Do)
	assert.NoError(t, err)
}
