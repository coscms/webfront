package xrole

import (
	"sync"

	"github.com/coscms/webcore/library/perm"
	"github.com/coscms/webfront/initialize/frontend/usernav"
	"github.com/webx-top/echo"
)

var (
	navTreeCached *perm.Map
	navTreeOnce   sync.Once
)

func initNavTreeCached() {
	navTreeCached = perm.NewMap(nil)
	navTreeCached.Import(usernav.LeftNavigate)
	navTreeCached.Import(usernav.TopNavigate)
}

func NavTreeCached() *perm.Map {
	navTreeOnce.Do(initNavTreeCached)
	return navTreeCached
}

func init() {
	echo.OnCallback(`nging.httpserver.run.before`, func(_ echo.Event) error {
		NavTreeCached()
		return nil
	})
}
