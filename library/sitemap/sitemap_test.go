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
	Registry.Add(`test`, `Test`, echo.KVxOptX[Sitemap, any](Sitemap{
		Do: func(add Adder) error {
			return add(&smg.SitemapLoc{
				Loc:        "news/2021-01-05/a-news-page",
				LastMod:    &now,
				ChangeFreq: smg.Weekly,
				Priority:   1,
			})
		},
	}))

	err := Generate(`https://www.webx.top`)
	assert.NoError(t, err)
}
