package xsocketio

import (
	"net/http"
	"strings"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/common"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/webx-top/echo"
	esi "github.com/webx-top/echo-socket.io"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware"
)

func RegisterRoute(e echo.RouteRegister, cfg *Config, s ...func(*middleware.CORSConfig)) {
	cfgCORS := &middleware.CORSConfig{
		AllowOrigins: []string{`Token`},
	}
	for _, f := range s {
		f(cfgCORS)
	}
	prefix := e.Prefix()
	nsp := strings.Trim(prefix, `/`)
	nsp = strings.ReplaceAll(nsp, `/`, `_`)
	socket := SocketIO(nsp, cfg)
	e.Any(`/socket.io/`, func(ctx echo.Context) error {
		if common.Setting(`socketio`).String(`enabled`) != `1` {
			return echo.ErrNotFound
		}
		return socket.Handle(ctx)
	}, middleware.CORSWithConfig(*cfgCORS))
}

var RequestChecker engineio.CheckerFunc = func(req *http.Request) (http.Header, error) {
	token := common.Setting(`socketio`).String(`token`)
	if len(token) == 0 {
		return nil, nil
	}
	post := req.Header.Get(`Token`)
	if len(post) == 0 {
		post = req.URL.Query().Get(`token`)
	}
	if token != post {
		if log.IsEnabled(log.LevelDebug) {
			log.Debugf(`[socketIO] invalid token: %q`, post)
			log.Debugf(`[socketIO] request headers: %+v`, req.Header)
		}
		return nil, echo.NewError(`invalid token`, code.InvalidToken)
	}
	return nil, nil
}

func socketIOWrapper(nsp string) *esi.Wrapper {
	wrapper := esi.NewWrapper(&engineio.Options{
		RequestChecker: RequestChecker,
	})
	wrapper.OnConnect(nsp, func(ctx echo.Context, conn socketio.Conn) error {
		if ev := namespaceEvents[nsp]; ev != nil {
			if err := ev.FireConnect(ctx, conn); err != nil {
				return err
			}
		}
		if err := globalEvents.FireConnect(ctx, conn); err != nil {
			return err
		}
		return nil
	})

	wrapper.OnError(nsp, func(ctx echo.Context, conn socketio.Conn, e error) {
		log.Error("[socketIO] meet error: ", e)
		if ev := namespaceEvents[nsp]; ev != nil {
			ev.FireError(ctx, conn, e)
		}
		globalEvents.FireError(ctx, conn, e)
		conn.Close()
	})

	wrapper.OnDisconnect(nsp, func(ctx echo.Context, conn socketio.Conn, msg string) {
		log.Debug("[socketIO] closed", msg)
		if ev := namespaceEvents[nsp]; ev != nil {
			ev.FireDisconnect(ctx, conn, msg)
		}
		globalEvents.FireDisconnect(ctx, conn, msg)
		conn.Close()
	})

	if ev := namespaceEvents[nsp]; ev != nil {
		ev.FireGlobal(wrapper)
	}
	globalEvents.FireGlobal(wrapper)

	return wrapper
}
