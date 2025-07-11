package customer

import (
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/perm"
	"github.com/coscms/webcore/model"
	multidevicesignin "github.com/coscms/webfront/library/multidevicesignin"
	"github.com/coscms/webfront/library/xrole"
)

// SignIn 用户登录
func (f *Customer) SignIn(user, pass, signInType string, options ...CustomerOption) error {
	if len(signInType) == 0 {
		signInType = `name`
	}
	co := NewCustomerOptions(nil)
	co.Name = user
	co.Password = pass
	co.SignInType = signInType
	co.ApplyOptions(options...)
	baseCfg := config.Setting(`base`)
	if baseCfg.String(`customerLogin`, `open`) == `close` {
		return f.Context().NewError(code.DataStatusIncorrect, `本站已经暂时关闭登录，请稍后再尝试`)
	}
	cond := db.NewCompounds()
	fieldName := `name`
	switch co.SignInType {
	case `email`:
		if _err := f.Context().Validate(fieldName, co.Name, `email`); _err != nil {
			return f.Context().NewError(code.InvalidParameter, `E-mail格式不正确`).SetZone(fieldName)
		}
		cond.AddKV(`email`, co.Name)
		cond.AddKV(`email_bind`, `Y`)
	case `mobile`:
		if _err := f.Context().Validate(fieldName, co.Name, `mobile`); _err != nil {
			return f.Context().NewError(code.InvalidParameter, `手机号码格式不正确`).SetZone(fieldName)
		}
		cond.AddKV(`mobile`, co.Name)
		cond.AddKV(`mobile_bind`, `Y`)
	case `name`:
		if !com.IsUsername(co.Name) {
			return f.Context().NewError(code.UserNotFound, `用户名无效`).SetZone(fieldName)
		}
		cond.AddKV(`name`, co.Name)
	default:
		return f.Context().NewError(code.Unsupported, `不支持登录方式: %v`, co.SignInType).SetZone(`type`)
	}
	if len(co.Name) == 0 {
		return f.Context().NewError(code.InvalidParameter, `请输入登录名称`).SetZone(co.SignInType)
	}
	if len(co.Password) == 0 {
		return f.Context().NewError(code.InvalidParameter, `请输入登录密码`).SetZone(`password`)
	}
	err := f.Get(nil, cond.And())
	if err != nil {
		loginLogM := f.NewLoginLog(co, model.AuthTypePassword)
		loginLogM.Errpwd = co.Password
		if err == db.ErrNoMoreRows {
			loginLogM.Failmsg = f.Context().T(`用户不存在`)
			loginLogM.Add()
			return f.Context().NewError(code.UserNotFound, `用户不存在`)
		}
		loginLogM.Failmsg = err.Error()
		loginLogM.Add()
		return err
	}
	if siteClose := baseCfg.Uint(`siteClose`); siteClose == 3 && f.Uid < 1 {
		return f.Context().NewError(code.NonPrivileged, `网站暂时关闭，仅供管理员访问`)
	}
	if err = f.CheckSignInPassword(co.Password); err != nil {
		if !echo.IsErrorCode(err, code.UserDisabled) {
			// 仅记录密码不正确的情况
			loginLogM := f.NewLoginLog(co, model.AuthTypePassword)
			loginLogM.OwnerId = f.Id
			loginLogM.Errpwd = co.Password
			loginLogM.Failmsg = err.Error()
			loginLogM.Add()
			f.IncrLoginFails()
		}
		return err
	}
	return f.FireSignInSuccess(co, model.AuthTypePassword)
}

func (f *Customer) FireSignInSuccess(co *CustomerOptions, authType string) (err error) {
	loginLogM := f.NewLoginLog(co, authType)
	loginLogM.OwnerId = f.Id
	set := echo.H{
		`login_fails`: 0,
	}
	ctx := f.Context()
	err = ctx.Begin()
	if err != nil {
		return
	}
	defer func() {
		ctx.End(err == nil)
		if err != nil {
			loginLogM.Failmsg = err.Error()
			loginLogM.Add()
		} else {
			loginLogM.Success = `Y`
			if f.disabledSession {
				loginLogM.InitLocation()
				loginLogM.Add()
			} else {
				loginLogM.AddAndSaveSession()
			}
		}
	}()
	if err = f.LevelUpOnSignIn(set); err != nil {
		return err
	}
	integral := config.Setting(`base`, `addExperience`).Float64(`login`)
	if err = f.AddRewardOnSignIn(integral); err != nil {
		return err
	}

	err = f.LinkOAuthUser()
	if err != nil {
		return err
	}
	err = FireSignIn(f.OfficialCustomer)
	if err != nil {
		return err
	}

	deviceM := NewDevice(f.Context())
	deviceM.SessionId = loginLogM.SessionId
	deviceM.CustomerId = f.Id
	deviceM.SetOptions(co)
	_, err = deviceM.Upsert()
	if err != nil {
		return err
	}
	if len(f.SessionId) > 0 {
		if f.SessionId != loginLogM.SessionId {
			set[`session_id`] = loginLogM.SessionId
			f.SessionId = loginLogM.SessionId
			err = deviceM.CleanCustomer(f.OfficialCustomer, co)
		} else {
			permission := CustomerPermission(f.Context(), f.OfficialCustomer)
			if permission != nil {
				if bev, ok := permission.Get(f.Context(), xrole.CustomerRolePermissionTypeBehavior).(perm.BehaviorPerms); ok {
					multideviceSignin, _ := bev.Get(multidevicesignin.BehaviorName).Value.(*multidevicesignin.MultideviceSignin)
					err = deviceM.CleanExceedLimit(deviceM.CustomerId, multideviceSignin)
				}
			}
		}
		if err != nil {
			return err
		}
	} else {
		set[`session_id`] = loginLogM.SessionId
		f.SessionId = loginLogM.SessionId
	}
	if len(set) > 0 {
		err = f.UpdateFields(nil, set, `id`, f.Id)
		if err != nil {
			return err
		}
	}

	f.SetSession()
	if !f.disabledSession {
		co.SetSession(f.Context())
	}
	return err
}

func (f *Customer) LinkOAuthUser() error {
	oAuthM := NewOAuth(f.Context())
	oAuthUser, exists, err := oAuthM.GetSession()
	if err != nil {
		if !exists {
			return nil
		}
		return err
	}
	if oAuthUser != nil {
		oAuthM.CustomerId = f.Id
		oAuthM.CopyFrom(oAuthUser)
		_, err := oAuthM.Add()
		if err != nil {
			log.Errorf(`%v: %#v (oAuthUser: %#v)`, err, *oAuthM.OfficialCustomerOauth, *oAuthUser)
			return err
		}
		oAuthM.DelSession()
	}
	return nil
}
