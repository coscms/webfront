package sitemap

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/admpub/sitemap-generator/smg"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func TestOSRoot(t *testing.T) {
	testDir := `./testdata`
	os.MkdirAll(testDir, os.ModePerm)
	defer os.RemoveAll(testDir)
	rootFS, err := os.OpenRoot(testDir)
	assert.NoError(t, err)
	for i := 0; i < 5; i++ {
		content := strconv.Itoa(i)
		rootFS.WriteFile(`test_`+content, []byte(content), os.ModePerm)
		b, err := rootFS.ReadFile(`test_` + content)
		assert.NoError(t, err)
		assert.Equal(t, content, string(b))
	}
	defer rootFS.Close()
}

func TestGenerate(t *testing.T) {
	eCtx := defaults.NewMockContext()
	now := time.Now().UTC()
	cfg := Sitemap{
		Do: func(ctx echo.Context, sm *smg.Sitemap, langCode string, subDirName string) error {
			sm.Add(&smg.SitemapLoc{
				Loc:        "https://www.coscms.com/news/2021-01-05/a-news-page",
				LastMod:    &now,
				ChangeFreq: smg.Weekly,
				Priority:   1,
				Images: []*smg.SitemapImage{
					{ImageLoc: `https://www.coscms.com/test.jpg`},
					{ImageLoc: `/test.jpg`},
					{ImageLoc: `test.jpg`},
				},
			})
			return sm.Add(&smg.SitemapLoc{
				Loc:        "news/2021-01-05/a-news-page",
				LastMod:    &now,
				ChangeFreq: smg.Weekly,
				Priority:   1,
			})
		},
	}
	Registry = echo.NewKVxData[Sitemap, any]()
	Registry.Add(`test`, `Test`, echo.KVxOptX[Sitemap, any](cfg))

	err := GenerateIndex(eCtx, `https://www.webx.top/`, `zh-CN`, true)
	assert.NoError(t, err)
	err = GenerateSingle(eCtx, `https://www.webx.top`, `zh-CN`, cfg.Do)
	assert.NoError(t, err)
}
