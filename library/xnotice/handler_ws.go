package xnotice

import (
	"github.com/admpub/websocket"
	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/webx-top/echo"
)

func MakeWSHandler(msgGetter NSender, msgSetter NReceiver) func(c *websocket.Conn, ctx echo.Context) error {
	return func(c *websocket.Conn, ctx echo.Context) error {
		customer := sessdata.Customer(ctx)
		//push(writer)
		close, ch, err := msgGetter(ctx, customer)
		if err != nil {
			return err
		}
		if close != nil {
			defer close()
		}
		if ch == nil {
			return nil
		}
		go func() {
			for {
				select {
				case message, ok := <-ch:
					if !ok || message == nil {
						c.Close()
						return
					}
					if send(c, message) != nil {
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		//echo
		execute := func(conn *websocket.Conn) error {
			for {
				mt, message, err := conn.ReadMessage()
				if err != nil {
					return err
				}

				if msgSetter != nil {
					var reply []byte
					reply, err = msgSetter(ctx, customer, message)
					if err != nil {
						return err
					}
					if reply != nil {
						if err = conn.WriteMessage(mt, reply); err != nil {
							return err
						}
					}
				}
			}
		}
		err = execute(c)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				ctx.Logger().Debug(err.Error())
			} else {
				ctx.Logger().Error(err.Error())
			}
		}
		return nil
	}
}
