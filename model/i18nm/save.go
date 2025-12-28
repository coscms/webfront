package i18nm

import (
	"strings"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/fileupdater/listener"
	"github.com/coscms/webcore/library/formbuilder"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/top"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
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
func SaveModelTranslations(ctx echo.Context, mdl Model, id uint64, options ...func(*SaveModelTranslationsOptions)) error {
	if !IsMultilingual() {
		return nil
	}
	rM := dbschema.NewOfficialI18nResource(ctx)
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	var err error
	table := mdl.Short_()
	cfg := SaveModelTranslationsOptions{
		ContentType: map[string]string{},
	}
	for _, fn := range options {
		fn(&cfg)
	}
	cfg.SetDefaults()
	var hasUpload func(field string) bool
	if pMap, ok := listener.UpdaterInfos[cfg.Project]; ok {
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
	var autoTranslate bool
	var forceTranslate bool
	var allowForceTranslate bool
	var trimOverflowText bool
	if cfg.AutoTranslate != nil {
		autoTranslate = *cfg.AutoTranslate
	}
	if cfg.AllowForceTranslate != nil {
		allowForceTranslate = cfg.AllowForceTranslate(ctx)
	}
	if allowForceTranslate {
		if cfg.ForceTranslate != nil {
			forceTranslate = *cfg.ForceTranslate
		} else {
			forceTranslate = ctx.Formx(`forceTranslate`).Bool()
		}
	}
	if cfg.TrimOverflowText != nil {
		trimOverflowText = *cfg.TrimOverflowText
	}
	langCfg := config.FromFile().Language
	for field, info := range dbschema.DBI.Fields[table] {
		if !info.Multilingual {
			continue
		}
		var resourceID uint
		if cfg.resourceIDsByField != nil {
			var ok bool
			resourceID, ok = cfg.resourceIDsByField[field]
			if !ok {
				continue
			}
		} else {
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
			resourceID = rM.Id
			tM.Reset()
		}
		ctx.Internal().Set(`i18n_translation_resource_field`, field)
		formNameL := com.CamelCase(field)
		formNameU := com.UpperCaseFirst(formNameL)
		langDefault := langCfg.Default
		originalText, _ := mdl.GetField(formNameU).(string)
		var picks []string
		if field == `content` {
			picks, originalText = top.PickoutHideTag(originalText)
		}
		for _, langCode := range langCfg.AllList {
			if langDefault == langCode {
				continue
			}
			translatedText := ctx.FormAny(cfg.FormNamePrefix+`[`+langCode+`][`+formNameL+`]`, cfg.FormNamePrefix+`[`+langCode+`][`+formNameU+`]`)
			cond := db.And(
				db.Cond{`lang`: langCode},
				db.Cond{`row_id`: id},
				db.Cond{`resource_id`: resourceID},
			)
			var contentType string
			if len(cfg.ContentType) > 0 {
				contentType, _ = cfg.ContentType[field]
			}
			if len(contentType) == 0 {
				contentType = `string` // 默认string类型(单行文本)
			}
			if forceTranslate {
				translatedText, err = cfg.Translate(ctx, field, translatedText, originalText, contentType, langCode, langDefault)
				if err != nil {
					return err
				}
				if len(picks) > 0 {
					translatedText = top.RestorePickoutedHideTag(translatedText, picks)
				}
			} else if len(translatedText) == 0 && autoTranslate {
				translatedText, err = cfg.Translate(ctx, field, translatedText, originalText, contentType, langCode, langDefault)
				if err != nil {
					return err
				}
				if len(picks) > 0 {
					translatedText = top.RestorePickoutedHideTag(translatedText, picks)
				}
			}
			if len(translatedText) == 0 {
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
			translatedText = common.ContentEncode(translatedText, contentType)
			if trimOverflowText && info.GetMaxSize() > 0 && len(translatedText) > info.GetMaxSize() {
				translatedText = info.TrimOverflowText(translatedText)
			}
			err = tM.Get(nil, cond)
			if err != nil {
				if err != db.ErrNoMoreRows {
					return err
				}
				tM.Lang = langCode
				tM.ResourceId = resourceID
				tM.RowId = id
				tM.Text = translatedText
				if hasUpload(field) {
					_, err = tM.Insert()
				} else {
					_, err = tM.EventOFF().Insert()
				}
			} else if translatedText != tM.Text {
				if hasUpload(field) {
					err = tM.UpdateField(nil, `text`, translatedText, cond)
				} else {
					err = tM.EventOFF().UpdateField(nil, `text`, translatedText, cond)
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
func SetModelTranslationsToForm(ctx echo.Context, mdl Model, id uint64, formNamePrefix ...string) error {
	if !IsMultilingual() {
		return nil
	}
	table := mdl.Short_()
	namePrefix := formbuilder.FormInputNamePrefixDefault
	if len(formNamePrefix) > 0 && len(formNamePrefix[0]) > 0 {
		namePrefix = formNamePrefix[0]
	}
	rM := dbschema.NewOfficialI18nResource(ctx)
	_, err := rM.ListByOffset(nil, nil, 0, -1, `code`, db.Like(table+`.%`))
	if err != nil {
		return err
	}
	rows := rM.Objects()
	if len(rows) == 0 {
		return nil
	}
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
func DeleteModelTranslations(ctx echo.Context, mdl Model, id uint64) error {
	if !IsMultilingual() {
		return nil
	}
	table := mdl.Short_()
	rM := dbschema.NewOfficialI18nResource(ctx)
	_, err := rM.ListByOffset(nil, nil, 0, -1, `code`, db.Like(table+`.%`))
	if err != nil {
		return err
	}
	rows := rM.Objects()
	if len(rows) == 0 {
		return nil
	}
	rIDs := make([]uint, len(rows))
	rCodes := map[uint]string{}
	for i, v := range rows {
		rIDs[i] = v.Id
		rCodes[v.Id] = v.Code
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
