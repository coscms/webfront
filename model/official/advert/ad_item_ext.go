package advert

import (
	"strings"
	"time"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/cache"
	"github.com/webx-top/echo"
)

var _ Adverter = (*ItemAndPosition)(nil)

type ItemAndPosition struct {
	*dbschema.OfficialAdItem
	Rendered   string
	AdPosition *dbschema.OfficialAdPosition `db:"-,relation=id:position_id|gtZero"`
}

func (i *ItemAndPosition) GetWidth() uint {
	return i.AdPosition.Width
}

func (i *ItemAndPosition) GetHeight() uint {
	return i.AdPosition.Height
}

func (i *ItemAndPosition) GetURL() string {
	return i.Url
}

func (i *ItemAndPosition) GetContent() string {
	return i.Content
}

func (i *ItemAndPosition) GetContype() string {
	return i.Contype
}

func (i *ItemAndPosition) GetTitle() string {
	if len(i.Title) == 0 {
		return i.AdPosition.Title
	}
	return i.Title
}

func (i *ItemAndPosition) GetDescription() string {
	if len(i.Description) == 0 {
		return i.AdPosition.Description
	}
	return i.Description
}

type ItemAndRendered struct {
	*dbschema.OfficialAdItem
	Rendered string
}

type CachedAdvert struct {
	List        PositionAdverts
	RefreshedAt time.Time
}

func (c *CachedAdvert) GenHTML() *CachedAdvert {
	c.List.GenHTML()
	return c
}

func DeleteCachedAdvert(ctx echo.Context, idents ...string) {
	identString := strings.Join(idents, `,`)
	for _, lc := range config.FromFile().Language.AllList {
		key := `advert:` + lc + `:` + identString
		cache.Delete(ctx, key)
	}
}

func GetCachedAdvert(ctx echo.Context, idents ...string) (*CachedAdvert, error) {
	res := &CachedAdvert{}
	if len(idents) > 0 {
		key := `advert:` + ctx.Lang().Normalize() + `:` + strings.Join(idents, `,`)
		err := cache.XFunc(ctx, key, res, func() error {
			m := NewAdPosition(ctx)
			var err error
			res.List, err = m.GetAdvertsByIdent(true, idents...)
			if err != nil {
				return err
			}
			res.RefreshedAt = time.Now()
			return err
		}, cache.GenOptions(ctx, 300)...)
		if err != nil {
			return nil, err
		}
	}
	if res.List == nil {
		res.List = PositionAdverts{}
	}
	//echo.Dump(echo.H{`idents`: idents, `result`: res})
	return res, nil
}

// GetAdvertForHTML(ctx,`home1`).Place(`home1`).HTML()
func GetAdvertForHTML(ctx echo.Context, idents ...string) interface{} {
	sz := len(idents)
	if sz < 1 || (sz == 1 && len(idents[0]) == 0) {
		return ItemsResponse{}
	}
	cc, err := GetCachedAdvert(ctx, idents...)
	if err != nil {
		return ItemsResponse{
			{Rendered: err.Error()},
		}
	}
	if sz == 1 {
		for _, item := range cc.List {
			return item
		}
		return ItemsResponse{}
	}
	return cc.List
}
