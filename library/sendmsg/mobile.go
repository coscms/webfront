package sendmsg

import (
	"fmt"
	"strings"
	"time"

	"github.com/coscms/sms"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webcore/model"
	uploadChecker "github.com/coscms/webcore/registry/upload/checker"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

// MobileSend 发送验证码短信
func MobileSend(ctx echo.Context, m *modelCustomer.Customer, purpose string, messages ...string) error {
	var err error
	vm := model.NewCode(ctx)
	now := time.Now()

	if err := vm.CheckFrequency(
		m.Id,
		`customer`,
		`mobile`,
		config.Setting().GetStore(`frequency`).GetStore(`mobile`),
	); err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
	}
	vm.Verification.Reset()
	data := captchabiz.VerifyCaptcha(ctx, httpserver.KindFrontend, `code`)
	if nerrors.IsFailureCode(data.GetCode()) {
		return nil
	}
	if m.MobileBind != `Y` {
		m.Mobile = ctx.Formx(`mobile`).String()
		if err := ctx.Validate(`mobile`, m.Mobile, `mobile`); err != nil {
			return ctx.NewError(code.InvalidParameter, `手机号码格式不正确`).SetZone(`mobile`)
		}
		customerM := modelCustomer.NewCustomer(ctx)
		exists, err := customerM.ExistsOther(m.Mobile, m.Id, `mobile`)
		if err != nil {
			return err
		}
		if exists {
			return ctx.NewError(code.DataAlreadyExists, `手机号码“%s”已被其他账号绑定`, m.Mobile).SetZone(`mobile`)
		}
	}

	verifyCode := com.RandomNumeric(VerifyCodeLength())
	//发送短信
	provider, smsProviderName := sms.AnyOne()
	if provider == nil || len(smsProviderName) == 0 {
		err = ctx.NewError(code.DataUnavailable, `找不到短信发送服务`).SetZone(`provider`)
		return err
	}
	smsConfig := sms.NewConfig()
	smsConfig.Mobile = m.Mobile

	//获取系统配置
	baseCfg := config.Setting().GetStore(`base`)

	//验证码有效期
	lifetime := baseCfg.Int64(`verifyCodeLifetime`, verifyCodeLifetime)
	if lifetime == 0 {
		lifetime = verifyCodeLifetime
	}
	expiry := now.Add(time.Duration(lifetime) * time.Minute)
	var message string
	if len(messages) > 0 {
		message = messages[0]
	}
	if len(message) == 0 {
		message = ctx.T(`亲爱的客户: %s，您正在进行手机号码验证，本次验证码为：%s (%d分钟内有效) [%s]`, m.Name, verifyCode, lifetime, baseCfg.String(`siteName`))
	} else {
		placeholders := map[string]string{
			`name`:     m.Name,
			`code`:     verifyCode,
			`lifeTime`: param.AsString(lifetime),
			`siteName`: baseCfg.String(`siteName`),
		}
		for find, to := range placeholders {
			message = strings.ReplaceAll(message, `{`+find+`}`, to)
		}
	}
	smsConfig.Content = message
	smsConfig.Template = ``
	smsConfig.SignName = ``
	smsConfig.ExtraData[`code`] = verifyCode

	// 记录日志
	vm.Verification.Code = verifyCode
	vm.Verification.OwnerId = m.Id
	vm.Verification.OwnerType = `customer`
	vm.Verification.Purpose = purpose
	vm.Verification.Start = uint(now.Unix())
	vm.Verification.End = uint(expiry.Unix())
	vm.Verification.Disabled = `N`
	vm.Verification.SendMethod = `mobile`
	vm.Verification.SendTo = m.Mobile
	if _, addErr := vm.AddVerificationCode(); addErr != nil {
		return addErr
	}
	logM := model.NewSendingLog(ctx)
	logM.Provider = smsProviderName
	logM.Method = `mobile`
	logM.To = m.Mobile
	logM.SourceType = `code_verification`
	logM.SourceId = vm.Verification.Id
	logM.Result = ctx.T(`发送成功`)
	logM.Status = `success`
	logM.Content = smsConfig.Content
	b, e := com.JSONEncode(smsConfig.ExtraData)
	if e != nil {
		return e
	}
	logM.Params = string(b)
	logM.AppointmentTime = 0
	if _, addErr := logM.Add(); addErr != nil {
		return addErr
	}
	timestamp := time.Now().Unix()
	smsConfig.CallbackURL = ctx.URLFor(`/verification/callback/` + smsProviderName + `/` + fmt.Sprint(vm.Verification.Id) + `/` + fmt.Sprint(timestamp) + `/` + uploadChecker.Token(smsProviderName, vm.Verification.Id, timestamp))

	err = provider.Send(smsConfig)
	if err != nil {
		logM.UpdateFields(nil, echo.H{
			`status`: `failure`,
			`result`: ctx.T(`发送失败`) + `: ` + err.Error(),
		}, `id`, logM.Id)
		return err
	}
	data.SetInfo(ctx.T(`短信发送成功`))
	return nil
}
