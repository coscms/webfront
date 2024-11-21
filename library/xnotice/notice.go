package xnotice

import (
	"context"
	"encoding/json"

	"github.com/admpub/log"
	"github.com/admpub/websocket"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/coscms/webfront/middleware/sessdata"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	ws "github.com/webx-top/echo/handler/websocket"
)

type NSender func(context.Context, *dbschema.OfficialCustomer) (close func(), msg <-chan *notice.Message, err error)
type NReceiver func(context.Context, *dbschema.OfficialCustomer, []byte) ([]byte, error)

var (
	DefaultSender   = MemoryNoticeSender
	DefaultReceiver NReceiver
)

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
			} else if ch == nil {
				return nil
			}
			if ch != nil {
				go func() {
					for {
						select {
						case message, ok := <-ch:
							if !ok || message == nil {
								c.Close()
								return
							}
							msgBytes, err := json.Marshal(message)
							message.Release()
							if err != nil {
								log.Error(`Push error (json.Marshal): `, err.Error())
								c.Close()
								return
							}
							log.Debugf(`Push message: %s`, msgBytes)
							c.WriteMessage(websocket.TextMessage, msgBytes)
						case <-ctx.Done():
							return
						}
					}
				}()
			}
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
		defer onlineM.Decr(1)
	}
	return frontend.Notify.MakeMessageGetter(customer.Name)
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

func RegisterRoute(r echo.RouteRegister) {
	cfg := config.FromFile().Extend.GetStore(`frontendNotice`)
	if cfg.Bool(`disabled`) {
		return
	}
	ws.New("/notice", MakeHandler(GetNSenderFromConfig(cfg), DefaultReceiver)).Wrapper(r)
}
