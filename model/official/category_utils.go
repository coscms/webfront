package official

import (
	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
)

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

// SortCategoryByParent sorts the categories in the given list by their parent_id.
// It will build a tree of categories, and then traverse the tree to get the sorted list.
// The root node of the tree is the category with parent_id of 0.
// The function will return a list of categories sorted by their parent_id, sort and id.
func SortCategoryByParent[idT com.Number, vT factory.Model](list []vT) []vT {
	mp := map[idT][]vT{} // {parent_id:[]}
	for _, row := range list {
		pid := row.GetField(`ParentId`).(idT)
		if _, ok := mp[pid]; !ok {
			mp[pid] = []vT{}
		}
		mp[pid] = append(mp[pid], row)
	}
	rows := make([]vT, 0, len(list))
	var appendFn func(children []vT)
	appendFn = func(children []vT) {
		for _, row := range children {
			id := row.GetField(`Id`).(idT)
			rows = append(rows, row)
			appendFn(mp[id])
		}
	}
	appendFn(mp[0])
	return rows
}

// ListAllParentByType returns a list of categories that match the given conditions.
// The conditions are: type = typ, disabled = N, level <= max_level.
// If exclude_id is not zero, it will also exclude the category with id equal to exclude_id.
// If extraConds is not empty, it will also add the conditions to the query.
// If max_level is zero, the categories will be sorted by sort and id.
// Otherwise, the categories will be sorted by level, parent_id, sort and id.
// The function will return a list of categories sorted by their parent_id, sort and id.
// The root node of the tree is the category with parent_id of 0.
func ListAllParentByType[idT com.Number, T factory.Model](f T, objects func() []T, typ string, excludeId idT, maxLevel uint, extraConds ...db.Compound) []T {
	return ListAllParentBy(f, objects, `type`, typ, excludeId, maxLevel, extraConds...)
}

// ListAllParentBy returns a list of categories that match the given conditions.
// The conditions are: column_name = column_value, disabled = N, level <= max_level.
// If exclude_id is not zero, it will also exclude the category with id equal to exclude_id.
// If extraConds is not empty, it will also add the conditions to the query.
// If max_level is zero, the categories will be sorted by sort and id.
// Otherwise, the categories will be sorted by level, parent_id, sort and id.
// The function will return a list of categories sorted by their parent_id, sort and id.
// The root node of the tree is the category with parent_id of 0.
func ListAllParentBy[idT com.Number, T factory.Model](f T, objects func() []T, columnName string, columnValue interface{}, excludeId idT, maxLevel uint, extraConds ...db.Compound) []T {
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
	if len(columnName) > 0 && columnValue != nil {
		cond.AddKV(columnName, columnValue)
	}
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
		return SortCategoryByParent[idT](objects())
	}
	return objects()
}
