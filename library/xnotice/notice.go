package xnotice

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/websocket"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/initialize/frontend"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

type NSender func(ctx echo.Context, customer *dbschema.OfficialCustomer) (close func(), msg <-chan *notice.Message, err error)
type NReceiver func(ctx echo.Context, customer *dbschema.OfficialCustomer, message []byte) ([]byte, error)

var (
	DefaultSender   = MemoryNoticeSender
	DefaultReceiver NReceiver
)

func send(c *websocket.Conn, message *notice.Message) error {
	defer message.Release()
	msgBytes, err := json.Marshal(message)
	if err != nil {
		message.Failure()
		log.Error(`Push error (json.Marshal): `, err.Error())
		c.Close()
		return err
	}
	log.Debugf(`Push message: %s`, msgBytes)
	err = c.WriteMessage(websocket.TextMessage, msgBytes)
	if err == nil {
		message.Success()
	}
	return err
}

func ResetClientCount() {
	ctx := defaults.NewMockContext()
	m := modelCustomer.NewOnline(ctx)
	m.ResetClientCount(true)
}

func MemoryNoticeSender(ctx echo.Context, customer *dbschema.OfficialCustomer) (func(), <-chan *notice.Message, error) {
	if ctx.Internal().Bool(`notice.clearClient`) {
		frontend.Notify.CloseAllClient(customer.Name)
	}
	clientID := ctx.Internal().String(`notice.clientID`)
	if len(clientID) == 0 {
		if lastEventID := ctx.Header(`Last-Event-Id`); len(lastEventID) > 0 {
			plaintext := config.FromFile().Decode256(lastEventID)
			if len(plaintext) > 0 {
				parts := strings.SplitN(plaintext, `|`, 4)
				if len(parts) == 4 && parts[1] == `f:`+com.Md5(customer.Name) {
					t, err := time.Parse(FMTDateTime, parts[2])
					if err == nil && !t.IsZero() && time.Since(t) < Day {
						clientID = parts[0]
					}
				}
			}
		}
	}
	if len(clientID) > 0 {
		_close, msgChan := frontend.Notify.MakeMessageGetterWithClientID(customer.Name, clientID, `message`)
		return _close, msgChan, nil
	}
	return frontend.Notify.MakeMessageGetter(customer.Name, `message`)
}

func OnlineStatusDBUpdater(c echo.Context, customer *dbschema.OfficialCustomer) (func(), <-chan *notice.Message, error) {
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

func OnlineStatusQueueUpdater(c echo.Context, customer *dbschema.OfficialCustomer) (func(), <-chan *notice.Message, error) {
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
