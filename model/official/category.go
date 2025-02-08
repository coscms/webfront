package official

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/top"
)

var (
	// CategoryMaxLevel 分类最大层数(层从0开始)
	CategoryMaxLevel uint = 3 //3代表最大允许4层
	// CategoryTypes 分类的类别
	CategoryTypes = echo.NewKVData()
)

func init() {
	CategoryTypes.AddItem(&echo.KV{K: `article`, V: `文章`})
	CategoryTypes.AddItem(&echo.KV{K: `friendlink`, V: `友情链接`})
}

// AddCategoryType 登记新的类别
func AddCategoryType(value, label string) {
	CategoryTypes.Add(value, label)
}

// ChineseSpace 中文全角空白字符
const ChineseSpace = `　`

func NewCategory(ctx echo.Context) *Category {
	m := &Category{
		OfficialCommonCategory: dbschema.NewOfficialCommonCategory(ctx),
		maxLevel:               CategoryMaxLevel,
	}
	return m
}

type Category struct {
	*dbschema.OfficialCommonCategory
	maxLevel uint
}

func (f *Category) MaxLevel() uint {
	return f.maxLevel
}

func (f *Category) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.Get(mw, args...)
	if err != nil {
		return err
	}
	if f.HasChild == common.BoolY {
		return f.Context().NewError(code.Failure, `删除失败：请先删除子分类`)
	}
	if err = f.Context().Begin(); err != nil {
		return err
	}
	defer func() {
		f.Context().End(err == nil)
	}()
	err = f.OfficialCommonCategory.Delete(mw, args...)
	if err != nil {
		return err
	}
	err = f.UpdateAllParents(f.OfficialCommonCategory)
	return err
}

func (f *Category) ListAllParent(typ string, excludeId uint, maxLevels ...uint) []*dbschema.OfficialCommonCategory {
	maxLevel := f.MaxLevel()
	if len(maxLevels) > 0 {
		maxLevel = maxLevels[0]
	}
	return f.ListAllParentBy(typ, excludeId, maxLevel)
}

func (f *Category) ListAllParentBy(typ string, excludeId uint, maxLevel uint, extraConds ...db.Compound) []*dbschema.OfficialCommonCategory {
	var queryMW func(r db.Result) db.Result
	if maxLevel == 0 {
		queryMW = func(r db.Result) db.Result {
			return r.OrderBy(`sort`, `id`)
		}
	} else {
		queryMW = func(r db.Result) db.Result {
			return r.OrderBy(`level`, `parent_id`, `sort`, `id`)
		}
	}
	cond := db.NewCompounds()
	cond.AddKV(`type`, typ)
	cond.AddKV(`disabled`, common.BoolN)
	cond.AddKV(`level`, db.Lte(maxLevel))
	if excludeId > 0 {
		cond.AddKV(`id`, db.NotEq(excludeId))
	}
	if len(extraConds) > 0 {
		cond.Add(extraConds...)
	}
	f.ListByOffset(nil, queryMW, 0, -1, cond.And())
	if maxLevel > 0 {
		return SortCategoryByParent(f.Objects())
	}
	return f.Objects()
}

func (f *Category) ListByParentID(typ string, parentID uint, extraConds ...db.Compound) []*dbschema.OfficialCommonCategory {
	queryMW := func(r db.Result) db.Result {
		return r.OrderBy(`sort`, `id`)
	}
	cond := db.NewCompounds()
	cond.AddKV(`type`, typ)
	cond.AddKV(`disabled`, common.BoolN)
	cond.AddKV(`parent_id`, parentID)
	if len(extraConds) > 0 {
		cond.Add(extraConds...)
	}
	f.ListByOffset(nil, queryMW, 0, -1, cond.And())
	return f.Objects()
}

