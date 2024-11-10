package official

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webfront/dbschema"
)

func NewCollection(ctx echo.Context) *Collection {
	return &Collection{
		OfficialCommonCollection: dbschema.NewOfficialCommonCollection(ctx),
	}
}

type Collection struct {
	*dbschema.OfficialCommonCollection
}

func (f *Collection) Add() (pk interface{}, err error) {
	var exists bool
	exists, err = f.Exists(f.TargetType, f.TargetId, f.CustomerId)
	if err != nil {
		return
	}
	if exists {
		err = f.Context().NewError(code.RepeatOperation, `您已经收藏过了`)
		return
	}
	return f.OfficialCommonCollection.Insert()
}

func (f *Collection) Exists(targetType string, targetID uint64, customerID uint64) (bool, error) {
	return f.OfficialCommonCollection.Exists(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`target_id`: targetID},
		db.Cond{`target_type`: targetType},
	))
}

func (f *Collection) Find(targetType string, targetID uint64, customerID uint64) error {
	return f.Get(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`target_id`: targetID},
		db.Cond{`target_type`: targetType},
	))
}

func (f *Collection) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return f.OfficialCommonCollection.Update(mw, args...)
}

func (f *Collection) DelByTarget(targetType string, targetID uint64) error {
	return f.OfficialCommonCollection.Delete(nil, db.And(
		db.Cond{`target_id`: targetID},
		db.Cond{`target_type`: targetType},
	))
}

func (f *Collection) DelByTargetOwner(targetType string, targetID uint64, customerID uint64) error {
	return f.OfficialCommonCollection.Delete(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`target_id`: targetID},
		db.Cond{`target_type`: targetType},
	))
}

func (f *Collection) ListByTargets(targetType string, targetIDs []uint64, customerID uint64) (map[uint64]*dbschema.OfficialCommonCollection, error) {
	conds := []db.Compound{
		db.Cond{`customer_id`: customerID},
		db.Cond{`target_id`: db.In(targetIDs)},
		db.Cond{`target_type`: targetType},
	}
	_, err := f.ListByOffset(nil, nil, 0, -1, db.And(conds...))
	if err != nil {
		return nil, err
	}
	result := map[uint64]*dbschema.OfficialCommonCollection{}
	for _, v := range f.Objects() {
		result[v.TargetId] = v
	}
	return result, err
}

func (f *Collection) ListPage(targetType string, targetID uint64, customerID uint64, sorts ...interface{}) ([]*CollectionResponse, error) {
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{`customer_id`: customerID},
		db.Cond{`target_id`: targetID},
		db.Cond{`target_type`: targetType},
	)
	list := []*CollectionResponse{}
	err := f.OfficialCommonCollection.ListPageAs(&list, cond, sorts...)
	if err != nil {
		return list, err
	}
	if ls, ok := CollectionTargets[targetType]; ok && ls.List != nil {
		list, err = ls.List.List(f.Context(), list)
	}
	return list, err
}

func (f *Collection) ListPageByOffset(targetType string, targetID uint64, customerID uint64, sorts ...interface{}) ([]*CollectionResponse, error) {
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{`customer_id`: customerID},
		db.Cond{`target_id`: targetID},
		db.Cond{`target_type`: targetType},
	)
	list := []*CollectionResponse{}
	err := f.OfficialCommonCollection.ListPageByOffsetAs(&list, cond, sorts...)
	if err != nil {
		return list, err
	}
	if ls, ok := CollectionTargets[targetType]; ok && ls.List != nil {
		list, err = ls.List.List(f.Context(), list)
	}
	return list, err
}
