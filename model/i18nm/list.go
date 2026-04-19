package i18nm

import (
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"
)

var TableTitles = echo.NewKVxData[[]string, any]()

type ListItem struct {
	*dbschema.OfficialI18nTranslation
	*dbschema.OfficialI18nResource
}

type ListQuery struct {
	Table string
	RowID uint64
}

func List(ctx echo.Context, query ListQuery) ([]*ListItem, error) {
	var list []*ListItem
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	rM := dbschema.NewOfficialI18nResource(ctx)
	cond := db.NewCompounds()
	if len(query.Table) > 0 {
		cond.AddKV(`R.code`, db.Like(query.Table+`.%`))
	}
	if query.RowID > 0 {
		cond.AddKV(`T.row_id`, query.RowID)
	}
	pr := tM.NewParam().SetAlias(`T`).SetArgs(cond.And()).AddJoin(`INNER`, rM.Name_(), `R`, `R.id=T.resource_id`).SetContext(ctx)
	tM.SetParam(pr)
	err := tM.ListPageByOffsetAs(&list, cond)
	if err != nil {
		return list, err
	}
	return list, err
}

func UpdateColumnTranslation(ctx echo.Context, table string, column string, rowID uint64, lang string, translatedText string) error {
	resourceIDs, _, err := getResourceIDs(ctx, table, column)
	if err != nil {
		return err
	}
	if len(resourceIDs) == 0 {
		return nil
	}

	tM := dbschema.NewOfficialI18nTranslation(ctx)
	cnd := db.NewCompounds()
	cnd.AddKV(`lang`, lang)
	cnd.AddKV(`row_id`, rowID)
	cnd.AddKV(`resource_id`, resourceIDs[0])
	var affected int64
	affected, err = tM.UpdatexFields(nil, echo.H{
		`text`: translatedText,
	}, cnd.And())
	if err != nil {
		return err
	}
	if affected < 1 {
		tM.Lang = lang
		tM.ResourceId = resourceIDs[0]
		tM.RowId = rowID
		tM.Text = translatedText
		_, err = tM.Insert()
	}
	return err
}

func ListByResource(ctx echo.Context, table string, sorts ...any) ([]echo.H, error) {
	var list []echo.H
	resourceIDs, resourceFields, err := getResourceIDs(ctx, table)
	if err != nil {
		return list, err
	}
	cnd := db.NewCompounds()
	if len(sorts) == 0 {
		sorts = append(sorts, `-id`)
	}
	columns := make([]any, 0, len(resourceFields)+1)
	columns = append(columns, `id`)
	for _, column := range resourceFields {
		columns = append(columns, column)
	}
	smw := func(r db.Result) db.Result {
		return r.Select(columns...).OrderBy(sorts...)
	}
	pr := factory.ParamPoolGet().SetContext(ctx)
	ls := pr.SetCollection(table).SetRecv(&list).NewLister()
	err = pagination.ListPageByOffset(ls, cnd, smw)
	pr.Release()
	if err != nil {
		return list, err
	}
	if len(list) == 0 {
		return list, err
	}
	rowIDs := make([]uint64, len(list))
	idxByRowID := map[uint64]int{}
	for idx, row := range list {
		rowID := row.Uint64(`id`)
		rowIDs[idx] = rowID
		idxByRowID[rowID] = idx
	}

	cnd.Reset()
	cnd.AddKV(`resource_id`, db.In(resourceIDs))
	cnd.AddKV(`row_id`, db.In(rowIDs))
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	_, err = tM.ListByOffset(nil, nil, 0, -1, cnd)
	if err != nil {
		return list, err
	}
	for _, row := range tM.Objects() {
		field := resourceFields[row.ResourceId]
		idx := idxByRowID[row.RowId]
		if _, ok := list[idx][`_translations`]; !ok {
			list[idx][`_translations`] = echo.H{
				row.Lang: echo.H{
					field: row.Text,
				},
			}
			continue
		}
		translations := list[idx][`_translations`].(echo.H)
		if _, ok := translations[row.Lang]; !ok {
			translations[row.Lang] = echo.H{
				field: row.Text,
			}
			continue
		}
		translations[row.Lang].(echo.H)[field] = row.Text
	}
	return list, err
}