func (f *Category) ListForSelected(typ string, id uint, extraConds ...db.Compound) []*SelectedCategory {
	if id == 0 {
		return []*SelectedCategory{
			{
				SelectedID: 0,
				ParentID:   0,
				Categories: f.ListByParentID(typ, 0, extraConds...),
			},
		}
	}
	posIds, err := f.PositionIds(id)
	if err == nil {
		return nil
	}
	categories := make([]*SelectedCategory, 0, len(posIds)+1)
	var parentID uint
	for _, catID := range posIds {
		sc := &SelectedCategory{
			SelectedID: catID,
			ParentID:   parentID,
			Categories: f.ListByParentID(typ, parentID, extraConds...),
		}
		categories = append(categories, sc)
		parentID = catID
	}
	if parentID > 0 {
		sc := &SelectedCategory{
			SelectedID: 0,
			ParentID:   parentID,
			Categories: f.ListByParentID(typ, parentID, extraConds...),
		}
		categories = append(categories, sc)
	}
	return categories
}

func SortCategoryByParent(list []*dbschema.OfficialCommonCategory) []*dbschema.OfficialCommonCategory {
	mp := map[uint][]*dbschema.OfficialCommonCategory{} // {parent_id:[]}
	for _, row := range list {
		if _, ok := mp[row.ParentId]; !ok {
			mp[row.ParentId] = []*dbschema.OfficialCommonCategory{}
		}
		mp[row.ParentId] = append(mp[row.ParentId], row)
	}
	rows := make([]*dbschema.OfficialCommonCategory, 0, len(list))
	var appendFn func(children []*dbschema.OfficialCommonCategory)
	appendFn = func(children []*dbschema.OfficialCommonCategory) {
		for _, row := range children {
			rows = append(rows, row)
			appendFn(mp[row.Id])
		}
	}
	appendFn(mp[0])
	return rows
}

func (f *Category) ListIndent(categoryList []*dbschema.OfficialCommonCategory) []*dbschema.OfficialCommonCategory {
	for idx, row := range categoryList {
		categoryList[idx].Name = strings.Repeat(ChineseSpace, int(row.Level)) + row.Name
	}
	return categoryList
}

func (f *Category) Parents(parentID uint, onlyID ...bool) ([]dbschema.OfficialCommonCategory, error) {
	categories := []dbschema.OfficialCommonCategory{}
	r := map[uint]bool{}
	var _mw func(db.Result) db.Result
	if len(onlyID) > 0 && onlyID[0] {
		_mw = func(r db.Result) db.Result {
			return r.Select(`id`, `parent_id`)
		}
	}
	for parentID > 0 && !r[parentID] {
		err := f.Get(_mw, `id`, parentID)
		if err != nil {
			if err == db.ErrNoMoreRows {
				break
			}
			return categories, err
		}
		categories = append(categories, *f.OfficialCommonCategory)
		r[parentID] = true
		parentID = f.ParentId
	}
	return categories, nil
}

func (f *Category) Positions(id uint) ([]dbschema.OfficialCommonCategory, error) {
	parents, err := f.Parents(id)
	if err != nil {
		return parents, err
	}
	if len(parents) == 0 {
		return parents, nil
	}
	positions := make([]dbschema.OfficialCommonCategory, len(parents))
	var index int
	for end := len(parents) - 1; end >= 0; end-- {
		positions[index] = parents[end]
		index++
	}
	return positions, err
}

func (f *Category) PositionIds(id uint) ([]uint, error) {
	parents, err := f.Parents(id, true)
	if err != nil {
		return nil, err
	}
	if len(parents) == 0 {
		return nil, nil
	}
	result := make([]uint, len(parents))
	for k, i := 0, len(parents)-1; i >= 0; i-- {
		result[k] = parents[i].Id
		k++
	}
	return result, nil
}

func (f *Category) setDefaults() {
	f.Name = strings.TrimSpace(f.Name)
	f.Keywords = strings.TrimSpace(f.Keywords)
	f.Description = strings.TrimSpace(f.Description)
	f.Template = strings.TrimSpace(f.Template)
	f.Disabled = common.GetBoolFlag(f.Disabled)
	f.ShowOnMenu = common.GetBoolFlag(f.ShowOnMenu, common.BoolY)
	f.Slugify = strings.TrimSpace(f.Slugify)
	if len(f.Slugify) == 0 {
		f.Slugify = top.Slugify(f.Name)
	}
}

