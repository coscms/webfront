package translate

import (
	"github.com/admpub/translate"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/config/extend"
	"github.com/coscms/webfront/model/i18nm"

	_ "github.com/admpub/translate/providers"
)

// init registers the default translation function and extends the system with a new translation configuration.
// It sets Translate as the default translator for model translations and registers a factory function
// for creating new translation configurations under the "translate" extension point.
func init() {
	i18nm.DefaultSaveModelTranslationsOptions.SetTranslator(Translate)
	extend.Register(`translate`, func() interface{} {
		return NewConfig()
	})
}

// Translate translates the given value from originalLangCode to langCode based on content type.
// It uses the configured translation provider and returns the translated string or original value if translation is not available.
// Parameters:
//   - ctx: echo context
//   - fieldName: name of the field being translated
//   - value: current value to be translated
//   - originalValue: original value to translate from
//   - contentType: content type (text/html/markdown)
//   - langCode: target language code
//   - originalLangCode: source language code
//
// Returns:
//   - translated string
//   - error if translation fails
func Translate(ctx echo.Context, fieldName string, value string, originalValue string, contentType string, langCode string, originalLangCode string) (string, error) {
	cfg := GetConfig()
	if len(cfg.Provider) == 0 {
		return value, nil
	}
	trs := translate.GetProvider(cfg.Provider)
	if trs == nil {
		return value, nil
	}
	translateConfig := translate.AcquireConfig()
	translateConfig.Input = originalValue
	translateConfig.From = originalLangCode
	translateConfig.To = langCode
	translateConfig.Format = `text`
	translateConfig.APIConfig = cfg.APIConfig[cfg.Provider]
	switch contentType {
	case `html`:
		translateConfig.Format = `html`
	case `markdown`:
		translateConfig.Format = `markdown`
	}
	translateConfig.SetDefaults()
	defer translateConfig.Release()
	return trs.Translate(ctx, translateConfig)
}
