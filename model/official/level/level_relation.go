package level

import (
	"time"

	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewRelation(ctx echo.Context) *Relation {
	m := &Relation{
		OfficialCustomerLevelRelation: dbschema.NewOfficialCustomerLevelRelation(ctx),
	}
	return m
}

// 客户的扩展组等级关联关系
type Relation struct {
	*dbschema.OfficialCustomerLevelRelation
}

func (f *Relation) ListByCustomerID(customerID uint64) ([]*dbschema.OfficialCustomerLevelRelation, error) {
	_, err := f.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`customer_id`: customerID},
		db.Cond{`status`: LevelStatusActived},
		db.Or(
			db.Cond{`expired`: 0},
			db.Cond{`expired`: db.Gt(time.Now().Unix())},
		),
	))
	if err != nil {
		return nil, err
	}
	return f.Objects(), nil
}

func (f *Relation) HasGroupLevelByCustomerID(customerID uint64, group string) bool {
	row, err := f.GetGroupLevelByCustomerID(customerID, group)
	return err == nil && row != nil
}

func (f *Relation) GetGroupLevelByCustomerID(customerID uint64, group string) (*dbschema.OfficialCustomerLevel, error) {
	row := dbschema.NewOfficialCustomerLevel(f.Context())
	p := f.NewParam().SetAlias(`r`).AddJoin(`INNER`, row.Name_(), `b`, `b.id=r.level_id`)
	p.SetArgs(db.And(
		db.Cond{`r.customer_id`: customerID},
		db.Cond{`r.status`: LevelStatusActived},
		db.Or(
			db.Cond{`r.expired`: 0},
			db.Cond{`r.expired`: db.Gt(time.Now().Unix())},
		),
		db.Cond{`b.group`: group},
	))
	err := p.SetCols(`b.*`).SetRecv(row).One()
	return row, err
}

func (f *Relation) ListByCustomerIDs(customerIDs []uint64) (map[uint64][]*RelationExt, error) {
	list := []*RelationExt{}
	var mw func(db.Result) db.Result
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{`customer_id`: db.In(customerIDs)},
		db.Cond{`status`: LevelStatusActived},
		db.Or(
			db.Cond{`expired`: 0},
			db.Cond{`expired`: db.Gt(time.Now().Unix())},
		),
	)
	_, err := f.ListByOffset(&list, mw, 0, -1, db.And())
	if err != nil {
		return nil, err
	}
	results := map[uint64][]*RelationExt{}
	for _, row := range list {
		if _, ok := results[row.CustomerId]; !ok {
			results[row.CustomerId] = []*RelationExt{}
		}
		results[row.CustomerId] = append(results[row.CustomerId], row)
	}
	return results, nil
}

func (f *Relation) LevelUp() error {
	return LevelUp(f.Context(), f.OfficialCustomerLevelRelation)
}

func LevelNew(ctx echo.Context, customerID uint64, group string, assetType string, expired uint) error {
	levelM := NewLevel(ctx)
	level, err := levelM.CanAutoLevelUpBy(customerID, group, assetType)
	if err != nil {
		return err
	}
	if level.Id < 1 {
		return nil
	}
	if level.Price > 0 {
		return nil
	}
	f := dbschema.NewOfficialCustomerLevelRelation(ctx)
	f.CustomerId = customerID
	f.LevelId = level.Id
	f.Status = LevelStatusActived
	f.Expired = expired
	_, err = f.Insert()
	return err
}

func LevelUp(ctx echo.Context, f *dbschema.OfficialCustomerLevelRelation) error {
	currentLevelM := dbschema.NewOfficialCustomerLevel(ctx)
	// 当前等级
	err := currentLevelM.Get(func(r db.Result) db.Result {
		return r.Select(`score`, `group`, `integral_asset`)
	}, db.And(
		db.Cond{`id`: f.LevelId},
		db.Cond{`disabled`: `N`},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return f.Delete(nil, `id`, f.Id)
	}
	// 当前积分可以匹配的等级
	levelM := NewLevel(ctx)
	var level *dbschema.OfficialCustomerLevel
	level, err = levelM.CanAutoLevelUpBy(f.CustomerId, currentLevelM.Group, currentLevelM.IntegralAsset)
	if err != nil {
		return err
	}
	if level.Id < 1 || level.Id == f.LevelId {
		return nil
	}
	if level.Price > 0 {
		return nil
	}
	if currentLevelM.Score == level.Score {
		return nil
	}
	f.LevelId = level.Id
	return f.UpdateField(nil, `level_id`, level.Id, `id`, f.Id)
}
