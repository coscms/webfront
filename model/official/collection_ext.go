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
	ls CollectionTargetList
	do CollectionTargetDo
}

func (c CollectionTarget) HasList() bool {
	return c.ls != nil
}

func (c CollectionTarget) HasDo() bool {
	return c.do != nil
}

func (c CollectionTarget) List(ctx echo.Context, rows []*CollectionResponse, targetIDs []uint64) ([]*CollectionResponse, error) {
	if c.ls == nil {
		return rows, nil
	}
	return c.ls.List(ctx, rows, targetIDs)
}

func (c CollectionTarget) Do(ctx echo.Context, id interface{}) (after func(isCancel ...bool) error, idGetter func() uint64, err error) {
	if c.do == nil {
		return nil, nil, nil
	}
	return c.do.Do(ctx, id)
}

var CollectionTargets = map[string]CollectionTarget{}

func AddCollectionTarget(name string, targetDo CollectionTargetDo, targetList CollectionTargetList) {
	CollectionTargets[name] = CollectionTarget{
		ls: targetList,
		do: targetDo,
	}
}

type CollectionResponse struct {
	*dbschema.OfficialCommonCollection
	Title string `db:"-" json:"title" xml:"title"`
	URL   string `db:"-" json:"url" xml:"url"`
	Extra echo.H `db:"-" json:"extra,omitempty" xml:"extra,omitempty"`
}
