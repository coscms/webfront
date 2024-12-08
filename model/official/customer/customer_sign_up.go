package customer

import (
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/model"
	modelLevel "github.com/coscms/webfront/model/official/level"
)

// SignUp 注册用户
func (f *Customer) SignUp(user, pass, mobile, email string, options ...CustomerOption) error {
	co := NewCustomerOptions(f.OfficialCustomer, true)
	co.Name = user
	co.Password = pass
	co.Mobile = mobile
	co.Email = email
	co.ApplyOptions(options...)
	if f.LevelId < 1 {
		levelM := modelLevel.NewLevel(f.Context())
		if level, err := levelM.CanAutoLevelUpByIntegral(0); err == nil {
			f.LevelId = level.Id
		}
	}
	if len(f.RoleIds) == 0 {
		roleM := NewRole(f.Context())
		if err := roleM.GetDefault(); err == nil {
			f.RoleIds = param.AsString(roleM.Id)
		}
	}
	f.SessionId = f.Context().Session().MustID()
	if inviter := f.Context().GetCookie(`inviter`); len(inviter) > 0 {
		if inviterID := param.AsUint64(inviter); inviterID > 0 {
			exists, _ := f.ExistsCustomerID(inviterID)
			if exists {
				f.InviterId = inviterID
			}
		}
	}
	_, err := f.Add()
	if err != nil {
		return err
	}

	return f.FireSignUpSuccess(co, model.AuthTypePassword)
}

func (f *Customer) FireSignUpSuccess(co *CustomerOptions, authType string) (err error) {
	integral := config.Setting(`base`, `addExperience`).Float64(`register`)
	if err = FireSignUp(f.OfficialCustomer); err != nil {
		return err
	}
	if err = f.AddRewardOnSignUp(integral); err != nil {
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
	deviceM.SessionId = f.SessionId
	deviceM.CustomerId = f.Id
	deviceM.SetOptions(co)
	_, err = deviceM.Upsert()
	if err != nil {
		return err
	}

	f.SetSession()
	if !f.disabledSession {
		co.SetSession(f.Context())
	}

	loginLogM := f.NewLoginLog(co, authType)
	loginLogM.Success = `Y`
	if f.disabledSession {
		loginLogM.InitLocation()
		loginLogM.Add()
	} else {
		loginLogM.AddAndSaveSession()
	}
	return
}

func (f *Customer) UpdateInviterID(customerID uint64, inviterID uint64) error {
	return f.OfficialCustomer.UpdateField(nil, `inviter_id`, inviterID, `id`, customerID)
}

func (f *Customer) ExistsCustomerID(customerID uint64) (bool, error) {
	return f.OfficialCustomer.Exists(nil, `id`, customerID)
}
