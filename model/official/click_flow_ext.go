package official

import "github.com/webx-top/echo"

type ClickFlowTargetInfo struct {
	ID       uint64
	AuthorID uint64
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
