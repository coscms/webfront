package frontend

import (
	"strings"
	"time"

	"github.com/admpub/events"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/echo/subdomains"

	"github.com/coscms/sms"
	dbschemaNging "github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/cron"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/top"
	"github.com/coscms/webfront/library/xcommon"
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

var noticeDefaultCallback = &notice.Callback{
	Failure: onSendMessageNotifyFail,
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
			).SetID(f.Id).SetCallback(noticeDefaultCallback),
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
		}, `id`, f.UserB)
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

func onSendMessageNotifyFail(m *notice.Message) {
	msgID := param.AsUint64(m.ID)
	data := m.Content.(echo.H)
	ctx := defaults.AcquireMockContext()
	defer defaults.ReleaseMockContext(ctx)
	msgM := dbschema.NewOfficialCommonMessage(ctx)
	err := msgM.Get(func(r db.Result) db.Result {
		return r.Select(`customer_b`, `user_b`)
	}, `id`, msgID)
	if err != nil {
		return
	}
	var email, mobile, username, siteURL string
	if msgM.CustomerB > 0 {
		custM := dbschema.NewOfficialCustomer(ctx)
		err = custM.Get(func(r db.Result) db.Result {
			return r.Select(`name`, `mobile`, `mobile_bind`, `email`, `email_bind`, `id`)
		}, `id`, msgM.CustomerB)
		if err != nil {
			return
		}
		mobile = custM.Mobile
		email = custM.Email
		username = custM.Name
		siteURL = xcommon.FrontendURL(ctx)
	} else if msgM.UserB > 0 {
		userM := dbschemaNging.NewNgingUser(ctx)
		err = userM.Get(func(r db.Result) db.Result {
			return r.Select(`username`, `mobile`, `email`, `id`)
		}, `id`, msgM.UserB)
		if err != nil {
			return
		}
		mobile = userM.Mobile
		email = userM.Email
		username = userM.Username
		siteURL = common.BackendURL(ctx)
	}
	visitURL := data.String(`url`)
	if strings.HasPrefix(visitURL, `/`) {
		visitURL = siteURL + visitURL
	}
	baseCfg := config.Setting().GetStore(`base`)
	if len(mobile) > 0 {
		notifySMSCfg := config.FromFile().Extend.GetStore(`notifySMS`)
		if notifySMSCfg.Bool(`on`) {
			//发送短信
			provider, smsProviderName := sms.AnyOne()
			if provider == nil || len(smsProviderName) == 0 {
				err = ctx.NewError(code.DataUnavailable, `找不到短信发送服务`).SetZone(`provider`)
				log.Error(err)
				return
			}
			message := notifySMSCfg.String(`messageTemplateContent`)
			if len(message) == 0 {
				message = ctx.T(`亲爱的客户: %s，用户「%s」给你发送了站内信，请进入网站查看 %s [%s]`, username, data.String(`author`), visitURL, baseCfg.String(`siteName`))
			} else {
				placeholders := map[string]string{
					`sender`:   data.String(`author`),
					`receiver`: username,
					`name`:     username,
					`url`:      visitURL,
					`siteName`: baseCfg.String(`siteName`),
				}
				for find, to := range placeholders {
					message = strings.ReplaceAll(message, `{`+find+`}`, to)
				}
			}
			smsConfig := sms.NewConfig()
			smsConfig.Mobile = mobile
			smsConfig.Content = message
			//smsConfig.Template = ``
			//smsConfig.SignName = ``
			//smsConfig.ExtraData =
			err = provider.Send(smsConfig)
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
	if len(email) > 0 {
		link := `<a href="` + visitURL + `" target="_blank">` + data.String(`content`) + `</a>`
		content := ctx.T(`亲爱的客户: %s，用户「%s」给你发送了站内信，请进入网站查看 %s<br /><br /> 来自：%s<br />时间：%s`, username, data.String(`author`), link, siteURL+`/`, time.Now().Format(time.RFC3339))
		err = cron.SendMail(email, username, m.Title, com.Str2bytes(content))
		if err != nil {
			log.Error(err)
			return
		}
	}
}
