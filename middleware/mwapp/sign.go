package mwapp

import (
	"bytes"
	"io"
	"net/url"

	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/engine"
)

func (a *AuthConfig) SignRequest(ctx echo.Context, appID string) (sign string, data url.Values, err error) {
	if a.secretGetter == nil {
		err = ctx.NewError(stdCode.Failure, ctx.T(`不支持获取密钥`))
		return
	}
	var secret string
	secret, err = a.secretGetter(ctx, appID)
	if err != nil {
		return
	}
	data = ctx.Forms()
	switch ctx.ResolveContentType() {
	case echo.MIMEApplicationJSON, echo.MIMEApplicationXML:
		body := ctx.Request().Body()
		var b []byte
		b, err = io.ReadAll(body)
		body.Close()
		if err != nil {
			return
		}
		ctx.Request().SetBody(io.NopCloser(bytes.NewBuffer(b)))
		if len(b) > 0 {
			data.Set(`data`, engine.Bytes2str(b))
		}
	}
	data.Del(a.FormSignKey)
	sign = a.signMaker(data, secret)
	return
}
