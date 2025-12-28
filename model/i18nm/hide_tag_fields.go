package i18nm

import (
	"slices"

	"github.com/coscms/webfront/library/top"
)

var HideTagFields = map[string][]string{
	`official_common_article`: []string{`content`},
}

func RegisterHideTagFields(table string, fields ...string) {
	if len(fields) == 0 {
		return
	}
	if _, ok := HideTagFields[table]; !ok {
		HideTagFields[table] = []string{}
	}
	HideTagFields[table] = append(HideTagFields[table], fields...)
}

func DefaultOriginalTextPickout(table string, fieldName string, originalText string) (string, func(translatedText string) string) {
	fields, ok := HideTagFields[table]
	if !ok {
		return originalText, nil
	}
	ok = slices.Contains(fields, fieldName)
	if !ok {
		return originalText, nil
	}
	var picks []string
	picks, originalText = top.PickoutHideTag(originalText)
	return originalText, func(translatedText string) string {
		translatedText = top.RestorePickoutedHideTag(translatedText, picks)
		return translatedText
	}
}
