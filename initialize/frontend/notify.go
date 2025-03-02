package frontend

import (
	"github.com/admpub/events"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	dbschemaNging "github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/top"
)

var Notify = notice.NewUserNotices(false, nil)

func init() {
	echo.OnCallback(`message.send`, func(e events.Event) error {
		msg := e.Context.Get(`data`).(*dbschema.OfficialCommonMessage)
		fromCustomer := e.Context.Get(`fromCustomer`).(*dbschema.OfficialCustomer)
		fromUser := e.Context.Get(`fromUser`).(*dbschemaNging.NgingUser)
		if err := sendMessageNotify(msg, fromCustomer, fromUser); err != nil {
			log.Error(err)
		}
		return nil
	})
}

func sendMessageNotify(f *dbschema.OfficialCommonMessage, fromCustomer *dbschema.OfficialCustomer, fromUser *dbschemaNging.NgingUser) error {
	ctx := f.Context()
	var sender string
	var isAdmin bool
	if fromCustomer != nil {
		sender = fromCustomer.Name
	} else if fromUser != nil {
		sender = fromUser.Username
		isAdmin = true
	}
	sendMessage := func(receiver string, visitURL string) {
		Notify.Send(
			receiver,
			notice.NewMessageWithValue(
				`message`,
				ctx.T(`收到新消息`),
				echo.H{
					`url`:     visitURL,
					`author`:  sender,
					`isAdmin`: isAdmin,
					`content`: com.IfTrue(len(f.Title) > 0, f.Title, ctx.T(`无标题`)),
				},
			).SetID(f.Id),
		)
	}
	if f.CustomerB > 0 {
		custM := dbschema.NewOfficialCustomer(ctx)
		err := custM.Get(func(r db.Result) db.Result {
			return r.Select(`name`, `id`)
		}, `id`, f.CustomerB)
		if err != nil {
			return err
		}
		if len(custM.Name) > 0 {
			visitURL := top.URLByName(`#frontend#user.message.view`, echo.H{`type`: `inbox`, `id`: f.Id})
			sendMessage(custM.Name, visitURL)
		}
	} else if f.UserB > 0 {
		userM := dbschemaNging.NewNgingUser(ctx)
		err := userM.Get(func(r db.Result) db.Result {
			return r.Select(`username`, `id`)
		}, `id`, f.CustomerB)
		if err != nil {
			return err
		}
		if len(userM.Username) > 0 {
			visitURL := top.URLByName(`#backend#admin.message.view`, echo.H{`type`: `inbox`, `id`: f.Id})
			sendMessage(userM.Username, visitURL)
		}
	}
	return nil
}
