package customer

import (
	"strings"
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/offlinepay"
)

func NewOfflinePay(ctx echo.Context) *OfflinePay {
	m := &OfflinePay{
		OfficialCustomerOfflinePay: dbschema.NewOfficialCustomerOfflinePay(ctx),
	}
	return m
}

type OfflinePay struct {
	*dbschema.OfficialCustomerOfflinePay
}

func (u *OfflinePay) check() error {
	u.TargetType = strings.TrimSpace(u.TargetType)
	u.PayMethod = strings.TrimSpace(u.PayMethod)
	u.PayAccount = strings.TrimSpace(u.PayAccount)
	if len(u.TargetType) == 0 {
		return u.Context().NewError(code.InvalidParameter, `目标类型无效`).SetZone(`targetType`)
	}
	if !OfflinePayTargetTypes.Has(u.TargetType) {
		return u.Context().NewError(code.InvalidParameter, `目标类型无效`).SetZone(`targetType`)
	}
	if len(u.PayAccount) == 0 {
		return u.Context().NewError(code.InvalidParameter, `付款账号无效`).SetZone(`payAccount`)
	}
	if offlinepay.GetMethod(u.PayMethod, nil) == nil {
		return u.Context().NewError(code.InvalidParameter, `付款方式无效`).SetZone(`payMethod`)
	}
	if u.PayAmount <= 0 {
		return u.Context().NewError(code.InvalidParameter, `请输入付款金额`).SetZone(`payAmount`)
	}
	u.PayBankBranch = strings.TrimSpace(u.PayBankBranch)
	u.PayTransactionNo = strings.TrimSpace(u.PayTransactionNo)
	if u.PayTime > 0 && int64(u.PayTime) > time.Now().AddDate(0, 0, 1).Unix() {
		return u.Context().NewError(code.InvalidParameter, `付款时间无效`).SetZone(`payTime`)
	}
	u.PayOwner = strings.TrimSpace(u.PayOwner)
	u.Status = OfflinePayStatusPending
	return nil
}

func (u *OfflinePay) SetVerified() error {
	ctx := u.Context()
	if err := ctx.Begin(); err != nil {
		return err
	}
	affected, err := u.OfficialCustomerOfflinePay.UpdatexFields(nil, echo.H{
		`status`:  OfflinePayStatusVerified,
		`updated`: time.Now().Unix(),
	}, db.And(
		db.Cond{`id`: u.Id},
		db.Cond{`status`: db.NotEq(OfflinePayStatusVerified)},
	))
	if err != nil {
		ctx.Rollback()
		return err
	}
	if affected < 1 {
		ctx.Commit()
		return err
	}
	u.Status = OfflinePayStatusVerified
	err = FireVerifiedOfflinePayTargetType(ctx, u.OfficialCustomerOfflinePay)
	ctx.End(err == nil)
	return err
}

func (u *OfflinePay) SetInvalid() error {
	ctx := u.Context()
	if err := ctx.Begin(); err != nil {
		return err
	}
	affected, err := u.OfficialCustomerOfflinePay.UpdatexFields(nil, echo.H{
		`status`:  OfflinePayStatusInvalid,
		`updated`: time.Now().Unix(),
	}, db.And(
		db.Cond{`id`: u.Id},
		db.Cond{`status`: db.NotEq(OfflinePayStatusPending)},
	))
	if err != nil {
		ctx.Rollback()
		return err
	}
	if affected < 1 {
		ctx.Commit()
		return err
	}
	u.Status = OfflinePayStatusInvalid
	err = FireInvalidOfflinePayTargetType(ctx, u.OfficialCustomerOfflinePay)
	ctx.End(err == nil)
	return err
}

func (u *OfflinePay) Add() (interface{}, error) {
	if err := u.check(); err != nil {
		return nil, err
	}
	return u.OfficialCustomerOfflinePay.Insert()
}
