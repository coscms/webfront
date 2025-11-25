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

func isMultilingual(){
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
func GetTranslations(ctx echo.Context, table string, ids []uint64) map[uint64]map[string]string {
 if !isMultilingual() {
  return map[uint64]map[string]string{}
 }
	if ttl, ok := ctx.Internal().Get(`translationsTTL`).(int64); ok {
		m := map[uint64]map[string]string{}
		cache.XFunc(ctx, `translations.`+table+`.`+com.JoinNumbers(ids, `_`), &m, func() error {
			r := getTranslations(ctx, table, ids)
			maps.Copy(m, r)
			return nil
		}, cache.AdminRefreshable(ctx, sessdata.Customer(ctx), cache.TTL(ttl)))
		return m
	}
	return getTranslations(ctx, table, ids)
}

func getTranslations(ctx echo.Context, table string, ids []uint64) map[uint64]map[string]string {
	m := map[uint64]map[string]string{}
	rM := dbschema.NewOfficialI18nResource(ctx)
	rM.ListByOffset(nil, nil, 0, -1, `code`, db.Like(table+`.%`))
	rows := rM.Objects()
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

// GetModelTranslations retrieves translations for multiple model instances by their IDs.
// It returns a map where each key is a model ID and the value is another map of language translations.
// The translations are fetched using the model's context and table name.
func GetModelTranslations(mdl factory.Model, ids []uint64) map[uint64]map[string]string {
	return GetTranslations(mdl.Context(), mdl.Short_(), ids)
}

// GetModelsTranslations retrieves translations for a slice of models and applies them to each model.
// It takes a slice of models as input and returns the same slice with translations applied.
// For each model, it extracts the ID, fetches translations using GetModelTranslations,
// and updates the model fields with the translated values.
// If the input slice is empty or any model lacks an ID field, it returns the original slice unchanged.
func GetModelsTranslations[T factory.Model](ctx echo.Context, models []T, tableName ...string) []T {
	if len(models) == 0 {
		return models
	} 
 if !isMultilingual() {
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
	if len(tableName) > 0 && len(tableName[0]) > 0 {
		table = tableName[0]
	}
	translations := GetTranslations(ctx, table, ids)
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
