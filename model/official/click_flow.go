package official

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webfront/dbschema"
)

func NewClickFlow(ctx echo.Context) *ClickFlow {
	return &ClickFlow{
		OfficialCommonClickFlow: dbschema.NewOfficialCommonClickFlow(ctx),
	}
}

type ClickFlow struct {
	*dbschema.OfficialCommonClickFlow
}

func (f *ClickFlow) Add() (pk interface{}, err error) {
	old := dbschema.NewOfficialCommonClickFlow(f.Context())
	err = old.Get(func(r db.Result) db.Result {
		return r.Select(`type`)
	}, db.And(
		db.Cond{`target_type`: f.TargetType},
		db.Cond{`target_id`: f.TargetId},
		db.Cond{`owner_id`: f.OwnerId},
		db.Cond{`owner_type`: f.OwnerType},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		err = nil
	} else {
		if old.Type == f.Type {
			err = f.Context().NewError(code.RepeatOperation, `您已经表过态了`)
		} else {
			err = f.Context().NewError(code.DataAlreadyExists, `您已经表态过其它观点`)
		}
		return
	}
	return f.OfficialCommonClickFlow.Insert()
}

func (f *ClickFlow) Exists(targetType string, targetID uint64, ownerID uint64, ownerType string) (bool, error) {
	return f.OfficialCommonClickFlow.Exists(nil, db.And(
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
		db.Cond{`owner_id`: ownerID},
		db.Cond{`owner_type`: ownerType},
	))
}

func (f *ClickFlow) Find(targetType string, targetID uint64, ownerID uint64, ownerType string) error {
	return f.Get(nil, db.And(
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
		db.Cond{`owner_id`: ownerID},
		db.Cond{`owner_type`: ownerType},
	))
}

func (f *ClickFlow) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCommonClickFlow.Update(mw, args...)
}

func (f *ClickFlow) DelByTarget(targetType string, targetID uint64) error {
	return f.OfficialCommonClickFlow.Delete(nil, db.And(
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
	))
}

func (f *ClickFlow) DelByTargetOwner(targetType string, targetID uint64, ownerID uint64, ownerType string) error {
	return f.OfficialCommonClickFlow.Delete(nil, db.And(
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
		db.Cond{`owner_id`: ownerID},
		db.Cond{`owner_type`: ownerType},
	))
}

func (f *ClickFlow) NewestTimeByCustomerID(targetType string, targetID uint64, ownerID uint64, ownerType string) (uint, error) {
	m := dbschema.NewOfficialCommonClickFlow(f.Context())
	err := m.Get(func(r db.Result) db.Result {
		return r.Select(`created`).OrderBy(`-id`)
	}, db.And(
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
		db.Cond{`owner_id`: ownerID},
		db.Cond{`owner_type`: ownerType},
	))
	if err != nil && err == db.ErrNoMoreRows {
		return 0, nil
	}
	return m.Created, err
}

func (f *ClickFlow) ListByCustomerTargets(targetType string, targetIDs []uint64, customerID uint64) (map[uint64]*dbschema.OfficialCommonClickFlow, error) {
	return f.ListByTargets(targetType, targetIDs, customerID, 0)
}

func (f *ClickFlow) ListByAdminTargets(targetType string, targetIDs []uint64, adminUID uint) (map[uint64]*dbschema.OfficialCommonClickFlow, error) {
	return f.ListByTargets(targetType, targetIDs, 0, adminUID)
}

func (f *ClickFlow) ListByTargets(targetType string, targetIDs []uint64, customerID uint64, adminUID uint) (map[uint64]*dbschema.OfficialCommonClickFlow, error) {
	conds := []db.Compound{
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: db.In(targetIDs)},
	}
	if customerID > 0 {
		conds = append(conds, db.Cond{`owner_id`: customerID})
		conds = append(conds, db.Cond{`owner_type`: `customer`})
	} else {
		if adminUID == 0 {
			return map[uint64]*dbschema.OfficialCommonClickFlow{}, nil
		}
		conds = append(conds, db.Cond{`owner_id`: adminUID})
		conds = append(conds, db.Cond{`owner_type`: `user`})
	}
	_, err := f.ListByOffset(nil, nil, 0, -1, db.And(conds...))
	if err != nil {
		return nil, err
	}
	result := map[uint64]*dbschema.OfficialCommonClickFlow{}
	for _, v := range f.Objects() {
		result[v.TargetId] = v
	}
	return result, err
}
