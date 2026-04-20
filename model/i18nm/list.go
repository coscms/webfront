package i18nm

import (
	"errors"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/cache"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"
)

type ListItem struct {
	*dbschema.OfficialI18nTranslation
	*dbschema.OfficialI18nResource
}

type ListQuery struct {
	Table string
	RowID uint64
	Lang  string
	Sorts []any
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

func UpdateColumnTranslation(ctx echo.Context, table string, column string, rowID uint64, lang string, translatedText string) (affected int64, err error) {
	if len(lang) == 0 || LangIsDefault(lang) {
		pr := factory.NewParam().SetContext(ctx)
		affected, err = pr.SetCollection(table).SetSend(echo.H{
			column: translatedText,
		}).SetArgs(`id`, rowID).Updatex()
		return
	}
	var resourceIDs []uint
	resourceIDs, _, err = getResourceIDs(ctx, table, column)
	if err != nil {
		return
	}
	if len(resourceIDs) == 0 {
		return
	}

	tM := dbschema.NewOfficialI18nTranslation(ctx)
	cnd := db.NewCompounds()
	cnd.AddKV(`lang`, lang)
	cnd.AddKV(`row_id`, rowID)
	cnd.AddKV(`resource_id`, resourceIDs[0])
	affected, err = tM.UpdatexFields(nil, echo.H{
		`text`: translatedText,
	}, cnd.And())
	if err != nil {
		return
	}
	if affected < 1 {
		tM.Lang = lang
		tM.ResourceId = resourceIDs[0]
		tM.RowId = rowID
		tM.Text = translatedText
		_, err = tM.Insert()
		if err != nil {
			return
		}
		affected = 1
	}
	return
}

func ListByResource(ctx echo.Context, query ListQuery) ([]echo.H, error) {
	var list []echo.H
	resourceIDs, resourceFields, err := getResourceIDs(ctx, query.Table)
	if err != nil {
		return list, err
	}
	cnd := db.NewCompounds()
	if query.RowID > 0 {
		cnd.AddKV(`id`, query.RowID)
	}
	if len(query.Sorts) == 0 {
		query.Sorts = append(query.Sorts, `-id`)
	}
	columns := make([]any, 0, len(resourceFields)+1)
	columns = append(columns, `id`)
	for _, column := range resourceFields {
		columns = append(columns, column)
	}
	smw := func(r db.Result) db.Result {
		return r.Select(columns...).OrderBy(query.Sorts...)
	}
	pr := factory.ParamPoolGet().SetContext(ctx)
	ls := pr.SetCollection(query.Table).SetRecv(&list).NewLister()
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
	if len(query.Lang) > 0 {
		cnd.AddKV(`lang`, query.Lang)
	}
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

func Batch(ctx echo.Context, query ListQuery, np notice.NProgressor, restartID ...uint64) error {
	cfg := DefaultSaveModelTranslationsOptions
	if cfg.AllowForceTranslate == nil {
		return errors.New("AllowForceTranslate function is not set in configuration")
	}
	forceTranslate := cfg.AllowForceTranslate(ctx)
	if cfg.translator == nil {
		return errors.New("Translator function is not set in configuration")
	}

	var list []echo.H
	_, resourceFields, err := getResourceIDs(ctx, query.Table)
	if err != nil {
		return err
	}
	cacheKey := `translation.` + query.Table
	var cacheExpire int64 = 86400 * 365 * 10
	var lastID uint64
	cnd := db.NewCompounds()
	if query.RowID > 0 {
		cnd.AddKV(`id`, query.RowID)
	} else if len(restartID) > 0 {
		lastID = restartID[0]
		cnd.AddKV(`id`, db.Gte(lastID))
	} else {
		err = cache.Get(ctx, cacheKey, &lastID)
		if err != nil && !cache.IsNotExist(err) {
			return err
		}
		cnd.AddKV(`id`, db.Gt(lastID))
	}
	if len(query.Sorts) == 0 {
		query.Sorts = append(query.Sorts, `id`)
	}
	columnsResourceID := map[string]uint{}
	columns := make([]any, 0, len(resourceFields)+1)
	columns = append(columns, `id`)
	for resourceID, column := range resourceFields {
		columns = append(columns, column)
		columnsResourceID[column] = resourceID
	}
	smw := func(r db.Result) db.Result {
		return r.Select(columns...).OrderBy(query.Sorts...)
	}
	pr := factory.ParamPoolGet().SetContext(ctx)
	ls := pr.SetCollection(query.Table).SetRecv(&list).NewLister()
	offsetLister := pagination.NewOffsetLister(ls, &list, smw, cnd.And())
	np.Reset()
	var initedNP bool
	offsetLister.SetProg(func(offset, total int64) {
		if initedNP {
			return
		}
		np.Add(total)
		initedNP = true
	})
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	translate := cfg.Translate
	langCfg := config.FromFile().Language
	err = offsetLister.ChunkListNoOffset(func() (db.Compound, error) {
		_lastID := lastID
		for _, row := range list {
			np.OnlyAdd(1)
			rowID := row.Uint64(`id`)
			if rowID == 0 {
				continue
			}
			lastID = rowID
			for column, value := range row {
				if column == `id` {
					continue
				}
				if value == nil {
					continue
				}
				originalText, ok := value.(string)
				if !ok {
					continue
				}
				resourceID, ok := columnsResourceID[column]
				if !ok {
					continue
				}
				var restoreFunc func(translatedText string) string
				if cfg.originalTextPickout != nil {
					originalText, restoreFunc = cfg.originalTextPickout(query.Table, column, originalText)
				}
				var langList []string
				if len(query.Lang) > 0 {
					langList = []string{query.Lang}
				} else {
					langList = langCfg.AllList
				}
				np.Success(ctx.T(`开始翻译“%s”...`, originalText))
				for _, langCode := range langList {
					if LangIsDefault(langCode) {
						continue
					}
					np.Success(ctx.T(`正在翻译成 %s ...`, langCode))
					translatedText, err := translateText(ctx, `string`, translate, restoreFunc, forceTranslate, true, column, originalText, ``, langCode, langCfg.Default)
					if err != nil {
						np.Failure(err.Error())
						return nil, err
					}
					tM.RowId = rowID
					tM.ResourceId = resourceID
					tM.Lang = langCode
					tM.Text = translatedText
					affected, err := tM.UpdatexFields(nil, echo.H{
						`text`: tM.Text,
					}, db.And(
						db.Cond{`row_id`: tM.RowId},
						db.Cond{`resource_id`: tM.ResourceId},
						db.Cond{`lang`: tM.Lang},
					))
					if err != nil {
						np.Failure(err.Error())
						return nil, err
					}
					if affected > 0 {
						np.Success(ctx.T(`更新成功`))
						continue
					}
					_, err = tM.Insert()
					if err != nil {
						np.Failure(err.Error())
						return nil, err
					}
					np.Success(ctx.T(`创建成功`))
				}
			}

			if query.RowID == 0 {
				continue
			}
			err = cache.Put(ctx, cacheKey, lastID, cacheExpire)
			if err != nil {
				return nil, err
			}
		}
		if _lastID == lastID {
			return nil, db.ErrNoMoreRows
		}
		return db.Cond{`id`: db.Gt(lastID)}, nil
	}, 100)
	pr.Release()
	return err
}
