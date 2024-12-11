package qrsignin

import "github.com/webx-top/echo"

type QRSignInCodec interface {
	Encode(echo.Context, QRSignIn) (string, error)
	Decode(echo.Context, string) (QRSignIn, error)
}

var qrSignInCodecs = map[string]QRSignInCodec{
	`cache`:   newCacheQRSignIn(),   // 缺点：占用存储空间；优点：字符串短，生成的二维码更容易识别
	`default`: newDefaultQRSignIn(), // 优点：不占用存储空间；缺点：加密字符串太长，生成的二维码元素图块小而多，不易识别
}

func Get(caseName string) QRSignInCodec {
	cs, ok := qrSignInCodecs[caseName]
	if ok {
		return cs
	}
	return qrSignInCodecs[`default`]
}

func Register(caseName string, qrsic QRSignInCodec) {
	qrSignInCodecs[caseName] = qrsic
}
