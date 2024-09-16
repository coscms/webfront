package sessdata

import (
	"github.com/webx-top/echo"

	"github.com/coscms/webfront/dbschema"
)

// Customer 前台会员客户信息
func Customer(c echo.Context) *dbschema.OfficialCustomer {
	customer, _ := c.Internal().Get(`customer`).(*dbschema.OfficialCustomer)
	return customer
}

// IsAdmin 是否是后台管理员或与后台管理员关联的客户
func IsAdmin(c echo.Context, onlyBackendAdmin ...bool) bool {
	return AdminUID(c, onlyBackendAdmin...) > 0
}

// IsBackendAdmin 是否是后台管理元
func IsBackendAdmin(c echo.Context) bool {
	if user := User(c); user != nil {
		return user.Id > 0
	}
	return false
}

// AdminUID 后台管理员用户ID
func AdminUID(c echo.Context, onlyBackendAdmin ...bool) uint {
	if user := User(c); user != nil {
		return user.Id
	}
	if len(onlyBackendAdmin) > 0 && onlyBackendAdmin[0] {
		return 0
	}
	if customer := Customer(c); customer != nil {
		return customer.Uid
	}
	return 0
}
