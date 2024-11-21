package xsocketio

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/webx-top/echo"
	esi "github.com/webx-top/echo-socket.io"
)

type events struct {
	events       []func(esi.IWrapper)
	onConnect    []func(ctx echo.Context, conn socketio.Conn) error
	onError      []func(ctx echo.Context, conn socketio.Conn, e error)
	onDisconnect []func(ctx echo.Context, conn socketio.Conn, msg string)
}

func (e *events) FireGlobal(wrapper esi.IWrapper) {
	for _, fn := range e.events {
		fn(wrapper)
	}
}

func (e *events) FireConnect(ctx echo.Context, conn socketio.Conn) (err error) {
	for _, fn := range e.onConnect {
		if err = fn(ctx, conn); err != nil {
			break
		}
	}
	return
}

func (e *events) FireError(ctx echo.Context, conn socketio.Conn, err error) {
	for _, fn := range e.onError {
		fn(ctx, conn, err)
	}
}

func (e *events) FireDisconnect(ctx echo.Context, conn socketio.Conn, msg string) {
	for _, fn := range e.onDisconnect {
		fn(ctx, conn, msg)
	}
}

var (
	globalEvents    = &events{}
	namespaceEvents = map[string]*events{}
)

// ----- globalEvents -----

func OnEvent(fns ...func(esi.IWrapper)) {
	globalEvents.events = append(globalEvents.events, fns...)
}

func OnConnect(fns ...func(ctx echo.Context, conn socketio.Conn) error) {
	globalEvents.onConnect = append(globalEvents.onConnect, fns...)
}

func OnError(fns ...func(ctx echo.Context, conn socketio.Conn, e error)) {
	globalEvents.onError = append(globalEvents.onError, fns...)
}

func OnDisconnect(fns ...func(ctx echo.Context, conn socketio.Conn, msg string)) {
	globalEvents.onDisconnect = append(globalEvents.onDisconnect, fns...)
}

// ----- namespaceEvents -----

func OnNSEvent(nsp string, fns ...func(esi.IWrapper)) {
	if _, ok := namespaceEvents[nsp]; !ok {
		namespaceEvents[nsp] = &events{}
	}
	namespaceEvents[nsp].events = append(namespaceEvents[nsp].events, fns...)
}

func OnNSConnect(nsp string, fns ...func(ctx echo.Context, conn socketio.Conn) error) {
	if _, ok := namespaceEvents[nsp]; !ok {
		namespaceEvents[nsp] = &events{}
	}
	namespaceEvents[nsp].onConnect = append(namespaceEvents[nsp].onConnect, fns...)
}

func OnNSError(nsp string, fns ...func(ctx echo.Context, conn socketio.Conn, e error)) {
	if _, ok := namespaceEvents[nsp]; !ok {
		namespaceEvents[nsp] = &events{}
	}
	namespaceEvents[nsp].onError = append(namespaceEvents[nsp].onError, fns...)
}

func OnNSDisconnect(nsp string, fns ...func(ctx echo.Context, conn socketio.Conn, msg string)) {
	if _, ok := namespaceEvents[nsp]; !ok {
		namespaceEvents[nsp] = &events{}
	}
	namespaceEvents[nsp].onDisconnect = append(namespaceEvents[nsp].onDisconnect, fns...)
}
