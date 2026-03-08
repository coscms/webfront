package official

import (
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/echo"
)

type ClickFlowTargetInfo struct {
	ID         uint64
	AuthorType string
	AuthorID   uint64
	IsAuthor   func(echo.Context, *dbschema.OfficialCustomer) (bool, error)
}

type ClickFlowTargetFunc func(ctx echo.Context, id interface{}) (after func(typ string, isCancel ...bool) error, infoGetter func() ClickFlowTargetInfo, err error)

func (c ClickFlowTargetFunc) Do(ctx echo.Context, id interface{}) (after func(typ string, isCancel ...bool) error, infoGetter func() ClickFlowTargetInfo, err error) {
	return c(ctx, id)
}

func MakeOperator(isCancel ...bool) string {
	if len(isCancel) > 0 && isCancel[0] {
		return `-`
	}
	return `+`
}

type ClickFlowTarget interface {
	Do(ctx echo.Context, id interface{}) (after func(typ string, isCancel ...bool) error, infoGetter func() ClickFlowTargetInfo, err error)
}

var ClickFlowTargets = map[string]ClickFlowTarget{}

func AddClickFlowTarget(name string, target ClickFlowTarget) {
	ClickFlowTargets[name] = target
}
