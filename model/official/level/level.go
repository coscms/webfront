package level

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webfront/dbschema"
)

// GroupList base-普通会员(基础组),其它名称为扩展组。客户只能有一个基础组等级,可以有多个扩展组等级
var GroupList = echo.NewKVData().Add(`base`, echo.T(`普通会员`))

func AddGroup(k string, v string) {
	GroupList.Add(k, v)
}

func NewLevel(ctx echo.Context) *Level {
	m := &Level{
		OfficialCustomerLevel: dbschema.NewOfficialCustomerLevel(ctx),
	}
	return m
}

type Level struct {
	*dbschema.OfficialCustomerLevel
}

func (f *Level) GreaterOrEqualThan(levelId uint, targetLevelId uint) (bool, error) {
	return f.Than(levelId, targetLevelId, compareGreaterOrEqual)
}

func (f *Level) LessThan(levelId uint, targetLevelId uint) (bool, error) {
	return f.Than(levelId, targetLevelId, compareLess)
}

func (f *Level) Than(
	levelId uint, targetLevelId uint,
	compare func(*dbschema.OfficialCustomerLevel, *dbschema.OfficialCustomerLevel) bool,
) (bool, error) {
	row := dbschema.NewOfficialCustomerLevel(f.Context())
	err := row.Get(func(r db.Result) db.Result {
		return r.Select(`score`, `group`)
	}, `id`, levelId)
	if err != nil {
		return false, err
	}
	err = f.Get(func(r db.Result) db.Result {
		return r.Select(`score`)
	}, db.And(
		db.Cond{`id`: targetLevelId},
		db.Cond{`group`: row.Group},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = nil
		}
		return false, err
	}
	return compare(row, f.OfficialCustomerLevel), err
}

func (f *Level) check() error {
	if len(f.Group) == 0 {
		return f.Context().NewError(code.InvalidParameter, `group is required`).SetZone(`group`)
	}
	if !GroupList.Has(f.Group) {
		var validGroups string
		for i, g := range GroupList.Slice() {
			if i > 0 {
				validGroups += `, `
			}
			validGroups += g.K + `(` + g.V + `)`
		}
		return f.Context().NewError(code.InvalidParameter, `group无效(仅支持: %v)`, validGroups).SetZone(`group`)
	}
	if len(f.IntegralAmountType) == 0 {
		f.IntegralAmountType = AmountTypeBalance
	} else if !AmountTypes.Has(f.IntegralAmountType) {
		var validAmountTypes string
		for i, g := range AmountTypes.Slice() {
			if i > 0 {
				validAmountTypes += `, `
			}
			validAmountTypes += g.K + `(` + g.V + `)`
		}
		return f.Context().NewError(code.InvalidParameter, `integralAmountType无效(仅支持: %v)`, validAmountTypes).SetZone(`integralAmountType`)
	}
	return nil
}

func (f *Level) Add() (pk interface{}, err error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	if err := f.Exists(f.Name); err != nil {
		return nil, err
	}
	return f.OfficialCustomerLevel.Insert()
}

func (f *Level) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	if err := f.ExistsOther(f.Name, f.Id); err != nil {
		return err
	}
	return f.OfficialCustomerLevel.Update(mw, args...)
}

func (f *Level) Exists(name string) error {
	exists, err := f.OfficialCustomerLevel.Exists(nil, db.Cond{`name`: name})
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().NewError(code.DataAlreadyExists, `等级名称“%s”已经使用过了`, name).SetZone(`name`)
	}
	return err
}

func (f *Level) ListByCustomer(customer *dbschema.OfficialCustomer) ([]*dbschema.OfficialCustomerLevel, error) {
	r := NewRelation(f.Context())
	rows, err := r.ListByCustomerID(customer.Id)
	if err != nil {
		return nil, err
	}
	levelIDs := make([]uint, len(rows))
	for index, row := range rows {
		levelIDs[index] = row.LevelId
	}
	if len(levelIDs) == 0 {
		return nil, err
	}
	_, err = f.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`id`: db.In(levelIDs)},
		db.Cond{`disabled`: `N`},
	))
	if err != nil {
		return nil, err
	}
	return f.Objects(), err
}

