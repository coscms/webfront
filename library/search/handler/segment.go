package handler

import (
	"strings"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webfront/library/search/segment"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func RegisterRoute(e echo.RouteRegister) {
	e.Post(`/segment`, segmentHandler, checkToken).SetMetaKV(`noAttack`, true)
}

func segmentHandler(ctx echo.Context) error {
	text := ctx.Form(`text`)
	if len(text) == 0 {
		return ctx.String(``)
	}
	args := ctx.Formx(`args`).Split(`,`).String()
	mode := ctx.Form(`mode`)
	var result []string
	if len(mode) > 0 {
		result = segment.Default().SegmentBy(text, mode, args...)
	} else {
		result = segment.Default().Segment(text, args...)
	}
	return ctx.String(strings.Join(result, ` `))
}

func checkToken(h echo.Handler) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		apiKey := config.Setting(`base`).String(`apiKey`)
		if len(apiKey) == 0 {
			return echo.ErrNotAcceptable
		}
		token := ctx.Header(`Authorization`)
		parts := strings.SplitN(token, ` `, 2)
		if len(parts) == 2 {
			token = parts[1]
		} else {
			token = parts[0]
		}

		text := ctx.Form(`text`)
		args := ctx.Form(`args`)
		mode := ctx.Form(`mode`)
		data := map[string]string{
			`text`: text,
		}
		if len(mode) > 0 {
			data[`mode`] = mode
		}
		if len(args) > 0 {
			data[`args`] = args
		}
		b, _ := com.JSONEncode(data)
		expectedToken := com.Token(apiKey, b)
		if expectedToken != token {
			return echo.ErrBadRequest
		}
		return h.Handle(ctx)
	}
}
