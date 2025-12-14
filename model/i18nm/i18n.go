package i18nm

import (
	"database/sql"
	"maps"
	"strings"

	"github.com/coscms/webcore/library/config"
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
	return lang == config.FromFile().Language.Default
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

// getTranslations retrieves translations for specified table rows and columns
// Returns a nested map where outer key is row ID and inner map contains column-value pairs
// Parameters:
//   - ctx: echo context containing request information
//   - table: database table name
//   - ids: slice of row IDs to retrieve translations for
//   - columns: optional list of specific columns to retrieve (empty for all columns)
//
// Returns empty map if no translations found or on error
func getTranslations(ctx echo.Context, table string, ids []uint64, columns ...string) map[uint64]map[string]string {
	m := map[uint64]map[string]string{}
	if len(ids) == 0 {
		return m
	}
	rows, err := getResources(ctx, table, columns...)
	if err != nil || len(rows) == 0 {
		return m
	}
	rIDs := make([]uint, len(rows))
	rKeys := map[uint]string{}
	for i, v := range rows {
		rIDs[i] = v.Id
		rKeys[v.Id] = strings.SplitN(v.Code, `.`, 2)[1]
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
	var row *sql.Row
	row, err = p.SetCollection(dbschema.WithPrefix(table)).SetMW(func(r db.Result) db.Result {
		return r.Select(column)
	}).SetArgs(db.Cond{`row_id`: rowID}).QueryRow()
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = nil
		}
		return ``, err
	}
	result := sql.NullString{}
	err = row.Scan(&result)
	return result.String, err
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
