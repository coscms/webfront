package qrsignin

import (
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	modelCustomer "github.com/coscms/webfront/model/official/customer"
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
