package official

import (
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/echo"
)

type CollectionTargetDoFunc func(ctx echo.Context, id interface{}) (after func(isCancel ...bool) error, idGetter func() uint64, err error)

func (c CollectionTargetDoFunc) Do(ctx echo.Context, id interface{}) (after func(isCancel ...bool) error, idGetter func() uint64, err error) {
	return c(ctx, id)
}

type CollectionTargetListFunc func(ctx echo.Context, rows []*CollectionResponse, targetIDs []uint64) ([]*CollectionResponse, error)

func (c CollectionTargetListFunc) List(ctx echo.Context, rows []*CollectionResponse, targetIDs []uint64) ([]*CollectionResponse, error) {
	return c(ctx, rows, targetIDs)
}

type CollectionTargetList interface {
	List(ctx echo.Context, rows []*CollectionResponse, targetIDs []uint64) ([]*CollectionResponse, error)
}

type CollectionTargetDo interface {
	Do(ctx echo.Context, id interface{}) (after func(isCancel ...bool) error, idGetter func() uint64, err error)
}

type CollectionTarget struct {
	List CollectionTargetList
	Do   CollectionTargetDo
}

var CollectionTargets = map[string]CollectionTarget{}

func AddCollectionTarget(name string, targetDo CollectionTargetDo, targetList CollectionTargetList) {
	CollectionTargets[name] = CollectionTarget{
		List: targetList,
		Do:   targetDo,
	}
}

type CollectionResponse struct {
	*dbschema.OfficialCommonCollection
	Title string `db:"-" json:"title" xml:"title"`
	URL   string `db:"-" json:"url" xml:"url"`
	Extra echo.H `db:"-" json:"extra,omitempty" xml:"extra,omitempty"`
}
