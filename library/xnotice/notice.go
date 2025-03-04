package xnotice

import (
	"context"
	"encoding/json"

	"github.com/admpub/log"
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	ws "github.com/webx-top/echo/handler/websocket"

	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/coscms/webfront/middleware/sessdata"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

type NSender func(ctx context.Context, customer *dbschema.OfficialCustomer) (close func(), msg <-chan *notice.Message, err error)
type NReceiver func(ctx context.Context, customer *dbschema.OfficialCustomer, message []byte) ([]byte, error)

var (
	DefaultSender   = MemoryNoticeSender
	DefaultReceiver NReceiver
)

func send(c *websocket.Conn, message *notice.Message) {
	defer message.Release()
	msgBytes, err := json.Marshal(message)
	if err != nil {
		message.Failure()
		log.Error(`Push error (json.Marshal): `, err.Error())
		c.Close()
		return
	}
	log.Debugf(`Push message: %s`, msgBytes)
	err = c.WriteMessage(websocket.TextMessage, msgBytes)
	if err == nil {
		message.Success()
	}
}

func MakeHandler(msgGetter NSender, msgSetter NReceiver) func(c *websocket.Conn, ctx echo.Context) error {
	return func(c *websocket.Conn, ctx echo.Context) error {
		customer := sessdata.Customer(ctx)
		//push(writer)
		if msgGetter != nil {
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
						send(c, message)
					case <-ctx.Done():
						return
					}
				}
			}()
		}

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
		err := execute(c)
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

func ResetClientCount() {
	ctx := defaults.NewMockContext()
	m := modelCustomer.NewOnline(ctx)
	m.ResetClientCount(true)
}

func MemoryNoticeSender(ctx context.Context, customer *dbschema.OfficialCustomer) (func(), <-chan *notice.Message, error) {
	return frontend.Notify.MakeMessageGetter(customer.Name, `message`)
}

func OnlineStatusDBUpdater(ctx context.Context, customer *dbschema.OfficialCustomer) (func(), <-chan *notice.Message, error) {
	c := ctx.(echo.Context)
	sessionID := c.Session().ID()
	if len(sessionID) > 0 || customer != nil {
		onlineM := modelCustomer.NewOnline(c)
		onlineM.SessionId = sessionID
		if customer != nil {
			onlineM.CustomerId = customer.Id
		}
		err := onlineM.Incr(1)
		if err != nil {
			return nil, nil, err
		}
		return func() {
			onlineM.Decr(1)
		}, nil, nil
	}
	return nil, nil, nil
}

func OnlineStatusQueueUpdater(ctx context.Context, customer *dbschema.OfficialCustomer) (func(), <-chan *notice.Message, error) {
	c := ctx.(echo.Context)
	sessionID := c.Session().ID()
	if len(sessionID) > 0 || customer != nil {
		err := SendOnlineStatusToQueue(sessionID, customer.Id, true)
		if err != nil {
			return nil, nil, err
		}
		return func() {
			SendOnlineStatusToQueue(sessionID, customer.Id, false)
		}, nil, nil
	}
	return nil, nil, nil
}

func GetNSenderFromConfig(cfg echo.H) NSender {
	var noticeNS NSender
	switch cfg.String(`store`) {
	case `database`, `db`:
		noticeNS = OnlineStatusDBUpdater
	case `queue`:
		noticeNS = OnlineStatusQueueUpdater
	case `memory`:
		noticeNS = MemoryNoticeSender
	default:
		noticeNS = DefaultSender
	}
	return noticeNS
}

func RegisterRoute(r echo.RouteRegister, cfg echo.H) echo.IRouter {
	return ws.New("/notice", MakeHandler(GetNSenderFromConfig(cfg), DefaultReceiver)).Wrapper(r)
}
