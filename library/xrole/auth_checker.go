package xrole

import (
	"strings"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/echo"
)

type (
	AuthChecker func(
		c echo.Context,
		rpath string,
		customer *dbschema.OfficialCustomer,
		permission *RolePermission,
	) (err error, ppath string, returning bool)
	AuthCheckers map[string]AuthChecker
)

func (a AuthCheckers) Check(
	c echo.Context,
	rpath string,
	customer *dbschema.OfficialCustomer,
	permission *RolePermission,
) (err error, ppath string, returning bool) {
	if checker, ok := a[rpath]; ok {
		return checker(c, rpath, customer, permission)
	}
	ppath = rpath
	return
}

var SpecialAuths = AuthCheckers{
	`/user/file/crop`: func(
		c echo.Context,
		rpath string,
		customer *dbschema.OfficialCustomer,
		permission *RolePermission,
	) (err error, ppath string, returning bool) {
		ppath = `/user/file/upload/:type`
		return
	},
}

func AuthRegister(ppath string, checker AuthChecker) {
	SpecialAuths[ppath] = checker
}

func AuthUnregister(ppath string) {
	if _, ok := SpecialAuths[ppath]; ok {
		delete(SpecialAuths, ppath)
	}
}

func CheckPermissionByRoutePath(ctx echo.Context, customer *dbschema.OfficialCustomer, permission *RolePermission, routePath string) error {
	err, route, ret := SpecialAuths.Check(ctx, routePath, customer, permission)
	if ret || err != nil {
		return err
	}
	if route == `/user/index` {
		return nil
	}
	route = strings.TrimPrefix(route, `/user/`)
	if !permission.Check(ctx, route) {
		return nerrors.ErrUserNoPerm
	}
	return nil
}

func CheckPermissionByRouteName(ctx echo.Context, customer *dbschema.OfficialCustomer, permission *RolePermission, name string) bool {
	route := ctx.Echo().GetRoutePathByName(name)
	if len(route) == 0 {
		log.Warnf(`the route named %s could not be found`, name)
		return false
	}
	return CheckPermissionByRoutePath(ctx, customer, permission, route) == nil
}
