package i18nm

import (
	"strings"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/fileupdater/listener"
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
	if !IsMultilingual() {
		return nil
	}
	ctx := mdl.Context()
	rM := dbschema.NewOfficialI18nResource(ctx)
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	var err error
	table := mdl.Short_()
	namePrefix := `Language`
	if len(formNamePrefix) > 0 && len(formNamePrefix[0]) > 0 {
		namePrefix = formNamePrefix[0]
	}
	var hasUpload func(field string) bool
	if pMap, ok := listener.UpdaterInfos[``]; ok {
		if updaterInfo, ok := pMap[table]; ok && updaterInfo != nil {
			hasUpload = func(field string) bool {
				_, ok := updaterInfo[field]
				return ok
			}
		}
	}
	if hasUpload == nil {
		hasUpload = func(field string) bool {
			return false
		}
	}
	ctx.Internal().Set(`i18n_translation_resource_table`, table)
	defer func() {
		ctx.Internal().Delete(`i18n_translation_resource_table`)
		ctx.Internal().Delete(`i18n_translation_resource_field`)
	}()
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
		ctx.Internal().Set(`i18n_translation_resource_field`, field)
		resourceID := rM.Id
		tM.Reset()
		formName := com.CamelCase(field)
		formName2 := com.UpperCaseFirst(formName)
		langDefault := config.FromFile().Language.Default
		for _, langCode := range config.FromFile().Language.AllList {
			if langDefault == langCode {
				continue
			}
			translate := ctx.FormAny(namePrefix+`[`+langCode+`][`+formName+`]`, namePrefix+`[`+langCode+`][`+formName2+`]`)
			cond := db.And(
				db.Cond{`lang`: langCode},
				db.Cond{`row_id`: id},
				db.Cond{`resource_id`: resourceID},
			)
			if len(translate) == 0 {
				if hasUpload(field) {
					err = tM.Delete(nil, cond)
				} else {
					err = tM.EventOFF().Delete(nil, cond)
				}
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
				if hasUpload(field) {
					_, err = tM.Insert()
				} else {
					_, err = tM.EventOFF().Insert()
				}
			} else if translate != tM.Text {
				if hasUpload(field) {
					err = tM.UpdateField(nil, `text`, translate, cond)
				} else {
					err = tM.EventOFF().UpdateField(nil, `text`, translate, cond)
				}
			}
			if err != nil {
				return err
			}
			rM.Reset()
		}
	}
	return err
}

// SetModelTranslationsToForm sets translation texts from a model to form fields.
// It retrieves translations for the given model ID and populates the form with language-specific values.
// The form field names are prefixed with the given prefix (default "Language") in the format: prefix[lang][field].
// Returns any error encountered during the operation.
func SetModelTranslationsToForm(mdl factory.Model, id uint64, formNamePrefix ...string) error {
	if !IsMultilingual() {
		return nil
	}
	ctx := mdl.Context()
	table := mdl.Short_()
	namePrefix := `Language`
	if len(formNamePrefix) > 0 && len(formNamePrefix[0]) > 0 {
		namePrefix = formNamePrefix[0]
	}
	rM := dbschema.NewOfficialI18nResource(ctx)
	_, err := rM.ListByOffset(nil, nil, 0, -1, `code`, db.Like(table+`.%`))
	if err != nil {
		return err
	}
	rows := rM.Objects()
	rIDs := make([]uint, len(rows))
	rKeys := map[uint]string{}
	for i, v := range rows {
		rIDs[i] = v.Id
		rKeys[v.Id] = com.CamelCase(strings.SplitN(v.Code, `.`, 2)[1])
	}
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	_, err = tM.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`row_id`: id},
		db.Cond{`resource_id`: db.In(rIDs)},
	))
	tRows := tM.Objects()
	for _, v := range tRows {
		field := rKeys[v.ResourceId]
		ctx.Request().Form().Set(namePrefix+`[`+v.Lang+`][`+field+`]`, v.Text)
	}
	langDefault := config.FromFile().Language.Default
	for _, info := range dbschema.DBI.Fields[table] {
		if !info.Multilingual {
			continue
		}
		ctx.Request().Form().Set(namePrefix+`[`+langDefault+`][`+com.LowerCaseFirst(info.GoName)+`]`, com.String(mdl.GetField(info.GoName)))
	}
	return err
}

// DeleteModelTranslations deletes all translations associated with a specific model instance.
// It removes both the resource entries and their corresponding translations from the database.
// Parameters:
//
//	mdl - the model instance containing context and table information
//	id  - the ID of the model instance whose translations should be deleted
//
// Returns:
//
//	error - any error encountered during the deletion process
func DeleteModelTranslations(mdl factory.Model, id uint64) error {
	if !IsMultilingual() {
		return nil
	}
	ctx := mdl.Context()
	table := mdl.Short_()
	rM := dbschema.NewOfficialI18nResource(ctx)
	_, err := rM.ListByOffset(nil, nil, 0, -1, `code`, db.Like(table+`.%`))
	if err != nil {
		return err
	}
	rows := rM.Objects()
	rIDs := make([]uint, len(rows))
	rCodes := map[uint]string{}
	for i, v := range rows {
		rIDs[i] = v.Id
		rCodes[v.Id] = v.Code
	}
	if len(rIDs) == 0 {
		return nil
	}
	cond := db.And(
		db.Cond{`row_id`: id},
		db.Cond{`resource_id`: db.In(rIDs)},
	)
	ctx.Internal().Set(`i18n_translation_resource_table`, table)
	ctx.Internal().Set(`i18n_translation_resource_codes`, rCodes)
	defer func() {
		ctx.Internal().Delete(`i18n_translation_resource_table`)
		ctx.Internal().Delete(`i18n_translation_resource_codes`)
	}()
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	return tM.Delete(nil, cond)
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
