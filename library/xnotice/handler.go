package xnotice

import (
	"github.com/coscms/webcore/library/httpserver"
	"github.com/webx-top/echo"
	ws "github.com/webx-top/echo/handler/websocket"
)

func RegisterRoute(r echo.RouteRegister, cfg echo.H) {
	sender := GetNSenderFromConfig(cfg)
	ws.New("/notice", MakeWSHandler(sender, DefaultReceiver)).Wrapper(r).SetMetaKV(httpserver.PermPublicKV())
	if cfg.Bool(`enableSSE`) {
		r.Get("/sse", MakeSSEHandler(sender)).SetMetaKV(httpserver.PermPublicKV())
	}
}
