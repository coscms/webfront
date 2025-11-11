package i18nm

import (
	"strings"

	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

// GetTranslations retrieves translations for the specified table and row IDs.
// It returns a map where keys are row IDs and values are maps of field names to translated texts.
// The translations are filtered by the current request language and the specified table prefix.
func GetTranslations(ctx echo.Context, table string, ids []uint64) map[uint64]map[string]string {
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
func GetModelsTranslations[T factory.Model](ctx echo.Context, models []T) []T {
	if len(models) == 0 {
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
	translations := GetTranslations(ctx, models[0].Short_(), ids)
	for id, row := range translations {
		index := idk[id]
		rowI := map[string]interface{}{}
		for field, text := range row {
			rowI[field] = text
		}
		models[index].FromRow(rowI)
	}
	return models
}

// Initialize scans all database fields marked as multilingual and ensures they have corresponding
// entries in the i18n resource table. It creates new i18n resource records for any multilingual
// fields that don't already exist in the resource table. Returns any error encountered during
// the process.
func Initialize(ctx echo.Context) error {
	rM := dbschema.NewOfficialI18nResource(ctx)
	var exists bool
	var err error
	for table, fieldInfo := range dbschema.DBI.Fields {
		for field, info := range fieldInfo {
			if !info.Multilingual {
				continue
			}
			exists, err = rM.Exists(nil, `code`, table+`.`+field)
			if err != nil {
				return err
			}
			if exists {
				continue
			}
			rM.Code = table + `.` + field
			_, err = rM.Insert()
			if err != nil {
				return err
			}
			rM.Reset()
		}
	}
	return err
}
