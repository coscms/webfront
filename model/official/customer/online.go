package customer

import (
	"time"

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

func (u *Online) Exists() (bool, error) {
	return u.OfficialCustomerOnline.Exists(nil, db.And(
		db.Cond{`session_id`: u.SessionId},
		db.Cond{`customer_id`: u.CustomerId},
	))
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
	exists, err := u.Exists()
	if err != nil {
		return err
	}
	if !exists {
		u.ClientCount = n
		u.Updated = uint(time.Now().Unix())
		_, err = u.OfficialCustomerOnline.Insert()
		return err
	}
	return u.OfficialCustomerOnline.UpdateFields(nil, echo.H{
		`client_count`: db.Raw("client_count+" + param.AsString(n)),
		`updated`:      time.Now().Unix(),
	}, db.And(
		db.Cond{`session_id`: u.SessionId},
		db.Cond{`customer_id`: u.CustomerId},
	))
}

func (u *Online) Decr(n uint64) error {
	exists, err := u.Exists()
	if err != nil || !exists {
		return err
	}
	return u.OfficialCustomerOnline.UpdateFields(nil, echo.H{
		`client_count`: db.Raw("client_count-" + param.AsString(n)),
		`updated`:      time.Now().Unix(),
	}, db.And(
		db.Cond{`session_id`: u.SessionId},
		db.Cond{`customer_id`: u.CustomerId},
		db.Cond{`client_count`: db.Gt(0)},
	))
}

func (u *Online) Cleanup() error {
	return u.OfficialCustomerOnline.Delete(nil, db.Cond{`client_count`: 0})
}

func (u *Online) IsOnlineCustomerIDs(customerIDs []uint64) map[uint64]bool {
	u.OfficialCustomerOnline.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`)
	}, 0, -1, db.And(
		db.Cond{`customer_id`: db.In(customerIDs)},
		db.Cond{`client_count`: db.Gt(0)},
	))
	exists := map[uint64]bool{}
	for _, row := range u.Objects() {
		exists[row.CustomerId] = true
	}
	return exists
}

func (u *Online) IsOnlineCustomerID(customerID uint64) bool {
	exists, _ := u.OfficialCustomerOnline.Exists(nil, func(r db.Result) db.Result {
		return r.Select(`id`)
	}, 0, -1, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`client_count`: db.Gt(0)},
	))
	return exists
}
