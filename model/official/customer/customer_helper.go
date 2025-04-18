package customer

import (
	"errors"
	"time"

	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/top"
	modelLevel "github.com/coscms/webfront/model/official/level"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

// ErrSameDay 同一天不能领取两次奖励
var ErrSameDay = errors.New("same day")

// AddRewardOnSignIn 登录时奖励
func (f *Customer) AddRewardOnSignIn(amount float64) error {
	if amount > 0 {
		err := f.IsSameDay(`experience`, `sign_in`)
		if err != nil {
			if err == ErrSameDay {
				return nil
			}
			return err
		}
	}
	return f.AddExperience(amount, f.Id, `sign_in`, `official_customer`, f.Id, `登录奖励`)
}

// AddRewardOnSignUp 注册时奖励
func (f *Customer) AddRewardOnSignUp(amount float64) error {
	return f.AddExperience(amount, f.Id, `sign_up`, `official_customer`, f.Id, `注册奖励`)
}

// IsSameDay 是否为同一天领取第二次
func (f *Customer) IsSameDay(assetType string, sourceType string) error {
	m := dbschema.NewOfficialCustomerWalletFlow(f.Context())
	m.CPAFrom(f.OfficialCustomer)
	err := m.Get(func(r db.Result) db.Result {
		return r.Select(`id`, `created`).OrderBy(`-id`)
	}, db.And(
		db.Cond{`customer_id`: f.Id},
		db.Cond{`asset_type`: assetType},
		db.Cond{`amount_type`: `balance`},
		db.Cond{`source_type`: sourceType},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
	} else {
		// 每天只奖励一次，避免刷积分作弊
		if top.IsSameDay(time.Unix(int64(m.Created), 0)) {
			return ErrSameDay
		}
	}
	return nil
}

func (f *Customer) ExistsAssetSource(assetType string, sourceType string, sourceTable string, sourceID uint64) (bool, error) {
	m := dbschema.NewOfficialCustomerWalletFlow(f.Context())
	m.CPAFrom(f.OfficialCustomer)
	return m.Exists(nil, db.And(
		db.Cond{`customer_id`: f.Id},
		db.Cond{`asset_type`: assetType},
		db.Cond{`amount_type`: `balance`},
		db.Cond{`source_type`: sourceType},
		db.Cond{`source_table`: sourceTable},
		db.Cond{`source_id`: sourceID},
	))
}

// AddExperience 增加经验值
func (f *Customer) AddExperience(amount float64, customerID uint64, sourceType string, sourceTable string, sourceID uint64, description string) error {
	if amount <= 0 {
		return nil
	}
	err := f.Context().Begin()
	if err != nil {
		return err
	}
	m := NewWallet(f.Context())
	m.Flow.CustomerId = customerID
	m.Flow.AssetType = AssetTypeExperience
	m.Flow.AmountType = AmountTypeBalance
	m.Flow.Amount = amount
	m.Flow.SourceType = sourceType
	m.Flow.SourceTable = sourceTable
	m.Flow.SourceId = sourceID
	m.Flow.TradeNo = ``
	m.Flow.Status = FlowStatusConfirmed //状态(pending-待确认;confirmed-已确认;canceled-已取消)
	m.Flow.Description = description
	err = m.AddRepeatableFlow()
	if err != nil {
		f.Context().Rollback()
		return err
	}
	return f.Context().Commit()
}

// AddIntegral 增加消费积分
func (f *Customer) AddIntegral(amount float64, customerID uint64, sourceType string, sourceTable string, sourceID uint64, description string) error {
	if amount <= 0 {
		return nil
	}
	err := f.Context().Begin()
	if err != nil {
		return err
	}
	m := NewWallet(f.Context())
	m.Flow.CustomerId = customerID
	m.Flow.AssetType = AssetTypeIntegral
	m.Flow.AmountType = AmountTypeBalance
	m.Flow.Amount = amount
	m.Flow.SourceType = sourceType
	m.Flow.SourceTable = sourceTable
	m.Flow.SourceId = sourceID
	m.Flow.TradeNo = ``
	m.Flow.Status = FlowStatusConfirmed //状态(pending-待确认;confirmed-已确认;canceled-已取消)
	m.Flow.Description = description
	err = m.AddRepeatableFlow()
	if err != nil {
		f.Context().Rollback()
		return err
	}
	return f.Context().Commit()
}

// LevelUpOnSignIn 登录时检查是否可升级
func (f *Customer) LevelUpOnSignIn(set echo.H) error {
	levelM := modelLevel.NewLevel(f.Context())

	// 当前积分可以匹配的等级
	level, err := levelM.CanAutoLevelUpByCustomerID(f.Id)
	if err != nil {
		return err
	}
	if level.Id < 1 || level.Id == f.LevelId {
		return nil
	}
	if f.LevelId > 0 {
		currentLevelM := dbschema.NewOfficialCustomerLevel(f.Context())
		// 当前等级
		err = currentLevelM.Get(func(r db.Result) db.Result {
			return r.Select(`score`)
		}, db.And(
			db.Cond{`id`: f.LevelId},
			db.Cond{`disabled`: `N`},
			db.Cond{`group`: `base`},
		))
		if err == nil {
			if currentLevelM.Score != level.Score {
				// currentLevelM.Score < level.Score 可以升级
				// currentLevelM.Score > level.Score 需要降级
				f.LevelId = level.Id
				set.Set(`level_id`, f.LevelId)
			}
		}
	} else {
		f.LevelId = level.Id
		set.Set(`level_id`, f.LevelId)
	}
	return nil
}

func (f *Customer) GetAdmins(customerIDs ...uint64) (map[uint64]struct{}, error) {
	r := map[uint64]struct{}{}
	_, err := f.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`)
	}, 0, -1, db.And(
		db.Cond{`id`: db.In(customerIDs)},
		db.Cond{`uid`: db.Gt(0)},
	))
	if err == nil {
		for _, v := range f.Objects() {
			r[v.Id] = struct{}{}
		}
	}
	return r, err
}