func (f *Level) ExistsOther(name string, id uint) error {
	exists, err := f.OfficialCustomerLevel.Exists(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`id`: db.NotEq(id)},
	))
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().NewError(code.DataAlreadyExists, `等级名称“%s”已经使用过了`, name).SetZone(`name`)
	}
	return err
}

func (f *Level) ListLevelGroup() ([]*LevelGroup, error) {
	var list []*LevelGroup
	_, err := f.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`group`, `-score`, `id`)
	}, 0, -1, db.Cond{`disabled`: `N`})
	if err != nil {
		return list, err
	}
	groupRows := map[string][]*dbschema.OfficialCustomerLevel{}
	var groups []string
	for _, v := range f.Objects() {
		_, ok := groupRows[v.Group]
		if !ok {
			groups = append(groups, v.Group)
			groupRows[v.Group] = []*dbschema.OfficialCustomerLevel{v}
		} else {
			groupRows[v.Group] = append(groupRows[v.Group], v)
		}
	}
	list = make([]*LevelGroup, len(groups))
	for index, group := range groups {
		title := GroupList.Get(group, group)
		list[index] = &LevelGroup{
			Group: group,
			Title: title,
			List:  groupRows[group],
		}
	}
	return list, err
}

func (f *Level) CanAutoLevelUpByCustomerID(customerID uint64) (*dbschema.OfficialCustomerLevel, error) {
	return f.CanAutoLevelUpBy(customerID, `base`, `integral`)
}

func (f *Level) CanAutoLevelUpBy(customerID uint64, group string, assetType string) (*dbschema.OfficialCustomerLevel, error) {
	walletM := dbschema.NewOfficialCustomerWallet(nil)
	walletM.CPAFrom(f.OfficialCustomerLevel)
	err := walletM.Get(func(r db.Result) db.Result {
		return r.Select(`balance`, `accumulated`)
	}, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`asset_type`: assetType},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return nil, err
		}
	}
	return f.CanAutoLevelUpByIntegralAsset(group, walletM.Balance, walletM.Accumulated, assetType)
}

func (f *Level) CanAutoLevelUpByIntegral(balance float64, accumulated float64) (*dbschema.OfficialCustomerLevel, error) {
	return f.CanAutoLevelUpByIntegralAsset(`base`, balance, accumulated, `integral`)
}

func (f *Level) CanAutoLevelUpByIntegralAsset(group string, balance float64, accumulated float64, asset string) (*dbschema.OfficialCustomerLevel, error) {
	cond := MakeFreeCond(group, balance, accumulated, asset)
	err := f.Get(func(r db.Result) db.Result {
		return r.OrderBy(`-score`)
	}, cond.And())
	if err != nil {
		if err != db.ErrNoMoreRows {
			return nil, err
		}
	}
	return f.OfficialCustomerLevel, nil
}

func (f *Level) CanPaymentLevelUpByIntegralAsset(group string, balance float64, accumulated float64, asset string) (*dbschema.OfficialCustomerLevel, error) {
	cond := MakePayCond(group, balance, accumulated, asset)
	err := f.Get(func(r db.Result) db.Result {
		return r.OrderBy(`-score`)
	}, cond.And())
	if err != nil {
		if err != db.ErrNoMoreRows {
			return nil, err
		}
	}
	return f.OfficialCustomerLevel, nil
}

func (f *Level) GetMinLevelByGroup(group string) error {
	return f.Get(func(r db.Result) db.Result {
		return r.OrderBy(`score`, `integral_min`)
	}, db.And(
		db.Cond{`disabled`: `N`},
		db.Cond{`group`: group},
	))
}
