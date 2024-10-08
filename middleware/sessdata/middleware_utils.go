package sessdata

import (
	"context"
	"fmt"

	"github.com/coscms/webfront/library/cache"
	"github.com/webx-top/echo"
)

var (
	// PermissionCacheKey 权限缓存Key前缀
	PermissionCacheKey = `customer:permission:`
	// LeftNavigateCacheKey 左边栏菜单缓存Key前缀
	LeftNavigateCacheKey = `customer:navigate:left:`
	// TopNavigateCacheKey 顶部菜单缓存Key前缀
	TopNavigateCacheKey = `customer:navigate:top:`
)

// ClearPermissionCache 删除用户的权限缓存
func ClearPermissionCache(ctx context.Context, customerID uint64) {
	cid := fmt.Sprint(customerID)
	cache.Delete(ctx, PermissionCacheKey+cid)
	cache.Delete(ctx, PermissionCacheKey+cid)
	cache.Delete(ctx, PermissionCacheKey+cid)
}

// CheckPerm 检查指定路由的权限
func CheckPerm(ctx echo.Context, route string) error {
	return ctx.GetFunc(`CheckPerm`).(func(string) error)(route)
}

func init() {
	echo.OnCallback(`webx.customer.role.change`, func(v echo.Event) error {
		customerID := v.Context.Uint64(`customerID`)
		if customerID == 0 {
			return nil
		}
		ClearPermissionCache(context.Background(), customerID)
		return nil
	})
}
