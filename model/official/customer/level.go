package customer

import (
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	modelLevel "github.com/coscms/webfront/model/official/level"
)

// NewLevel 客户关联等级信息
func NewLevel(ctx echo.Context) *Level {
	m := &Level{
		OfficialCustomerLevelRelation: dbschema.NewOfficialCustomerLevelRelation(ctx),
	}
	return m
}

type Level struct {
	*dbschema.OfficialCustomerLevelRelation
}

func (f *Level) HasLevel(customerID uint64, levelIds ...interface{}) (bool, error) {
	err := f.Get(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`level_id`: db.In(levelIds)},
		db.Cond{`status`: modelLevel.LevelStatusActived},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}

func (f *Level) GetByCustomerID(group string, customerID uint64) (result *modelLevel.RelationExt, err error) {
	levelM := dbschema.NewOfficialCustomerLevel(f.Context())
	_, err = levelM.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`group`: group},
		db.Cond{`disabled`: common.BoolN},
	))
	if err != nil {
		return
	}
	rows := levelM.Objects()
	levelIDs := make([]uint, len(rows))
	levelMap := map[uint]int{}
	for i, lv := range rows {
		levelIDs[i] = lv.Id
		levelMap[lv.Id] = i
	}
	err = f.Get(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`level_id`: db.In(levelIDs)},
	))
	if err != nil {
		return
	}
	result = &modelLevel.RelationExt{
		OfficialCustomerLevelRelation: f.OfficialCustomerLevelRelation,
		Level:                         rows[levelMap[f.LevelId]],
	}
	return
}

func (f *Level) Add() (pk interface{}, err error) {
	m := dbschema.NewOfficialCustomerLevelRelation(f.Context())
	err = m.Get(nil, db.And(
		db.Cond{`customer_id`: f.CustomerId},
		db.Cond{`level_id`: f.LevelId},
	))
	if f.AccumulatedDays == 0 && f.Expired > 0 {
		f.AccumulatedDays = (f.Expired - uint(time.Now().Unix())) / 86400
	}
	m.Status = modelLevel.LevelStatusActived
	m.LastRenewalAt = uint(time.Now().Unix())
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		pk, err = f.OfficialCustomerLevelRelation.Insert()
		return
	}
	m.AccumulatedDays += f.AccumulatedDays
	err = f.Edit(nil, `id`, m.Id)
	return
}

func (f *Level) HasPassedDays() (days uint) {
	nowTs := uint(time.Now().Unix())
	days = (nowTs - f.Created) / 86400
	return
}

func (f *Level) Exists(customerID uint64, levelID uint) (bool, error) {
	cond := db.NewCompounds()
	cond.AddKV(`customer_id`, customerID)
	cond.AddKV(`level_id`, levelID)
	return f.OfficialCustomerLevelRelation.Exists(nil, cond.And())
}

func (f *Level) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCustomerLevelRelation.Update(mw, args...)
}

func (f *Level) UpdateAllExpired() error {
	return f.OfficialCustomerLevelRelation.UpdateField(nil, `status`, modelLevel.LevelStatusExpired, db.And(
		db.Cond{`status`: modelLevel.LevelStatusActived},
		db.Cond{`expired`: db.Gt(0)},
		db.Cond{`expired`: db.Lte(time.Now().Unix())},
	))
}

func (f *Level) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.OfficialCustomerLevelRelation.Delete(mw, args...)
	return err
}

func (f *Level) LevelUp() error {
	return modelLevel.LevelUp(f.Context(), f.OfficialCustomerLevelRelation)
}
