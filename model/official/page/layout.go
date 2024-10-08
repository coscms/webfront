package page

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
)

func NewLayout(ctx echo.Context) *Layout {
	m := &Layout{
		OfficialPageLayout: dbschema.NewOfficialPageLayout(ctx),
	}
	return m
}

type Layout struct {
	*dbschema.OfficialPageLayout
}

func (f *Layout) ListPage(cond *db.Compounds, orderby ...interface{}) ([]*LayoutExt, error) {
	list := []*LayoutExt{}
	_, err := common.NewLister(f, &list, func(r db.Result) db.Result {
		return r.OrderBy(orderby...)
	}, cond.And()).Paging(f.Context())
	return list, err
}
