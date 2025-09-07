package underattack

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

type Rule struct {
	Path        string
	UAWhitelist string
	Headers     map[string]string
	QueryString string // a=b&b=c
	regexpUA    *regexp.Regexp
	kvURIQuery  map[string][]string
}

func (r *Rule) init() {
	if len(r.UAWhitelist) > 0 {
		rows := param.Unique(com.TrimSpaceForRows(r.UAWhitelist))
		if len(rows) > 0 {
			r.regexpUA = regexp.MustCompile(strings.Join(rows, `|`))
		}
	}
	r.kvURIQuery = map[string][]string{}
	if len(r.QueryString) > 0 {
		r.kvURIQuery, _ = url.ParseQuery(r.QueryString)
	}
}

func (r Rule) Validate(ctx echo.Context) error {
	var err error
	if len(r.UAWhitelist) > 0 {
		rows := com.TrimSpaceForRows(r.UAWhitelist)
		for _, row := range rows {
			_, err = regexp.Compile(row)
			if err != nil {
				return ctx.NewError(code.InvalidParameter, `正则表达式语法错误: %s`, row)
			}
		}
	}
	return err
}

// IsAllowed 没有任何白名单规则时直接放行; 全部白名单规则匹配时直接放行
func (r Rule) IsAllowed(ctx echo.Context) bool {
	if len(r.Headers) > 0 {
		for k, v := range r.Headers {
			if ctx.Header(k) != v {
				return false
			}
		}
	}
	if r.kvURIQuery != nil {
		for key, values := range r.kvURIQuery {
			inputs := ctx.FormValues(key)
			for _, value := range values {
				if !com.InSlice(value, inputs) {
					return false
				}
			}
		}
	}
	if r.regexpUA != nil {
		if !r.regexpUA.MatchString(ctx.Request().UserAgent()) {
			return false
		}
	}
	return true
}
