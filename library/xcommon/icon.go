package xcommon

import (
	"regexp"
	"slices"
	"strings"
)

func ParseIconFromCSS(content string, prefixes ...string) ([]string, error) {
	if len(prefixes) == 0 {
		prefixes = append(prefixes, `icon`, `fa`, `glyphicon`)
	} else {
		for i, v := range prefixes {
			prefixes[i] = regexp.QuoteMeta(v)
		}
	}

	r, err := regexp.Compile(`\.((?:` + strings.Join(prefixes, `|`) + `)-(?:[\w-]+))`)
	if err != nil {
		return nil, err
	}
	matches := r.FindAllStringSubmatch(content, -1)
	var results []string
	for _, match := range matches {
		if slices.Contains(results, match[1]) {
			continue
		}
		results = append(results, match[1])
	}
	return results, err
}
