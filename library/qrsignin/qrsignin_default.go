package qrsignin

import (
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func newDefaultQRSignIn() defaultQRSignIn {
	c := defaultQRSignIn{}
	return c
}

type defaultQRSignIn struct {
}

func (c defaultQRSignIn) Encode(_ echo.Context, signInData QRSignIn) (string, error) {
	plaintext, err := com.JSONEncodeToString(signInData)
	if err != nil {
		return ``, err
	}
	qrcode := config.FromFile().Encode256(plaintext)
	qrcode = com.URLSafeBase64(qrcode, true)
	return qrcode, nil
}

func (c defaultQRSignIn) Decode(ctx echo.Context, encrypted string) (QRSignIn, error) {
	signInData := QRSignIn{}
	encrypted = com.URLSafeBase64(encrypted, false)
	plaintext := config.FromFile().Decode256(encrypted)
	var err error
	if len(plaintext) == 0 {
		err = ctx.NewError(code.InvalidParameter, `解密失败`).SetZone(`data`)
	} else {
		err = com.JSONDecodeString(plaintext, &signInData)
	}
	return signInData, err
}
