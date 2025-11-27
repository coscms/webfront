package article

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webfront/dbschema"
	modelAuthor "github.com/coscms/webfront/model/author"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/coscms/webfront/model/official"
)

var Contype = official.Contype

type ArticleWithOwner struct {
	*dbschema.OfficialCommonArticle
	User     *modelAuthor.User     `db:"-,relation=id:owner_id|gtZero|eq(owner_type:user),columns=id&username&avatar" json:",omitempty"`
	Customer *modelAuthor.Customer `db:"-,relation=id:owner_id|gtZero|eq(owner_type:customer),columns=id&name&avatar" json:",omitempty"`
	Category *Category             `db:"-,relation=id:category_id|gtZero,columns=id&name" json:",omitempty"`
}

func MultilingualArticlesWithOwner(ctx echo.Context, list []*ArticleWithOwner) {
	if !i18nm.IsMultilingual() {
		return
	}
	i18nm.GetModelsTranslations(ctx, list)
	cateIDs := map[uint][]int{}
	categories := []*Category{}
	for i, a := range list {
		if a.Category != nil {
			if _, ok := cateIDs[a.Category.Id]; !ok {
				cateIDs[a.Category.Id] = []int{}
				categories = append(categories, a.Category)
			}
			cateIDs[a.Category.Id] = append(cateIDs[a.Category.Id], i)
		}
	}
	if len(categories) == 0 {
		return
	}
	i18nm.GetModelsTranslations(ctx, categories, `name`)
	for _, v := range categories {
		for _, i := range cateIDs[v.Id] {
			list[i].Category.Name = v.Name
		}
	}
}

type Category struct {
	Id   uint   `db:"id"`
	Name string `db:"name"`
}

func (c *Category) Short_() string {
	return `official_common_category`
}

func (c *Category) Name_() string {
	return dbschema.WithPrefix(`official_common_category`)
}

func (a *Category) GetField(field string) interface{} {
	switch field {
	case "Id":
		return a.Id
	case "Name":
		return a.Name
	default:
		return nil
	}
}

func (a *Category) FromRow(row map[string]interface{}) {
	for key, value := range row {
		if _, ok := value.(db.RawValue); ok {
			continue
		}
		switch key {
		case "id":
			a.Id = param.AsUint(value)
		case "name":
			a.Name = param.AsString(value)
		}
	}
}

type ArticleAndSourceInfo struct {
	*dbschema.OfficialCommonArticle
	Categories []*dbschema.OfficialCommonCategory
	SourceInfo echo.KV `db:"-"`
}

func (a *ArticleAndSourceInfo) GetCategory1() uint {
	return a.Category1
}

func (a *ArticleAndSourceInfo) GetCategory2() uint {
	return a.Category2
}

func (a *ArticleAndSourceInfo) GetCategory3() uint {
	return a.Category3
}

func (a *ArticleAndSourceInfo) GetCategoryID() uint {
	return a.CategoryId
}

func (a *ArticleAndSourceInfo) AddCategory(g *dbschema.OfficialCommonCategory) {
	a.Categories = append(a.Categories, g)
}

func WithSourceInfo(ctx echo.Context, list []*ArticleAndSourceInfo) error {
	sourceTableIds := map[string][]string{}
	sourceIDIndexes := map[string][]int{}
	for i, v := range list {
		if len(v.SourceTable) == 0 || len(v.SourceId) == 0 {
			continue
		}
		if _, ok := sourceTableIds[v.SourceTable]; !ok {
			sourceTableIds[v.SourceTable] = []string{}
		}
		sourceTableIds[v.SourceTable] = append(sourceTableIds[v.SourceTable], v.SourceId)
		scKey := v.SourceTable + `:` + v.SourceId
		if _, ok := sourceIDIndexes[scKey]; !ok {
			sourceIDIndexes[scKey] = []int{}
		}
		sourceIDIndexes[scKey] = append(sourceIDIndexes[scKey], i)
	}
	if len(sourceTableIds) > 0 {
		for sourceTable, sourceIds := range sourceTableIds {
			infoMapGetter := Source.GetInfoMapGetter(sourceTable)
			if infoMapGetter == nil {
				continue
			}
			infoMap, err := infoMapGetter(ctx, sourceIds...)
			if err != nil {
				return err
			}
			if infoMap == nil {
				continue
			}
			for sourceID, info := range infoMap {
				keys, ok := sourceIDIndexes[sourceTable+`:`+sourceID]
				if !ok {
					continue
				}
				for _, index := range keys {
					list[index].SourceInfo = info
				}
			}
		}
	}
	return nil
}
