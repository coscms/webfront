package i18nm

import (
	"maps"
	"strings"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/cache"
	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

// IsMultilingual returns true if the application supports multiple languages based on the configuration
func IsMultilingual() bool {
	return len(config.FromFile().Language.AllList) > 1
}

// LangIsDefault checks if the given language is the default language configured in the system.
func LangIsDefault(lang string) bool {
	return lang == LangDefault()
}

// LangDefault 返回配置文件中设置的默认语言代码
func LangDefault() string {
	return config.FromFile().Language.Default
}

// IsDefaultLang checks if the current language in context is the default language
func IsDefaultLang(ctx echo.Context) bool {
	return LangIsDefault(ctx.Lang().Normalize())
}

// SetTranstationsTTL sets the TTL (Time To Live) for translations in the given echo context.
// The ttl parameter specifies the duration in seconds that translations should be cached.
func SetTranstationsTTL(ctx echo.Context, ttl int64) {
	ctx.Internal().Set(`translationsTTL`, ttl)
}

// GetTranslations retrieves translations for the specified table and row IDs.
// It returns a map where keys are row IDs and values are maps of field names to translated texts.
// The translations are filtered by the current request language and the specified table prefix.
func GetTranslations(ctx echo.Context, table string, ids []uint64, columns ...string) map[uint64]map[string]string {
	if !IsMultilingual() {
		return map[uint64]map[string]string{}
	}
	if ttl, ok := ctx.Internal().Get(`translationsTTL`).(int64); ok {
		m := map[uint64]map[string]string{}
		cache.XFunc(ctx, `translations.`+table+`.`+com.JoinNumbers(ids, `_`), &m, func() error {
			r := getTranslations(ctx, table, ids, columns...)
			maps.Copy(m, r)
			return nil
		}, cache.AdminRefreshable(ctx, sessdata.Customer(ctx), cache.TTL(ttl)))
		return m
	}
	return getTranslations(ctx, table, ids, columns...)
}

// GetAllTranslations 获取指定表和ID集合的多语言翻译数据
//
// 参数:
//   - ctx: echo上下文对象
//   - table: 数据库表名
//   - ids: 需要查询的ID集合
//   - columns: 可选字段列表，指定要获取的列
//
// 返回值:
//   - 三层嵌套map结构: ID -> 语言代码 -> 字段名 -> 翻译值
//   - 如果未启用多语言功能，返回空map
func GetAllTranslations(ctx echo.Context, table string, ids []uint64, columns ...string) map[uint64]map[string]map[string]string {
	if !IsMultilingual() {
		return map[uint64]map[string]map[string]string{}
	}
	if ttl, ok := ctx.Internal().Get(`translationsTTL`).(int64); ok {
		m := map[uint64]map[string]map[string]string{}
		cache.XFunc(ctx, `allTranslations.`+table+`.`+com.JoinNumbers(ids, `_`), &m, func() error {
			r := getAllTranslations(ctx, table, ids, columns...)
			maps.Copy(m, r)
			return nil
		}, cache.AdminRefreshable(ctx, sessdata.Customer(ctx), cache.TTL(ttl)))
		return m
	}
	return getAllTranslations(ctx, table, ids, columns...)
}

// getResources retrieves i18n resources from the specified table with optional column filtering.
// It returns a slice of OfficialI18nResource objects and any error encountered during the operation.
// Parameters:
//   - ctx: echo context for database operations
//   - table: name of the table to query
//   - columns: optional list of columns to filter (if empty, matches all columns with LIKE pattern)
//
// Returns:
//   - []*dbschema.OfficialI18nResource: slice of retrieved resources
//   - error: any error that occurred during the query
func getResources(ctx echo.Context, table string, columns ...string) ([]*dbschema.OfficialI18nResource, error) {
	rM := dbschema.NewOfficialI18nResource(ctx)
	var condVal interface{}
	if len(columns) > 0 {
		for i, v := range columns {
			columns[i] = table + `.` + v
		}
		if len(columns) == 1 {
			condVal = columns[0]
		} else {
			condVal = db.In(columns)
		}
	} else {
		condVal = db.Like(table + `.%`)
	}
	_, err := rM.ListByOffset(nil, nil, 0, -1, `code`, condVal)
	if err != nil {
		return nil, err
	}
	return rM.Objects(), err
}

// getResourceIDs retrieves a list of resource IDs and a map of resource IDs to column names for the specified table and columns.
// It returns a slice of uint64 resource IDs and a map where the key is the resource ID and the value is the column name
// If an error occurs during the query, it returns an empty list and map and the error.
func getResourceIDs(ctx echo.Context, table string, columns ...string) ([]uint, map[uint]string, error) {
	rows, err := getResources(ctx, table, columns...)
	if err != nil {
		return nil, nil, err
	}
	rIDs := make([]uint, len(rows))
	rKeys := map[uint]string{}
	for i, v := range rows {
		rIDs[i] = v.Id
		rKeys[v.Id] = strings.SplitN(v.Code, `.`, 2)[1]
	}
	return rIDs, rKeys, err
}

// getTranslations retrieves translations for specified table rows and columns
// Returns a nested map where outer key is row ID and inner map contains column-value pairs
// Parameters:
//   - ctx: echo context containing request information
//   - table: database table name
//   - ids: slice of row IDs to retrieve translations for
//   - columns: optional list of specific columns to retrieve (empty for all columns)
//
// Returns empty map if no translations found or on error
//
//	{
//	  rowID1: {
//	    column1: "translated text 1",
//	    column2: "translated text 2",
//	  },
//	  rowID2: {
//	    column1: "translated text 3",
//	  },
//	}
func getTranslations(ctx echo.Context, table string, ids []uint64, columns ...string) map[uint64]map[string]string {
	m := map[uint64]map[string]string{}
	if len(ids) == 0 {
		return m
	}
	rIDs, rKeys, err := getResourceIDs(ctx, table, columns...)
	if err != nil || len(rIDs) == 0 {
		return m
	}
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	tM.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`lang`: ctx.Lang().Normalize()},
		db.Cond{`row_id`: db.In(ids)},
		db.Cond{`resource_id`: db.In(rIDs)},
	))
	tRows := tM.Objects()
	for _, v := range tRows {
		if _, ok := m[v.RowId]; !ok {
			m[v.RowId] = map[string]string{}
		}
		m[v.RowId][rKeys[v.ResourceId]] = v.Text
	}
	return m
}

