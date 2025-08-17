package xhtml

import (
	"strings"

	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
)

func GetLangCodeByPath(path string) string {
	langCode := config.FromFile().Language.Default
	lcInPath := strings.SplitN(strings.TrimPrefix(path, `/`), `/`, 2)[0]
	if langCode != lcInPath &&
		com.InSlice(lcInPath, config.FromFile().Language.AllList) {
		langCode = lcInPath
	}
	return langCode
}
