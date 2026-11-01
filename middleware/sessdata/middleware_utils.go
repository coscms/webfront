package sessdata

import (
	"context"
	"strconv"

	"github.com/coscms/webcore/library/nerrors"
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
	cid := strconv.FormatUint(customerID, 10)
	cache.Delete(ctx, PermissionCacheKey+cid)
	cache.Delete(ctx, LeftNavigateCacheKey+cid)
	cache.Delete(ctx, TopNavigateCacheKey+cid)
}

// CheckPerm 检查指定路由的权限
func CheckPerm(ctx echo.Context, route string) error {
	fn, ok := ctx.GetFunc(`CheckPerm`).(func(string) error)
	if ok {
		return fn(route)
	}
	return nerrors.ErrUserNoPerm
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
