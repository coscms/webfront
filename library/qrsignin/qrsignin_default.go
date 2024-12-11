package qrsignin

import "github.com/webx-top/echo"

func newDefaultQRSignIn() defaultQRSignIn {
	c := defaultQRSignIn{}
	return c
}

type defaultQRSignIn struct {
}

func (c defaultQRSignIn) Encode(_ echo.Context, signInData QRSignIn) (string, error) {
	return signInData.Encode()
}

func (c defaultQRSignIn) Decode(ctx echo.Context, encrypted string) (QRSignIn, error) {
	signInData := QRSignIn{}
	err := signInData.Decode(ctx, encrypted)
	return signInData, err
}
