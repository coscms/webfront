package apiutils

import (
	"github.com/admpub/null"
	"github.com/coscms/sdk/sdk_options"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewOptions(ctx echo.Context, typ Type, generators ...sdk_options.URLValuesGenerator) *Options {
	var generator sdk_options.URLValuesGenerator
	if len(generators) > 0 {
		generator = generators[0]
	}
	base := sdk_options.New(typ, nil)
	if generator != nil {
		base.SetGenerator(generator)
	}
	return &Options{
		Options: base,
		ctx:     ctx,
	}
}

type AppInfoGetter func(ctx echo.Context, cond db.Compound) (appInfo AppInfo, err error)

type Options struct {
	*sdk_options.Options
	ctx            echo.Context
	applied        bool
	accountID      null.Uint64
	_appInfoGetter AppInfoGetter
}

func (o *Options) SetAppInfoGetter(appInfoGetter AppInfoGetter) *Options {
	o._appInfoGetter = appInfoGetter
	return o
}

func (o *Options) GetAccountID() uint64 {
	if o.accountID.Valid {
		return o.accountID.Uint64
	}
	o.accountID.Valid = true
	o.accountID.Uint64 = config.Setting(`thirdparty`).Uint64(string(o.Type))
	return o.accountID.Uint64
}

func (o *Options) ApplySetting() (err error) {
	if o.applied {
		return
	}
	o.applied = true
	accountID := o.GetAccountID()
	if accountID > 0 {
		err = o.getAccount(db.Cond{`id`: accountID})
	} else {
		err = o.getApp(db.Cond{`owner_type`: `official`})
	}
	return
}

func (o *Options) Context() echo.Context {
	return o.ctx
}
