package sitemap

import (
	"testing"
	"time"

	"github.com/admpub/sitemap-generator/smg"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func TestGenerate(t *testing.T) {
	eCtx := defaults.NewMockContext()
	now := time.Now().UTC()
	cfg := Sitemap{
		Do: func(ctx echo.Context, sm *smg.Sitemap) error {
			return sm.Add(&smg.SitemapLoc{
				Loc:        "news/2021-01-05/a-news-page",
				LastMod:    &now,
				ChangeFreq: smg.Weekly,
				Priority:   1,
			})
		},
	}
	Registry.Add(`test`, `Test`, echo.KVxOptX[Sitemap, any](cfg))

	err := GenerateIndex(eCtx, `https://www.webx.top/`, `zh-CN`, false)
	assert.NoError(t, err)
	err = GenerateSingle(eCtx, `https://www.webx.top`, `zh-CN`, cfg.Do)
	assert.NoError(t, err)
}