func (f *Category) Add() (pk interface{}, err error) {
	f.Context().Begin()
	defer func() {
		f.Context().End(err == nil)
	}()
	f.setDefaults()
	err = f.Exists(f.Name)
	if err != nil {
		return
	}
	if f.ParentId > 0 {
		parent := dbschema.NewOfficialCommonCategory(f.Context())
		err = parent.Get(nil, `id`, f.ParentId)
		if err != nil {
			if err != db.ErrNoMoreRows {
				return
			}
			err = f.Context().NewError(code.DataNotFound, `父级分类不存在`)
			return
		}
		f.Level = parent.Level + 1
		f.Type = parent.Type
		if parent.HasChild == common.BoolN {
			err = parent.UpdateField(nil, `has_child`, common.BoolY, `id`, f.ParentId)
			if err != nil {
				return
			}
		}
	} else {
		f.Level = 0
	}
	if f.Level > f.MaxLevel() {
		err = f.Context().NewError(code.Failure, `操作失败！分类超过最大层数: %d`, f.MaxLevel())
		return
	}
	return f.OfficialCommonCategory.Insert()
}

func (f *Category) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	f.Context().Begin()
	defer func() {
		f.Context().End(err == nil)
	}()
	f.setDefaults()
	if err = f.ExistsOther(f.Name, f.Id); err != nil {
		return err
	}
	oldData := dbschema.NewOfficialCommonCategory(f.Context())
	err = oldData.Get(nil, args...)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return f.Context().NewError(code.DataNotFound, `分类不存在`)
	}
	if oldData.Id == f.ParentId {
		return f.Context().NewError(code.Failure, `不能选择自己为上级分类`)
	}
	if oldData.ParentId != f.ParentId {
		if f.ParentId > 0 {
			parent := dbschema.NewOfficialCommonCategory(f.Context())
			err = parent.Get(nil, `id`, f.ParentId)
			if err != nil {
				if err != db.ErrNoMoreRows {
					return err
				}
				return f.Context().NewError(code.DataNotFound, `父级分类不存在`)
			}
			f.Level = parent.Level + 1
			f.Type = parent.Type
			if parent.HasChild == common.BoolN {
				err = parent.UpdateField(nil, `has_child`, common.BoolY, `id`, f.ParentId)
				if err != nil {
					return
				}
			}
		} else {
			f.Level = 0
		}
		if f.Level > f.MaxLevel() {
			return f.Context().NewError(code.Failure, `操作失败！分类超过最大层数: %d`, f.MaxLevel())
		}
		err = f.UpdateAllParents(oldData)
		if err != nil {
			return err
		}
		err = f.UpdateAllChildren(f.OfficialCommonCategory)
		if err != nil {
			return err
		}
	}
	return f.OfficialCommonCategory.Update(mw, args...)
}

// UpdateAllParents 更新所有父级分类(目前仅仅用于更新父级has_child值)
func (f *Category) UpdateAllParents(oldData *dbschema.OfficialCommonCategory) (err error) {
	if oldData.ParentId == 0 {
		return
	}
	oldParent := dbschema.NewOfficialCommonCategory(f.Context())
	err = oldParent.Get(nil, `id`, oldData.ParentId)
	if err == nil {
		ocond := db.And(
			db.Cond{`parent_id`: oldData.ParentId},
			db.Cond{`id`: db.NotEq(oldData.Id)},
		)
		var n int64
		n, err = oldData.Count(nil, ocond)
		if err != nil {
			return
		}
		if n == 0 && oldParent.HasChild != common.BoolN {
			err = oldData.UpdateField(nil, `has_child`, common.BoolN, `id`, oldData.ParentId)
		} else if n > 0 && oldParent.HasChild == common.BoolN {
			err = oldData.UpdateField(nil, `has_child`, common.BoolY, `id`, oldData.ParentId)
		}
	}
	if err != nil && err == db.ErrNoMoreRows {
		err = nil
	}
	return
}

