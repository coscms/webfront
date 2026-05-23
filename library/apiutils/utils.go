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
	code = sdk_options.GetValueByKey(data, `Code`, `code`)
	msg = sdk_options.GetValueByKey(data, `Info`, `info`, `msg`, `message`)
	return
}
