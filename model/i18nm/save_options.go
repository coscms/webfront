package i18nm

import (
	"strings"

	"github.com/coscms/webcore/library/formbuilder"
	"github.com/webx-top/echo"
)

// DefaultSaveModelTranslationsOptions is the default options for SaveModelTranslations.
var DefaultSaveModelTranslationsOptions = SaveModelTranslationsOptions{
	FormNamePrefix: formbuilder.FormInputNamePrefixDefault,
	ContentType:    map[string]string{},
	Project:        "",
}

// Translator is a function that translates a field value to the specified language code.
type Translator = func(ctx echo.Context, fieldName string, value string, originalValue string, contentType string, langCode string, originalLangCode string) (string, error)

// SaveModelTranslationsOptions is a struct that holds options for saving model translations.
type SaveModelTranslationsOptions struct {
	FormNamePrefix      string
	ContentType         map[string]string // map[fieldName]contentType
	Project             string
	AutoTranslate       *bool
	ForceTranslate      *bool
	TrimOverflowText    *bool
	AllowForceTranslate func(echo.Context) bool
	translator          Translator
	resourceIDsByField  map[string]uint
}

// SetDefaults sets default values for SaveModelTranslationsOptions fields
// when they are not explicitly set. It copies values from DefaultSaveModelTranslationsOptions
// to the receiver for empty/zero fields including FormNamePrefix, Project,
// translator, AutoTranslate and ContentType.
func (o *SaveModelTranslationsOptions) SetDefaults() {
	d := DefaultSaveModelTranslationsOptions
	if len(o.FormNamePrefix) == 0 {
		o.FormNamePrefix = d.FormNamePrefix
	}
	if len(o.Project) == 0 && len(d.Project) > 0 {
		o.Project = d.Project
	}
	if o.translator == nil && d.translator != nil {
		o.translator = d.translator
	}
	if o.AutoTranslate == nil && d.AutoTranslate != nil {
		o.AutoTranslate = d.AutoTranslate
	}
	if o.ForceTranslate == nil && d.ForceTranslate != nil {
		o.ForceTranslate = d.ForceTranslate
	}
	if o.AllowForceTranslate == nil && d.AllowForceTranslate != nil {
		o.AllowForceTranslate = d.AllowForceTranslate
	}
	if o.ContentType == nil {
		o.ContentType = map[string]string{}
	}
	if len(o.ContentType) == 0 && len(d.ContentType) > 0 {
		for fieldName, contentType := range d.ContentType {
			o.ContentType[fieldName] = contentType
		}
	}
}

// SetContentType sets the content type for the specified field name
func (o *SaveModelTranslationsOptions) SetContentType(fieldName string, contentType string) {
	o.ContentType[fieldName] = contentType
}

// SetFormNamePrefix sets the prefix for form field names
func (o *SaveModelTranslationsOptions) SetFormNamePrefix(formNamePrefix string) {
	o.FormNamePrefix = formNamePrefix
}

// SetProject sets the project name for the translation options
func (o *SaveModelTranslationsOptions) SetProject(project string) {
	o.Project = project
}

// SetAutoTranslate sets whether translations should be automatically generated when missing
func (o *SaveModelTranslationsOptions) SetAutoTranslate(autoTranslate bool) {
	o.AutoTranslate = &autoTranslate
}

// SetTrimOverflowText sets whether to trim overflow text when saving translations
func (o *SaveModelTranslationsOptions) SetTrimOverflowText(trimOverflowText bool) {
	o.TrimOverflowText = &trimOverflowText
}

// SetTranslator sets the translator function for converting field values
func (o *SaveModelTranslationsOptions) SetTranslator(translator Translator) {
	o.translator = translator
}

// Translate translates the given field value for the specified language code.
// If a translator function is set in options, it will be used for translation.
// Returns the translated value or the original value if no translator is set.
func (o *SaveModelTranslationsOptions) Translate(ctx echo.Context, fieldName string, value string, originalValue string, contentType string, langCode string, originalLangCode string) (string, error) {
	if len(strings.TrimSpace(originalValue)) == 0 {
		return value, nil
	}
	if o.translator != nil {
		return o.translator(ctx, fieldName, value, originalValue, contentType, langCode, originalLangCode)
	}
	return value, nil
}

