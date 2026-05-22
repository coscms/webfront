package apiutils

import (
	"github.com/coscms/sdk/sdk_options"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func AsRawData(resp *sdk_options.Response) *echo.RawData {
	return &echo.RawData{
		Code:  code.Code(resp.Code),
		State: resp.State,
		Data:  resp.Data,
		Info:  resp.Info,
		URL:   resp.URL,
		Zone:  resp.Zone,
	}
}

func MakeFormValueGetter(ctx echo.Context) func(string) string {
	return func(name string) string {
		return ctx.Form(name)
	}
}

func GetCodeAndMsgFromAPIData(data map[string]any) (code, msg any) {
	for _, key := range []string{`Code`, `code`} {
		var ok bool
		code, ok = data[key]
		if ok {
			break
		}
	}
	for _, key := range []string{`Info`, `info`, `msg`, `message`} {
		var ok bool
		msg, ok = data[key]
		if ok {
			break
		}
	}
	return
}
