package xrole

import (
	"strings"

	"github.com/coscms/webcore/library/perm"
	"github.com/coscms/webfront/initialize/frontend/usernav"
	"github.com/webx-top/echo"
)

func PermPageGenerator(ctx echo.Context) (string, error) {
	return perm.BuildPermActions(ctx.FormValues(`permAction[]`)), nil
}

func PermPageChecker(ctx echo.Context, parsed interface{}, current string) (interface{}, error) {
	current = strings.TrimPrefix(current, `/`)
	if len(current) == 0 {
		return NavTreeCached().Check(current), nil
	}
	//echo.Dump(parsed)
	permPages, ok := parsed.(*perm.Map)
	if !ok {
		return false, nil
	}
	return permPages.Check(current), nil
}

func PermPageParser(ctx echo.Context, rule string) (interface{}, error) {
	navTree := NavTreeCached()
	permPages := perm.NewMap(navTree)
	permPages.Parse(rule)
	return permPages, nil
}

var PermPageList = func(ctx echo.Context) ([]interface{}, error) {
	return nil, nil
}

func PermPageOnRender(ctx echo.Context) error {
	ctx.Set(`customerLeftNavigate`, usernav.LeftNavigate)
	ctx.Set(`customerTopNavigate`, usernav.TopNavigate)
	return nil
}

func PermBehaviorGenerator(ctx echo.Context) (string, error) {
	values := map[string][]string{}
	for _, permName := range ctx.FormValues(`permBehavior[]`) {
		values[permName] = ctx.FormValues(`permBehaviorConfig[` + permName + `]`)
	}
	if len(values) == 0 {
		return ``, nil
	}
	return perm.SerializeBehaviorValues(values, Behaviors)
}

func PermBehaviorChecker(ctx echo.Context, parsed interface{}, current string) (interface{}, error) {
	permBehaviors, ok := parsed.(perm.BehaviorPerms)
	if !ok {
		return &perm.CheckedBehavior{}, nil
	}
	return permBehaviors.CheckBehavior(current), nil
}

func PermBehaviorParser(ctx echo.Context, rule string) (interface{}, error) {
	return perm.ParseBehavior(rule, Behaviors)
}

func PermBehaviorList(ctx echo.Context) ([]interface{}, error) {
	behaviors := Behaviors.Slice()
	list := make([]interface{}, len(behaviors))
	for k, v := range behaviors {
		list[k] = v
	}
	return list, nil
}

func PermBehaviorOnRender(ctx echo.Context) error {
	behaviorList, err := PermBehaviorList(ctx)
	if err != nil {
		return err
	}
	ctx.Set(`customerBehaviorList`, behaviorList)
	return nil
}

func PermBehaviorIsValid(ctx echo.Context) bool {
	if list, ok := ctx.Get(`customerBehaviorList`).([]interface{}); ok {
		return len(list) > 0
	}
	return false
}
