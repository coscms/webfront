package customer

import (
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webfront/dbschema"
)

func NewOnline(ctx echo.Context) *Online {
	m := &Online{
		OfficialCustomerOnline: dbschema.NewOfficialCustomerOnline(ctx),
	}
	return m
}

type Online struct {
	*dbschema.OfficialCustomerOnline
}

func (u *Online) MakeCond(sessionID string, customerID uint64) db.Compound {
	if customerID == 0 {
		return db.Cond{`session_id`: sessionID}
	}
	return db.Or(
		db.Cond{`session_id`: sessionID},
		db.Cond{`customer_id`: customerID},
	)
}

func (u *Online) Exists() (bool, error) {
	return u.OfficialCustomerOnline.Exists(nil, u.MakeCond(u.SessionId, u.CustomerId))
}

func (u *Online) check() error {
	exists, err := u.Exists()
	if err != nil {
		return err
	}
	if exists {
		err = u.Context().NewError(code.DataAlreadyExists, `数据已经存在`)
	}
	return err
}

func (u *Online) Add() (interface{}, error) {
	if err := u.check(); err != nil {
		return nil, err
	}
	return u.OfficialCustomerOnline.Insert()
}

func (u *Online) Incr(n uint) error {
	old := dbschema.NewOfficialCustomerOnline(u.Context())
	err := old.Get(func(r db.Result) db.Result {
		return r.Select(`session_id`, `customer_id`, `updated`)
	}, u.MakeCond(u.SessionId, u.CustomerId))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		u.ClientCount = n
		u.Updated = uint(time.Now().Unix())
		_, err = u.OfficialCustomerOnline.Insert()
	} else {
		kvset := echo.H{
			`updated`: time.Now().Unix(),
		}
		if old.SessionId != u.SessionId {
			kvset[`session_id`] = u.SessionId
		}
		if old.CustomerId != u.CustomerId {
			kvset[`customer_id`] = u.CustomerId
		}
		if int64(old.Updated) < com.StartTime.Unix() { // 程序启动之前的状态需要重新设置为1
			kvset[`client_count`] = 1
		} else {
			kvset[`client_count`] = db.Raw("client_count+" + param.AsString(n))
		}
		err = u.OfficialCustomerOnline.UpdateFields(nil, kvset, u.MakeCond(u.SessionId, u.CustomerId))
	}
	if err != nil {
		return err
	}
	customerM := dbschema.NewOfficialCustomer(u.Context())
	err = customerM.UpdateField(nil, `online`, `Y`, db.And(
		db.Cond{`id`: u.CustomerId},
		db.Cond{`online`: `N`},
	))
	return err
}

func (u *Online) Decr(n uint64) error {
	old := dbschema.NewOfficialCustomerOnline(u.Context())
	err := old.Get(func(r db.Result) db.Result {
		return r.Select(`client_count`, `session_id`, `customer_id`)
	}, u.MakeCond(u.SessionId, u.CustomerId))
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return err
	}
	if old.ClientCount <= 1 {
		customerM := dbschema.NewOfficialCustomer(u.Context())
		customerM.UpdateField(nil, `online`, `N`, db.And(
			db.Cond{`id`: u.CustomerId},
			db.Cond{`online`: `Y`},
		))
	}
	kvset := echo.H{
		`client_count`: db.Raw("client_count-" + param.AsString(n)),
		`updated`:      time.Now().Unix(),
	}
	if old.SessionId != u.SessionId {
		kvset[`session_id`] = u.SessionId
	}
	if old.CustomerId != u.CustomerId {
		kvset[`customer_id`] = u.CustomerId
	}
	return u.OfficialCustomerOnline.UpdateFields(nil, kvset, db.And(
		db.Cond{`session_id`: old.SessionId},
		db.Cond{`customer_id`: old.CustomerId},
		db.Cond{`client_count`: db.Gt(0)},
	))
}

func (u *Online) Cleanup() error {
	return u.OfficialCustomerOnline.Delete(nil, db.Cond{`client_count`: 0})
}

func (u *Online) IsOnlineCustomerIDs(customerIDs []uint64) map[uint64]bool {
	u.OfficialCustomerOnline.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`customer_id`)
	}, 0, -1, db.And(
		db.Cond{`customer_id`: db.In(customerIDs)},
		db.Cond{`client_count`: db.Gt(0)},
		db.Cond{`updated`: db.Gte(com.StartTime.Unix())},
	))
	exists := map[uint64]bool{}
	for _, row := range u.Objects() {
		exists[row.CustomerId] = true
	}
	return exists
}

func (u *Online) IsOnlineCustomerID(customerID uint64) bool {
	exists, _ := u.OfficialCustomerOnline.Exists(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`client_count`: db.Gt(0)},
		db.Cond{`updated`: db.Gte(com.StartTime.Unix())},
	))
	return exists
}

func (u *Online) ResetClientCount(isDelete ...bool) error {
	var err error
	if len(isDelete) > 0 && isDelete[0] {
		err = u.OfficialCustomerOnline.Delete(nil, db.Cond{`updated`: db.Lt(com.StartTime.Unix())})
	} else {
		err = u.OfficialCustomerOnline.UpdateField(nil, `client_count`, 0, db.And(
			db.Cond{`client_count`: db.NotEq(0)},
			db.Cond{`updated`: db.Lt(com.StartTime.Unix())},
		))
	}
	if err != nil {
		return err
	}
	customerM := dbschema.NewOfficialCustomer(u.Context())
	err = customerM.UpdateField(nil, `online`, `N`, db.Cond{`online`: `Y`})
	return err
}
