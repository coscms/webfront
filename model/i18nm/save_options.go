package i18nm

import "github.com/webx-top/echo"

type SaveModelTranslationsOptions struct {
	FormNamePrefix string
	ContentType    map[string]string // map[fieldName]contentType
	Project        string
	AutoTranslate  bool
	Translator     func(ctx echo.Context, fieldName string, value string) (string, error)
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
	o.AutoTranslate = autoTranslate
}

// SetTranslator sets the translator function for converting field values
func (o *SaveModelTranslationsOptions) SetTranslator(translator func(ctx echo.Context, fieldName string, value string) (string, error)) {
	o.Translator = translator
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
func OptionTranslator(translator func(ctx echo.Context, fieldName string, value string) (string, error)) func(*SaveModelTranslationsOptions) {
	return func(o *SaveModelTranslationsOptions) {
		o.SetTranslator(translator)
	}
}
