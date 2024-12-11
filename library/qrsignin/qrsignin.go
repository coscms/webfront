package qrsignin

import (
	"time"

	"github.com/coscms/webcore/library/config"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

type QRSignIn struct {
	SessionID     string `json:"sID"`
	SessionMaxAge int    `json:"sAge"`
	Expires       int64  `json:"dExp"` // 数据过期时间戳(秒)
	IPAddress     string `json:"ip"`
	UserAgent     string `json:"ua"`
	Platform      string `json:"pf"`
	Scense        string `json:"ss"`
	DeviceNo      string `json:"dn"`
}

func (q QRSignIn) Encode() (string, error) {
	plaintext, err := com.JSONEncodeToString(q)
	if err != nil {
		return ``, err
	}
	qrcode := config.FromFile().Encode256(plaintext)
	qrcode = com.URLSafeBase64(qrcode, true)
	return qrcode, nil
}

func (q *QRSignIn) Decode(ctx echo.Context, encrypted string) error {
	encrypted = com.URLSafeBase64(encrypted, false)
	plaintext := config.FromFile().Decode256(encrypted)
	if len(plaintext) == 0 {
		return ctx.NewError(code.InvalidParameter, `解密失败`).SetZone(`data`)
	}
	return com.JSONDecodeString(plaintext, q)
}

func NewQRSignIn(ctx echo.Context, cookieMaxAge int, expireTime time.Time) QRSignIn {
	qsi := QRSignIn{
		SessionID:     ctx.Session().MustID(),
		SessionMaxAge: cookieMaxAge,
		Expires:       expireTime.Unix(),
		IPAddress:     ctx.RealIP(),
		UserAgent:     ctx.Request().UserAgent(),
		Platform:      ctx.Header(`X-Platform`),
		Scense:        ctx.Header(`X-Scense`),
		DeviceNo:      ctx.Header(`X-Device-Id`),
	}
	if len(qsi.Scense) > 0 {
		qsi.Scense = `qrcode_` + qsi.Scense
	} else {
		qsi.Scense = `qrcode_` + modelCustomer.DefaultDeviceScense
	}
	return qsi
}

func GenerateUniqueKey(ip, ua string) string {
	return com.Md5(com.String(time.Now().UnixMicro())+ua+ip) + com.RandomAlphanumeric(2)
}
