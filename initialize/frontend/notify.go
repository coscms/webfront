package frontend

import (
	"strings"

	"github.com/admpub/events"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"

	dbschemaNging "github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/config"
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
	var senderAvatar string
	var senderGender string
	var isAdmin bool
	if fromCustomer != nil {
		sender = fromCustomer.Name
		senderAvatar = fromCustomer.Avatar
		senderGender = fromCustomer.Gender
	} else if fromUser != nil {
		sender = fromUser.Username
		senderAvatar = fromUser.Avatar
		senderGender = fromUser.Gender
		isAdmin = true
	}
	notifyAudioCfg := config.FromFile().Extend.GetStore(`notifyAudio`)
	var notifyAudio string
	disabled := notifyAudioCfg.Bool(`disabled`)
	if !disabled {
		notifyAudio = notifyAudioCfg.String(`audio`)
		if len(notifyAudio) == 0 {
			notifyAudio = subdomains.Default.URL(`/public/assets/backend/audio/notify-dingdong.mp3`, `backend`)
		} else if !strings.Contains(notifyAudio, `/`) {
			notifyAudio = subdomains.Default.URL(`/public/assets/backend/audio/`+notifyAudio, `backend`)
		}
	}
	sendMessage := func(receiver, visitURL string) {
		Notify.Send(
			receiver,
			notice.NewMessageWithValue(
				`message`,
				ctx.T(`收到新消息`),
				echo.H{
					`url`:     visitURL,
					`author`:  sender,
					`avatar`:  senderAvatar,
					`gender`:  senderGender,
					`isAdmin`: isAdmin,
					`content`: com.IfTrue(len(f.Title) > 0, f.Title, ctx.T(`无标题`)),
					`sound`:   notifyAudio,
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
