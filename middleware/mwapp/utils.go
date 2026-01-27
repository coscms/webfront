package mwapp

import (
	"strings"
	"time"

	"github.com/coscms/webfront/library/xcode"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// GetAppID returns the AppID of the current request.
// It first checks the form data, and if not found, it checks the header.
// The AppID is trimmed of whitespace.
func GetAppID(ctx echo.Context, a *AuthConfig) string {
	appID := ctx.Formx(a.FormAppIDKey).String()
	if len(appID) == 0 {
		appID = ctx.Header(a.HeaderAppIDKey)
		appID = strings.TrimSpace(appID)
	}
	return appID
}

// GetSign returns the sign of the current request.
// It first checks the form data, and if not found, it checks the header.
// The sign is trimmed of whitespace.
func GetSign(ctx echo.Context, a *AuthConfig) string {
	sign := ctx.Formx(a.FormSignKey).String()
	if len(sign) == 0 {
		sign = ctx.Header(a.HeaderSignKey)
		sign = strings.TrimSpace(sign)
	}
	return sign
}

// GetTimestamp returns the timestamp of the current request.
// It first checks the form data, and if not found, it checks the header.
// The timestamp is parsed as an int64.
// If the timestamp is not found or is invalid, it returns 0.
//
// The timestamp is expected to be in seconds since the Unix epoch (January 1, 1970, 00:00:00 UTC).
//
// The timestamp is used to verify the lifetime of the request.
// If the request is older than the lifetime specified in the AuthConfig,
// the request is rejected.
func GetTimestamp(ctx echo.Context, a *AuthConfig) int64 {
	timestamp := ctx.Formx(a.FormTimeKey).Int64()
	if timestamp <= 0 {
		ts := ctx.Header(a.HeaderTimeKey)
		timestamp = param.AsInt64(ts)
	}
	return timestamp
}

// VerifyLifetime verifies the lifetime of the current request.
// If the request has a sign and the lifetime is specified in the AuthConfig,
// it checks if the request is older than the lifetime.
// If the request is older than the lifetime, it returns an error.
// The error is set to be in the zone of the FormTimeKey of the AuthConfig.
func VerifyLifetime(ctx echo.Context, a *AuthConfig, timestamp int64, sign string) error {
	if len(sign) > 0 && a.LifeSeconds > 0 {
		if time.Now().Unix()-timestamp > a.LifeSeconds {
			return ctx.NewError(xcode.SignatureHasExpired, ctx.T(`页面已经失效，请返回重新提交`)).SetZone(a.FormTimeKey)
		}
	}
	return nil
}
