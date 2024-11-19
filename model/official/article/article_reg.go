package article

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/mapping"
	"github.com/coscms/webfront/model/official"
)

const GroupName = `article`

func init() {
	official.AddTagGroup(GroupName, `文章`)
	official.AddClickFlowTarget(GroupName, official.ClickFlowTargetFunc(func(c echo.Context, id interface{}) (func(typ string, isCancel ...bool) error, func() official.ClickFlowTargetInfo, error) {
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
	official.AddCollectionTarget(GroupName, `文章`, official.CollectionTargetDoFunc(func(c echo.Context, id interface{}) (func(isCancel ...bool) error, func() (uint64, string), error) {
		articleM := NewArticle(c)
		err := articleM.Get(nil, `id`, id)
		if err != nil {
			if err == db.ErrNoMoreRows {
				err = c.NewError(code.DataNotFound, `文章不存在`)
			}
			return nil, nil, err
		}
		return func(_ ...bool) error { return nil }, func() (uint64, string) {
			return articleM.Id, articleM.Title
		}, nil
	}), official.CollectionTargetListFunc(func(c echo.Context, rows []*official.CollectionResponse, targetIDs []uint64) ([]*official.CollectionResponse, error) {
		articleM := NewArticle(c)
		_, err := articleM.ListByOffset(nil, func(r db.Result) db.Result {
			return r.Select(`id`, `title`, `image`)
		}, 0, -1, `id`, db.In(targetIDs))
		if err != nil {
			return rows, err
		}
		urlGenerator := func(v *dbschema.OfficialCommonArticle) interface{} {
			return c.Echo().URI(`article.detail`, v.Id)
		}
		extra := func(v *dbschema.OfficialCommonArticle) interface{} {
			return echo.H{
				`thumb`: v.Image,
			}
		}
		return mapping.Slice(articleM.Objects(), rows, `Id`, `TargetId`,
			mapping.M{`Title`, `Title`},
			mapping.M{urlGenerator, `URL`},
			mapping.M{extra, `Extra`},
		), nil
	}))
}
