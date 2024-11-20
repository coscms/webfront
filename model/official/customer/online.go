package customer

import (
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

func (u *Online) Incr(target string, n uint) error {
	exists, err := u.Exists()
	if err != nil {
		return err
	}
	if !exists {
		u.ClientCount = n
		_, err = u.OfficialCustomerOnline.Insert()
		return err
	}
	return u.OfficialCustomerOnline.UpdateField(nil, `client_count`, db.Raw("client_count+"+param.AsString(n)), db.And(
		db.Cond{`session_id`: u.SessionId},
		db.Cond{`customer_id`: u.CustomerId},
	))
}

func (u *Online) Decr(target string, n uint64) error {
	exists, err := u.Exists()
	if err != nil || !exists {
		return err
	}
	return u.OfficialCustomerOnline.UpdateField(nil, `client_count`, db.Raw("client_count-"+param.AsString(n)), db.And(
		db.Cond{`session_id`: u.SessionId},
		db.Cond{`customer_id`: u.CustomerId},
	))
}

func (u *Online) Cleanup() error {
	return u.OfficialCustomerOnline.Delete(nil, db.Cond{`client_count`: 0})
}
