package xnotice

import (
	"time"

	"github.com/admpub/sse"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/middleware/sessdata"
)

const FMTDateTime = `20060102150405`
const Day = time.Hour * 24

func EncodeClientID(clientID string, customer *dbschema.OfficialCustomer) string {
	return config.FromFile().Encode256(clientID + `|f:` + com.Md5(customer.Name) + `|` + time.Now().Format(FMTDateTime) + `|` + com.RandomAlphanumeric(6))
}

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
		data := make(chan interface{}, 1)
		defer close(data)
		exec := func(msgChan <-chan *notice.Message) {
			var encodedClientID string
			getEncodedClientID := func(msg *notice.Message) string {
				if len(encodedClientID) == 0 {
					encodedClientID = EncodeClientID(msg.ClientID, customer)
				}
				return encodedClientID
			}
			for {
				select {
				case msg, ok := <-msgChan:
					if !ok || msg == nil {
						return
					}
					select {
					case data <- sse.Event{
						Event: notice.SSEventName,
						Data:  msg,
						Id:    getEncodedClientID(msg),
					}:
					case <-ctx.Done():
						msg.Release()
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}
		if msgChan == nil {
			ch := make(chan *notice.Message, 1)
			ch <- notice.NewMessage().SetMode(`-`).SetType(`clientID`).SetClientID(param.AsString(time.Now().UnixMilli()))
			go exec(ch)
		} else {
			go exec(msgChan)
		}
		ctx.SetRenderer(notice.SSERender)
		err = ctx.SSEvent(notice.SSEventName, data)
		return err
	}
}
