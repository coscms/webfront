package customer

import (
	"github.com/admpub/log"
	"github.com/webx-top/db"

	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webcore/library/sessionguard"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/middleware/sessdata"
)

// SetSession 记录登录信息
func (f *Customer) SetSession(customers ...*dbschema.OfficialCustomer) {
	if f.disabledSession {
		return
	}
	customerCopy := f.ClearPasswordData(customers...)
	f.Context().Session().Set(`customer`, &customerCopy)
}

// UnsetSession 退出登录
func (f *Customer) UnsetSession() error {
	if f.disabledSession {
		return nil
	}
	deviceM := NewDevice(f.Context())
	hasSignedInOtherDevice, err := deviceM.ExistsCustomerID(f.Id)
	if err != nil {
		f.Context().Logger().Error(err)
	} else if hasSignedInOtherDevice {
		f.Context().Internal().Set(`hasSignedInOtherDevice`, true)
	}
	err = FireSignOut(f.OfficialCustomer)
	f.Context().Session().Delete(`customer`)
	if !hasSignedInOtherDevice {
		sessdata.ClearPermissionCache(f.Context(), f.OfficialCustomer.Id)
	}
	return err
}

func (f *Customer) VerifySession(customers ...*dbschema.OfficialCustomer) error {
	var customer *dbschema.OfficialCustomer
	if len(customers) > 0 {
		customer = customers[0]
	} else {
		customer, _ = f.Context().Session().Get(`customer`).(*dbschema.OfficialCustomer)
	}
	if customer == nil {
		return nerrors.ErrUserNotLoggedIn
	}
	detail, err := f.GetDetail(db.Cond{`id`: customer.Id})
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		f.UnsetSession()
		return nerrors.ErrUserNotFound
	}
	sessionID := f.Context().Session().ID()
	if detail.OfficialCustomer.SessionId != sessionID {
		var exists bool
		if len(sessionID) > 0 {
			deviceM := NewDevice(f.Context())
			exists, err = deviceM.ExistsSessionID(sessionID)
			if err != nil {
				return err
			}
		}
		if !exists {
			f.UnsetSession()
			return nerrors.ErrUserNotLoggedIn
		}
	}
	if !sessionguard.Validate(f.Context(), ``, `customer`, detail.Id) {
		log.Warn(f.Context().T(`客户“%s”的会话环境发生改变，需要重新登录`, detail.Name))
		f.UnsetSession()
		return nerrors.ErrUserNotLoggedIn
	}
	if detail.OfficialCustomer.Updated != customer.Updated {
		f.SetSession(detail.OfficialCustomer)
		f.Context().Internal().Set(`customer`, detail.OfficialCustomer)
	}
	if detail.OfficialCustomer.RoleIds != customer.RoleIds {
		sessdata.ClearPermissionCache(f.Context(), detail.OfficialCustomer.Id)
	}

	safeCustomer := f.ClearPasswordData(detail.OfficialCustomer)
	detail.OfficialCustomer = &safeCustomer
	f.Context().Internal().Set(`customerDetail`, detail)
	return nil
}