// getAllTranslations 获取指定表、ID列表和列的多语言翻译结果
// 返回格式为 map[行ID]map[语言]map[字段名]翻译文本
// 如果ids为空或查询无结果，返回空map
func getAllTranslations(ctx echo.Context, table string, ids []uint64, columns ...string) map[uint64]map[string]map[string]string {
	m := map[uint64]map[string]map[string]string{}
	if len(ids) == 0 {
		return m
	}
	rIDs, rKeys, err := getResourceIDs(ctx, table, columns...)
	if err != nil || len(rIDs) == 0 {
		return m
	}
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	tM.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`row_id`: db.In(ids)},
		db.Cond{`resource_id`: db.In(rIDs)},
	))
	tRows := tM.Objects()
	for _, v := range tRows {
		if _, ok := m[v.RowId]; !ok {
			m[v.RowId] = map[string]map[string]string{}
		}
		if _, ok := m[v.RowId][v.Lang]; !ok {
			m[v.RowId][v.Lang] = map[string]string{}
		}
		m[v.RowId][v.Lang][rKeys[v.ResourceId]] = v.Text
	}
	return m
}

// GetModelTranslationsByIDs retrieves translations for multiple model instances by their IDs.
// It returns a map where each key is a model ID and the value is another map of language translations.
// The translations are fetched using the model's context and table name.
func GetModelTranslationsByIDs(ctx echo.Context, mdl Model, ids []uint64, columns ...string) map[uint64]map[string]string {
	return GetTranslations(ctx, mdl.Short_(), ids, columns...)
}

// GetModelRowID retrieves the row ID of a model based on the specified column and text in the current language context.
// It first gets the resource ID from the model and column, then looks up the translation matching the given text.
// Returns the row ID if found, or 0 with an error if the resource or translation cannot be found.
func GetModelRowID(ctx echo.Context, mdl Model, column string, text string) (uint64, error) {
	return GetResourceRowID(ctx, mdl.Short_(), column, text)
}

// GetResourceRowID retrieves the row ID for a given internationalized text in the specified table and column.
// It first looks up the resource ID for the table-column pair, then finds the translation matching the given text and language.
// Returns the row ID (0 if not found) and any error that occurred during the lookup.
func GetResourceRowID(ctx echo.Context, table string, column string, text string) (uint64, error) {
	rows, err := getResources(ctx, table, column)
	if err != nil || len(rows) == 0 {
		return 0, err
	}
	resID := rows[0].Id
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	err = tM.Get(func(r db.Result) db.Result {
		return r.Select(`row_id`)
	}, db.And(
		db.Cond{`lang`: ctx.Lang().Normalize()},
		db.Cond{`resource_id`: resID},
		db.Cond{`text`: text},
	))
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = nil
		}
	}
	return tM.RowId, nil
}

// GetColumnDefaultLangText retrieves the default language text for a specific column in a table.
// It first gets the resource row ID using GetResourceRowID, then queries the database for the column value.
// Returns the text in default language or empty string if not found, along with any error encountered.
func GetColumnDefaultLangText(ctx echo.Context, table string, column string, text string) (string, error) {
	if !IsMultilingual() || IsDefaultLang(ctx) {
		return text, nil
	}
	rowID, err := GetResourceRowID(ctx, table, column, text)
	if err != nil {
		return ``, err
	}
	p := factory.ParamPoolGet()
	defer p.Release()
	row := struct {
		V string `db:"_v"`
	}{}
	err = p.SetCollection(dbschema.WithPrefix(table)).SetMW(func(r db.Result) db.Result {
		return r.Select(db.Raw(column + ` AS _v`))
	}).SetArgs(db.Cond{`id`: rowID}).SetRecv(&row).One()
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = nil
		}
		return ``, err
	}
	return row.V, err
}

