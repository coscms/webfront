package xroleutils

import (
	"github.com/coscms/webfront/library/xrole"
	"github.com/coscms/webfront/middleware/sessdata"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/echo"
)

var (
	CustomerPermTTL    = modelCustomer.CustomerPermTTL
	CustomerPermission = modelCustomer.CustomerPermission
	CustomerRoles      = modelCustomer.CustomerRoles
)

func AllowedByRouteName(ctx echo.Context, name string) bool {
	customer := sessdata.Customer(ctx)
	if customer == nil {
		return false
	}
	permission := CustomerPermission(ctx, customer)
	if permission == nil {
		return false
	}
	return xrole.CheckPermissionByRouteName(ctx, customer, permission, name)
}

func AllowedByRoutePath(ctx echo.Context, routePath string) bool {
	customer := sessdata.Customer(ctx)
	if customer == nil {
		return false
	}
	permission := CustomerPermission(ctx, customer)
	if permission == nil {
		return false
	}
	return xrole.CheckPermissionByRoutePath(ctx, customer, permission, routePath) == nil
}
