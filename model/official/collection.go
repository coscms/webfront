package official

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/mysql"
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
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
	))
}

func (f *Collection) Find(targetType string, targetID uint64, customerID uint64) error {
	return f.Get(nil, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
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
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
	))
}

func (f *Collection) NewestTimeByCustomerID(targetType string, targetID uint64, customerID uint64) (uint, error) {
	m := dbschema.NewOfficialCommonCollection(f.Context())
	err := m.Get(func(r db.Result) db.Result {
		return r.Select(`created`).OrderBy(`-id`)
	}, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
	))
	if err != nil && err == db.ErrNoMoreRows {
		return 0, nil
	}
	return m.Created, err
}

func (f *Collection) ListByTargets(targetType string, targetIDs []uint64, customerID uint64) (map[uint64]*dbschema.OfficialCommonCollection, error) {
	conds := []db.Compound{
		db.Cond{`customer_id`: customerID},
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: db.In(targetIDs)},
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

func (f *Collection) ListPage(targetType string, customerID uint64, title string, sorts ...interface{}) ([]*CollectionResponse, error) {
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{`customer_id`: customerID},
		db.Cond{`target_type`: targetType},
	)
	if len(title) > 0 {
		cond.Add(mysql.SearchField(`~title`, title))
	}
	list := []*CollectionResponse{}
	err := f.OfficialCommonCollection.ListPageAs(&list, cond, sorts...)
	if err != nil {
		return list, err
	}
	if ls, ok := CollectionTargets[targetType]; ok && ls.HasList() {
		targetIDs := make([]uint64, len(list))
		for index, row := range list {
			targetIDs[index] = row.TargetId
		}
		list, err = ls.List(f.Context(), list, targetIDs)
	}
	return list, err
}

func (f *Collection) ListPageByOffset(targetType string, customerID uint64, title string, sorts ...interface{}) ([]*CollectionResponse, error) {
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{`customer_id`: customerID},
		db.Cond{`target_type`: targetType},
	)
	if len(title) > 0 {
		cond.Add(mysql.SearchField(`~title`, title))
	}
	list := []*CollectionResponse{}
	err := f.OfficialCommonCollection.ListPageByOffsetAs(&list, cond, sorts...)
	if err != nil {
		return list, err
	}
	if ls, ok := CollectionTargets[targetType]; ok && ls.HasList() {
		targetIDs := make([]uint64, len(list))
		for index, row := range list {
			targetIDs[index] = row.TargetId
		}
		list, err = ls.List(f.Context(), list, targetIDs)
	}
	return list, err
}
