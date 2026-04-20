package i18nm

import (
	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/echo"
)

// Model is an interface for model instances.
type Model interface {
	Short_() string
	GetField(string) interface{}
	FromRow(map[string]interface{})
}

// GetRowID returns the ID of a model instance.
func GetRowID(mdl Model) uint64 {
	var id uint64
	switch v := mdl.GetField(`Id`).(type) {
	case uint:
		id = uint64(v)
	case uint8:
		id = uint64(v)
	case uint16:
		id = uint64(v)
	case uint32:
		id = uint64(v)
	case uint64:
		id = v
	default:
	}
	return id
}

func translateText(ctx echo.Context, contentType string, translate Translator, restoreFunc func(string) string,
	forceTranslate bool, autoTranslate bool, field string, originalText string, translatedText string,
	langCode string, langDefault string) (string, error) {
	if len(contentType) == 0 {
		contentType = `string` // 默认string类型(单行文本)
	}
	if !common.CanTranslateContent(contentType) {
		return translatedText, nil
	}
	var err error
	if forceTranslate {
		translatedText, err = translate(ctx, field, translatedText, originalText, contentType, langCode, langDefault)
	} else if len(translatedText) == 0 && autoTranslate {
		translatedText, err = translate(ctx, field, translatedText, originalText, contentType, langCode, langDefault)
	} else {
		return translatedText, nil
	}
	if err != nil {
		return translatedText, err
	}
	if restoreFunc != nil {
		translatedText = restoreFunc(translatedText)
	}
	return translatedText, nil
}