// GetModelTranslations retrieves translations for a model instance by its ID.
// It returns translations as a map where each key is a field name and the value is the translated text.
// If the model lacks an ID field, it does nothing.
// If the context is nil, it uses the model's context.
// It fetches translations using the model's context and table name.
// If translations are found, it applies them to the model instance using the FromRow method.
func GetModelTranslations(ctx echo.Context, mdl Model, columns ...string) {
	if !IsMultilingual() {
		return
	}
	id := GetRowID(mdl)
	if id == 0 {
		return
	}
	translations := GetTranslations(ctx, mdl.Short_(), []uint64{id}, columns...)
	if len(translations) > 0 && len(translations[id]) > 0 {
		rowI := map[string]interface{}{}
		for field, text := range translations[id] {
			if len(text) == 0 {
				continue
			}
			rowI[field] = text
		}
		mdl.FromRow(rowI)
	}
}

// GetModelsTranslations retrieves translations for a slice of models and applies them to each model.
// It takes a slice of models as input and returns the same slice with translations applied.
// For each model, it extracts the ID, fetches translations using GetModelTranslations,
// and updates the model fields with the translated values.
// If the input slice is empty or any model lacks an ID field, it returns the original slice unchanged.
func GetModelsTranslations[T Model](ctx echo.Context, models []T, columns ...string) []T {
	if len(models) == 0 {
		return models
	}
	if !IsMultilingual() || IsDefaultLang(ctx) {
		return models
	}
	ids := make([]uint64, 0, len(models))
	idk := map[uint64][]int{}
	for index, row := range models {
		id := GetRowID(row)
		if id == 0 {
			return models
		}
		if _, ok := idk[id]; !ok {
			idk[id] = []int{}
			ids = append(ids, id)
		}
		idk[id] = append(idk[id], index)
	}
	if len(ids) == 0 {
		return models
	}
	table := models[0].Short_()
	translations := GetTranslations(ctx, table, ids, columns...)
	for id, row := range translations {
		mp := map[string]interface{}{}
		for field, text := range row {
			if len(text) == 0 {
				continue
			}
			mp[field] = text
		}
		for _, index := range idk[id] {
			models[index].FromRow(mp)
		}
	}
	return models
}

// GetModelsAllTranslations 获取多个模型的所有翻译数据
//
// 参数:
//
//	ctx - echo 上下文对象
//	models - 模型实例切片
//	columns - 可选的要获取的字段列表
//
// 返回值:
//
//	返回一个切片，每个元素是一个映射，包含对应模型的多语言翻译数据
//	映射的键是语言代码，值是字段到翻译文本的映射
func GetModelsAllTranslations[T Model](ctx echo.Context, models []T, columns ...string) []map[string]echo.H {
	var result []map[string]echo.H
	if len(models) == 0 {
		return result
	}
	if !IsMultilingual() {
		return result
	}
	ids := make([]uint64, 0, len(models))
	idk := map[uint64][]int{}
	for index, row := range models {
		id := GetRowID(row)
		if id == 0 {
			return result
		}
		if _, ok := idk[id]; !ok {
			idk[id] = []int{}
			ids = append(ids, id)
		}
		idk[id] = append(idk[id], index)
	}
	if len(ids) == 0 {
		return result
	}
	table := models[0].Short_()
	translations := GetAllTranslations(ctx, table, ids, columns...)
	result = make([]map[string]echo.H, len(models))
	for id, row := range translations {
		mp := map[string]echo.H{}
		for lang, texts := range row {
			vals := echo.H{}
			for field, text := range texts {
				if len(text) == 0 {
					continue
				}
				vals[field] = text
			}
			mp[lang] = vals
		}
		for _, index := range idk[id] {
			result[index] = mp
		}
	}
	return result
}

// TranslationsToMapByIndex converts a slice of language sets to a map.
// It takes a slice of language sets and an index as input and returns a map where each key is a field name
// and the value is the translated text for the given index.
// The map keys are in the format: "Language[lang][field]" where lang is the language code and field is the field name.
func TranslationsToMapByIndex(langsets []map[string]echo.H, index int) echo.H {
	result := echo.H{}
	for lang, langset := range langsets[index] {
		for key, val := range langset {
			result[formbuilder.FormInputNamePrefixDefault+`[`+lang+`][`+key+`]`] = val
		}
	}
	return result
}

// TranslationsToMaps converts a slice of language sets to a slice of maps.
// Each map in the output slice contains translations for a single language set.
// The map keys are in the format: "Language[lang][field]" where lang is the language code
// and field is the field name.
func TranslationsToMaps(langsets []map[string]echo.H) []echo.H {
	results := make([]echo.H, 0, len(langsets))
	for _, items := range langsets {
		result := echo.H{}
		for lang, langset := range items {
			for key, val := range langset {
				result[formbuilder.FormInputNamePrefixDefault+`[`+lang+`][`+key+`]`] = val
			}
		}
		results = append(results, result)
	}
	return results
}
