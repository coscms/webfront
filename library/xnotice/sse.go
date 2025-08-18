package xnotice

import (
	"github.com/admpub/sse"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webfront/middleware/sessdata"
)

func MakeSSEHandler(msgGetter NSender) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		customer := sessdata.Customer(ctx)
		_close, msgChan, err := msgGetter(ctx, customer)
		if err != nil {
			return err
		}
		if _close != nil {
			defer _close()
		}
		if msgChan == nil {
			return nil
		}
		data := make(chan interface{})
		var encodedClientID string
		go func() {
			defer close(data)
			for {
				select {
				case msg, ok := <-msgChan:
					if !ok || msg == nil {
						return
					}
					if len(encodedClientID) == 0 {
						encodedClientID = config.FromFile().Encode256(msg.ClientID + `|` + com.RandomAlphanumeric(16))
					}
					data <- sse.Event{
						Event: notice.SSEventName,
						Data:  msg,
						Id:    encodedClientID,
					}
				case <-ctx.Done():
					return
				}
			}
		}()
		ctx.SetRenderer(notice.SSERender)
		err = ctx.SSEvent(notice.SSEventName, data)
		return err
	}
}
