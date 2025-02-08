package official

import "github.com/coscms/webfront/dbschema"

type ICategory interface {
	GetCategory1() uint
	GetCategory2() uint
	GetCategory3() uint
	GetCategoryID() uint
	AddCategory(*dbschema.OfficialCommonCategory)
}

type SelectedCategory struct {
	SelectedID uint
	ParentID   uint
	Categories []*dbschema.OfficialCommonCategory
}
