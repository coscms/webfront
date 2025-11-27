package i18nm

import (
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

func IsMultilingual() bool {
	return len(config.FromFile().Language.AllList) > 1
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
		condVal = db.In(columns)
	} else {
		condVal = db.Like(table + `.%`)
	}
	_, err := rM.ListByOffset(nil, nil, 0, -1, `code`, condVal)
	if err != nil {
		return nil, err
	}
	return rM.Objects(), err
}

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
func GetModelTranslationsByIDs(mdl factory.Model, ids []uint64, columns ...string) map[uint64]map[string]string {
	return GetTranslations(mdl.Context(), mdl.Short_(), ids, columns...)
}

// GetModelTranslations retrieves translations for a model instance by its ID.
// It returns translations as a map where each key is a field name and the value is the translated text.
// If the model lacks an ID field, it does nothing.
// If the context is nil, it uses the model's context.
// It fetches translations using the model's context and table name.
// If translations are found, it applies them to the model instance using the FromRow method.
func GetModelTranslations(ctx echo.Context, mdl factory.Model, columns ...string) {
	if !IsMultilingual() {
		return
	}
	var id uint64
	switch v := mdl.GetField(`Id`).(type) {
	case uint64:
		id = v
	case uint:
		id = uint64(v)
	default:
		return
	}
	if id == 0 {
		return
	}
	if ctx == nil {
		ctx = mdl.Context()
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
func GetModelsTranslations[T factory.Model](ctx echo.Context, models []T, columns ...string) []T {
	if len(models) == 0 {
		return models
	}
	if !IsMultilingual() {
		return models
	}
	if config.FromFile().Language.Default == ctx.Lang().Normalize() {
		return models
	}
	ids := make([]uint64, len(models))
	idk := map[uint64]int{}
	for index, row := range models {
		var id uint64
		switch v := row.GetField(`Id`).(type) {
		case uint64:
			id = v
		case uint:
			id = uint64(v)
		default:
			return models
		}
		if id == 0 {
			return models
		}
		ids[index] = id
		idk[id] = index
	}
	if len(ids) == 0 {
		return models
	}
	if ctx == nil {
		ctx = models[0].Context()
	}
	table := models[0].Short_()
	translations := GetTranslations(ctx, table, ids, columns...)
	for id, row := range translations {
		index := idk[id]
		rowI := map[string]interface{}{}
		for field, text := range row {
			if len(text) == 0 {
				continue
			}
			rowI[field] = text
		}
		models[index].FromRow(rowI)
	}
	return models
}
