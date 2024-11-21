package xnotice

import (
	"context"
	"encoding/json"

	"github.com/admpub/once"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/cmder"
	"github.com/coscms/webfront/library/queue/asynq"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/echo/defaults"
)

type CustomerOnline struct {
	SessionID  string `json:"sessID"`
	CustomerID uint64 `json:"custID"`
	Online     bool   `json:"online"`
}

func (c *CustomerOnline) Apply() error {
	eCtx := defaults.NewMockContext()
	onlineM := modelCustomer.NewOnline(eCtx)
	onlineM.SessionId = c.SessionID
	onlineM.CustomerId = c.CustomerID
	var err error
	if c.Online {
		err = onlineM.Incr(1)
	} else {
		err = onlineM.Decr(1)
	}
	return err
}

var qonce once.Once

func initialize() {
	asynq.WorkerHandleFunc(`customer:online`, func(ctx context.Context, t *asynq.Task) error {
		data := CustomerOnline{}
		err := json.Unmarshal(t.Payload(), &data)
		if err != nil {
			return err
		}
		return data.Apply()
	})
	asynq.StartWorker()
}

type queueCmder struct {
	*cmder.Simple
}

func (q *queueCmder) Boot() error {
	config.FromCLI().ParseConfig()
	initialize()
	return nil
}

// ./webx --config ./config/config.yaml --type queue:worker
func RegsterCmder() {
	cmder.Register(`queue:worker`, &queueCmder{Simple: cmder.NewSimple()})
}

func SendOnlineStatusToQueue(sessionId string, customerID uint64, online bool) error {
	qonce.Do(initialize)
	_, err := asynq.SendJSON(`customer:online`, CustomerOnline{
		SessionID:  sessionId,
		CustomerID: customerID,
		Online:     online,
	})
	return err
}
