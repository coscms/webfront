package customer

import (
	"slices"
	"strings"
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/offlinepay"
	"github.com/coscms/webfront/library/top"
	"github.com/coscms/webfront/library/xcommon"
	"github.com/coscms/webfront/library/xrole"
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
	if len(u.Status) == 0 {
		u.Status = OfflinePayStatusPending
	} else if !slices.Contains(OfflinePayStatusAll, u.Status) {
		return u.Context().NewError(code.InvalidParameter, `状态无效`).SetZone(`status`)
	}
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

func (u *OfflinePay) SetPending() error {
	return u.OfficialCustomerOfflinePay.UpdateFields(nil, echo.H{
		`status`:  OfflinePayStatusPending,
		`updated`: time.Now().Unix(),
	}, `id`, u.Id)
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
		db.Cond{`status`: OfflinePayStatusPending},
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
	ctx := u.Context()
	if err := ctx.Begin(); err != nil {
		return nil, err
	}
	pk, err := u.OfficialCustomerOfflinePay.Insert()
	if err != nil {
		ctx.Rollback()
		return pk, err
	}
	if u.Id > 0 {
		err = u.fireEvent(ctx)
	}
	ctx.End(err == nil)
	return pk, err
}

func (u *OfflinePay) fireEvent(ctx echo.Context) error {
	switch u.Status {
	case OfflinePayStatusVerified:
		return FireVerifiedOfflinePayTargetType(ctx, u.OfficialCustomerOfflinePay)
	case OfflinePayStatusInvalid:
		return FireInvalidOfflinePayTargetType(ctx, u.OfficialCustomerOfflinePay)
	default:
		return nil
	}
}

func (u *OfflinePay) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := u.check(); err != nil {
		return err
	}
	ctx := u.Context()
	if err := ctx.Begin(); err != nil {
		return err
	}
	old := dbschema.NewOfficialCustomerOfflinePay(ctx)
	err := old.Get(nil, args...)
	if err != nil {
		ctx.Rollback()
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, ``)
		}
		return err
	}
	if old.Status == OfflinePayStatusVerified {
		ctx.Rollback()
		return ctx.NewError(code.DataUnavailable, `不能修改已经核实过的信息`).SetZone(`status`)
	}
	var affected int64
	affected, err = u.OfficialCustomerOfflinePay.Updatex(nil, args...)
	if err != nil {
		ctx.Rollback()
		return err
	}
	if old.Status != u.Status && affected > 0 {
		err = u.fireEvent(ctx)
	}
	ctx.End(err == nil)
	return err
}

func (f *OfflinePay) CustomerTodayCount(customerID interface{}) (int64, error) {
	startTs, endTs := top.TodayTimestamp()
	return f.Count(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`created`: db.Between(startTs, endTs)},
	))
}

func (f *OfflinePay) CustomerPendingCount(customerID interface{}) (int64, error) {
	return f.Count(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`status`: OfflinePayStatusPending},
	))
}

func (f *OfflinePay) CustomerPendingTodayCount(customerID interface{}) (int64, error) {
	startTs, endTs := top.TodayTimestamp()
	return f.Count(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`display`: OfflinePayStatusPending},
		db.Cond{`created`: db.Between(startTs, endTs)},
	))
}

func (f *OfflinePay) CheckCustomerAdd() error {
	permission := CustomerPermission(f.Context())
	return f.checkCustomerAdd(permission)
}

func (f *OfflinePay) checkCustomerAdd(permission *xrole.RolePermission) error {
	err := xcommon.CheckRoleCustomerAdd(f.Context(), permission, OfflinePayBehaviorName, f.CustomerId, f)
	if err == nil {
		return err
	}
	switch err {
	case xcommon.ErrCustomerRoleDisabled:
		return f.Context().E(`当前角色不支持线下转账`)
	case xcommon.ErrCustomerAddClosed:
		return f.Context().E(`线下转账功能已关闭`)
	case xcommon.ErrCustomerAddMaxPerDay:
		return f.Context().E(`提交线下转账信息失败。您的账号已达到今日最大提交数量`)
	case xcommon.ErrCustomerAddMaxPending:
		return f.Context().E(`提交线下转账信息失败。您的待核实信息数量已达上限，请等待审核通过后再提交`)
	default:
		return err
	}
}
