package customer

import (
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/echo"
)

const (
	OfflinePayStatusPending  = `pending`
	OfflinePayStatusVerified = `verified`
	OfflinePayStatusInvalid  = `invalid`
)

var OfflinePayStatusAll = []string{
	OfflinePayStatusPending,
	OfflinePayStatusVerified,
	OfflinePayStatusInvalid,
}

var OfflinePayStatuses = echo.NewKVData().
	Add(OfflinePayStatusPending, echo.T(`待处理`)).
	Add(OfflinePayStatusVerified, echo.T(`已核实`)).
	Add(OfflinePayStatusInvalid, echo.T(`无效`))

const (
	OfflinePayTargetTypeRecharge = `recharge`
)

type OfflinePayTargetTypeX interface {
	Verified(ctx echo.Context, row *dbschema.OfficialCustomerOfflinePay) error
	Invalid(ctx echo.Context, row *dbschema.OfficialCustomerOfflinePay) error
}

var OfflinePayTargetTypes = echo.NewKVxData[OfflinePayTargetTypeX, any]().
	Add(OfflinePayTargetTypeRecharge, echo.T(`钱包充值`), echo.KVxOptX[OfflinePayTargetTypeX, any](offlinePayTargetTypeRecharge{}))

func RegisterOfflinePayTargetType(k, v string, x OfflinePayTargetTypeX) {
	OfflinePayTargetTypes.Add(k, v, echo.KVxOptX[OfflinePayTargetTypeX, any](x))
}

func FireVerifiedOfflinePayTargetType(ctx echo.Context, row *dbschema.OfficialCustomerOfflinePay) error {
	item := OfflinePayTargetTypes.GetItem(row.TargetType)
	if item == nil {
		return nil
	}
	if item.X == nil {
		return nil
	}
	return item.X.Verified(ctx, row)
}

func FireInvalidOfflinePayTargetType(ctx echo.Context, row *dbschema.OfficialCustomerOfflinePay) error {
	item := OfflinePayTargetTypes.GetItem(row.TargetType)
	if item == nil {
		return nil
	}
	if item.X == nil {
		return nil
	}
	return item.X.Invalid(ctx, row)
}

type offlinePayTargetTypeRecharge struct{}

func (a offlinePayTargetTypeRecharge) Verified(ctx echo.Context, row *dbschema.OfficialCustomerOfflinePay) error {
	walletM := NewWallet(ctx)
	walletM.Flow.CustomerId = row.CustomerId
	walletM.Flow.AssetType = AssetTypeMoney
	walletM.Flow.AmountType = AmountTypeBalance
	walletM.Flow.Amount = row.PayAmount
	walletM.Flow.SourceType = `recharge`
	walletM.Flow.SourceTable = `official_customer_offline_pay`
	walletM.Flow.SourceId = row.Id
	walletM.Flow.TradeNo = row.PayTransactionNo
	walletM.Flow.Status = FlowStatusConfirmed //状态(pending-待确认;confirmed-已确认;canceled-已取消)
	walletM.Flow.Description = `线下转账充值`
	return walletM.AddFlow()
}

func (a offlinePayTargetTypeRecharge) Invalid(ctx echo.Context, row *dbschema.OfficialCustomerOfflinePay) error {
	return nil
}