// UpdateAllChildren 更新所有子孙分类level值
func (f *Category) UpdateAllChildren(row *dbschema.OfficialCommonCategory) error {
	children := dbschema.NewOfficialCommonCategory(f.Context())
	_, err := children.ListByOffset(nil, nil, 0, -1, `parent_id`, row.Id)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return nil
	}
	for _, child := range children.Objects() {
		child.Level = row.Level + 1
		err = child.UpdateField(nil, `level`, child.Level, `id`, child.Id)
		if err != nil {
			return err
		}
		err = f.UpdateAllChildren(child)
		if err != nil {
			return err
		}
	}
	return err
}

func (f *Category) Exists(name string) error {
	exists, err := f.OfficialCommonCategory.Exists(nil, db.Cond{`name`: name})
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().NewError(code.DataAlreadyExists, `分类名称“%s”已经使用过了`, name)
	}
	return err
}

func (f *Category) ExistsOther(name string, id uint) error {
	exists, err := f.OfficialCommonCategory.Exists(nil, db.And(
		db.Cond{`name`: name},
		db.Cond{`id`: db.NotEq(id)},
	))
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().NewError(code.DataAlreadyExists, `分类名称“%s”已经使用过了`, name)
	}
	return err
}

// ListChildren 查询子分类
func (f *Category) ListChildren(parentID uint) ([]*dbschema.OfficialCommonCategory, error) {
	_, err := f.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`parent_id`: parentID},
		db.Cond{`disabled`: common.BoolN},
	))
	if err != nil {
		return nil, err
	}
	return f.Objects(), nil
}

func (f *Category) FillTo(tg []ICategory) error {
	var categoryIds []uint
	for _, u := range tg {
		if u.GetCategory1() > 0 {
			if !com.InUintSlice(u.GetCategory1(), categoryIds) {
				categoryIds = append(categoryIds, u.GetCategory1())
			}
		}
		if u.GetCategory2() > 0 {
			if !com.InUintSlice(u.GetCategory2(), categoryIds) {
				categoryIds = append(categoryIds, u.GetCategory2())
			}
		}
		if u.GetCategory3() > 0 {
			if !com.InUintSlice(u.GetCategory3(), categoryIds) {
				categoryIds = append(categoryIds, u.GetCategory3())
			}
		}
		if u.GetCategoryID() > 0 {
			if !com.InUintSlice(u.GetCategoryID(), categoryIds) {
				categoryIds = append(categoryIds, u.GetCategoryID())
			}
		}
	}
	_, err := f.ListByOffset(nil, nil, 0, -1, db.Cond{`id IN`: categoryIds})
	if err != nil {
		return err
	}
	categoryList := f.Objects()
	categoryMap := map[uint]*dbschema.OfficialCommonCategory{}
	for _, g := range categoryList {
		categoryMap[g.Id] = g
	}
	for _, v := range tg {
		if v.GetCategory1() > 0 {
			if g, y := categoryMap[v.GetCategory1()]; y {
				v.AddCategory(g)
			}
		}
		if v.GetCategory2() > 0 {
			if g, y := categoryMap[v.GetCategory2()]; y {
				v.AddCategory(g)
			}
		}
		if v.GetCategory3() > 0 {
			if g, y := categoryMap[v.GetCategory3()]; y {
				v.AddCategory(g)
			}
		}
		if v.GetCategoryID() > 0 && (v.GetCategoryID() != v.GetCategory1() && v.GetCategoryID() != v.GetCategory2() && v.GetCategoryID() != v.GetCategory3()) {
			if g, y := categoryMap[v.GetCategoryID()]; y {
				v.AddCategory(g)
			}
		}
	}
	return nil
}

func CollectionCategoryIDs() func(...uint) []uint {
	categoryIds := []uint{}
	return func(cIds ...uint) []uint {
		for _, cID := range cIds {
			if !com.InUintSlice(cID, categoryIds) {
				categoryIds = append(categoryIds, cID)
			}
		}
		return categoryIds
	}
}
