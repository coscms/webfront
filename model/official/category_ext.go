package official

import (
	"github.com/coscms/tree"
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/echo/param"
)

type ICategory interface {
	GetCategory1() uint
	GetCategory2() uint
	GetCategory3() uint
	GetCategoryID() uint
	AddCategory(*dbschema.OfficialCommonCategory)
}

type SelectedCategory struct {
	SelectedID uint `json:"selected_id"`
	ParentID   uint `json:"parent_id"`
	Categories []*dbschema.OfficialCommonCategory
}

func NewCategoryTreeRow(c FieldValueGetter) *CategoryTreeRow {
	return &CategoryTreeRow{FieldValueGetter: c}
}

type FieldValueGetter interface {
	GetField(string) interface{}
}
type CategoryTreeRow struct {
	FieldValueGetter
}

func (t *CategoryTreeRow) GetID() int64 {
	return param.AsInt64(t.GetField(`Id`))
}

func (t *CategoryTreeRow) GetParentID() int64 {
	return param.AsInt64(t.GetField(`ParentId`))
}

func (t *CategoryTreeRow) GetObject() interface{} {
	return t.FieldValueGetter
}

func BuildCategoryTree[T FieldValueGetter](list []T) (*tree.Tree, error) {
	rows := make([]tree.Row, len(list))
	for k, v := range list {
		rows[k] = NewCategoryTreeRow(v)
	}
	return tree.Build(rows)
}