// SetForceTranslate sets whether to force translation updates regardless of existing translations
func (o *SaveModelTranslationsOptions) SetForceTranslate(forceTranslate bool) {
	o.ForceTranslate = &forceTranslate
}

// SetAllowForceTranslate sets whether to allow force translation updates
func (o *SaveModelTranslationsOptions) SetAllowForceTranslate(allowForceTranslate func(echo.Context) bool) {
	o.AllowForceTranslate = allowForceTranslate
}

// SetResourceIDsByField sets the resource IDs for each field
func (o *SaveModelTranslationsOptions) SetResourceIDsByField(resourceIDsByField map[string]uint) {
	o.resourceIDsByField = resourceIDsByField
}

// OptionContentType returns a function that sets the content type for the specified field
// in SaveModelTranslationsOptions. The returned function can be used as an option when saving
// model translations.
func OptionContentType(fieldName string, contentType string) func(*SaveModelTranslationsOptions) {
	return func(o *SaveModelTranslationsOptions) {
		o.SetContentType(fieldName, contentType)
	}
}

// OptionContentTypes returns a function that sets content types for multiple fields in SaveModelTranslationsOptions.
// The input map associates field names with their corresponding content types.
func OptionContentTypes(ct map[string]string) func(*SaveModelTranslationsOptions) {
	return func(o *SaveModelTranslationsOptions) {
		for fieldName, contentType := range ct {
			o.SetContentType(fieldName, contentType)
		}
	}
}

// OptionFormNamePrefix sets the form name prefix for SaveModelTranslationsOptions
func OptionFormNamePrefix(formNamePrefix string) func(*SaveModelTranslationsOptions) {
	return func(o *SaveModelTranslationsOptions) {
		o.SetFormNamePrefix(formNamePrefix)
	}
}

// OptionProject sets the project name for SaveModelTranslationsOptions
func OptionProject(project string) func(*SaveModelTranslationsOptions) {
	return func(o *SaveModelTranslationsOptions) {
		o.SetProject(project)
	}
}

// OptionAutoTranslate sets whether translations should be automatically generated when missing
func OptionAutoTranslate(autoTranslate bool) func(*SaveModelTranslationsOptions) {
	return func(o *SaveModelTranslationsOptions) {
		o.SetAutoTranslate(autoTranslate)
	}
}

// OptionTranslator sets the translator function for SaveModelTranslationsOptions
func OptionForceTranslate(forceTranslate bool) func(*SaveModelTranslationsOptions) {
	return func(o *SaveModelTranslationsOptions) {
		o.SetForceTranslate(forceTranslate)
	}
}

// OptionAllowForceTranslate 返回一个配置函数，用于设置是否允许强制翻译
// allowForceTranslate: 是否允许强制覆盖现有翻译
func OptionAllowForceTranslate(allowForceTranslate func(echo.Context) bool) func(*SaveModelTranslationsOptions) {
	return func(o *SaveModelTranslationsOptions) {
		o.SetAllowForceTranslate(allowForceTranslate)
	}
}

// OptionTrimOverflowText returns a function option that sets whether to trim overflow text when saving translations
// trimOverflowText: if true, will trim text that exceeds the field length limit
func OptionTrimOverflowText(trimOverflowText bool) func(*SaveModelTranslationsOptions) {
	return func(o *SaveModelTranslationsOptions) {
		o.SetTrimOverflowText(trimOverflowText)
	}
}

// OptionTranslator sets the translator function for SaveModelTranslationsOptions
func OptionTranslator(translator Translator) func(*SaveModelTranslationsOptions) {
	return func(o *SaveModelTranslationsOptions) {
		o.SetTranslator(translator)
	}
}

// OptionResourceIDsByField returns a function option that sets the resource IDs for each field.
// The input map associates field names with their corresponding resource IDs.
// The resource IDs are used to retrieve translations for the associated fields.
func OptionResourceIDsByField(resourceIDsByField map[string]uint) func(*SaveModelTranslationsOptions) {
	return func(o *SaveModelTranslationsOptions) {
		o.SetResourceIDsByField(resourceIDsByField)
	}
}
