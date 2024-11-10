package article

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webfront/library/mapping"
	"github.com/coscms/webfront/model/official"
)

const GroupName = `article`

func init() {
	official.AddClickFlowTarget(GroupName, official.ClickFlowTargetFunc(func(c echo.Context, id interface{}) (func(typ string, isCancel ...bool) error, func() uint64, error) {
		articleM := NewArticle(c)
		err := articleM.Get(nil, `id`, id)
		if err != nil {
			if err == db.ErrNoMoreRows {
				err = c.NewError(code.DataNotFound, `文章不存在`)
			}
			return nil, nil, err
		}
		return func(typ string, isCancel ...bool) error {
			field := typ + `s`
			return articleM.UpdateField(nil, field, db.Raw(field+official.MakeOperator(isCancel...)+`1`), db.Cond{`id`: id})
		}, nil, nil
	}))
	official.AddCollectionTarget(GroupName, official.CollectionTargetDoFunc(func(c echo.Context, id interface{}) (func(isCancel ...bool) error, func() uint64, error) {
		articleM := NewArticle(c)
		err := articleM.Get(nil, `id`, id)
		if err != nil {
			if err == db.ErrNoMoreRows {
				err = c.NewError(code.DataNotFound, `文章不存在`)
			}
			return nil, nil, err
		}
		return func(_ ...bool) error { return nil }, nil, nil
	}), official.CollectionTargetListFunc(func(c echo.Context, rows []*official.CollectionResponse, targetIDs []uint64) ([]*official.CollectionResponse, error) {
		articleM := NewArticle(c)
		err := articleM.ListByOffset(nil, func(r db.Result) db.Result {
			return r.Select(`id`, `title`)
		}, `id`, db.In(targetIDs))
		if err != nil {
			return rows, err
		}
		detailURL := c.Echo().URI(`article.detail`, `{Id}`)
		return mapping.Slice(articleM.Objects(), rows, `Id`, `TargetId`, map[interface{}]string{
			`Title`:                   `Title`,
			mapping.Layout(detailURL): `URL`,
		}), nil
	}))
}
