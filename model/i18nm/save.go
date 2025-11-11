package i18nm

import (
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

// SaveModelTranslations saves multilingual field translations for a model.
// It processes all multilingual fields of the given model, creating/updating translations
// for each language configured in the system.
// Parameters:
//   - mdl: The model instance containing multilingual fields
//   - id: The row ID of the model to save translations for
//   - formNamePrefix: Optional prefix for form field names (defaults to "Language")
//
// Returns:
//   - error: Any error encountered during the save process
func SaveModelTranslations(mdl factory.Model, id uint64, formNamePrefix ...string) error {
	ctx := mdl.Context()
	rM := dbschema.NewOfficialI18nResource(ctx)
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	var err error
	table := mdl.Short_()
	fnpre := `Language`
	if len(formNamePrefix) > 0 && len(formNamePrefix[0]) > 0 {
		fnpre = formNamePrefix[0]
	}
	for field, info := range dbschema.DBI.Fields[table] {
		if !info.Multilingual {
			continue
		}
		err = rM.Get(nil, `code`, table+`.`+field)
		if err != nil {
			if err != db.ErrNoMoreRows {
				return err
			}
			rM.Code = table + `.` + field
			_, err = rM.Insert()
			if err != nil {
				return err
			}
		}
		resourceID := rM.Id
		tM.Reset()
		formName := com.CamelCase(field)
		formName2 := com.UpperCaseFirst(formName)
		for _, langCode := range config.FromFile().Language.AllList {
			translate := ctx.FormAny(fnpre+`[`+langCode+`][`+formName+`]`, fnpre+`[`+langCode+`][`+formName2+`]`)
			cond := db.And(
				db.Cond{`lang`: langCode},
				db.Cond{`row_id`: id},
				db.Cond{`resource_id`: resourceID},
			)
			if len(translate) == 0 {
				err = tM.Delete(nil, cond)
				if err != nil {
					return err
				}
				continue
			}
			err = tM.Get(nil, cond)
			if err != nil {
				if err != db.ErrNoMoreRows {
					return err
				}
				tM.Lang = langCode
				tM.ResourceId = resourceID
				tM.RowId = id
				tM.Text = translate
				_, err = tM.Insert()
			} else if translate != tM.Text {
				err = tM.UpdateField(nil, `text`, translate, cond)
			}
			if err != nil {
				return err
			}
			rM.Reset()
		}
	}
	return err
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
