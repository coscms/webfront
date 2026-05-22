package apiutils

import (
	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/com"
)

func Decrypt(decrypted string, secret string) string {
	raw := decrypted
	if len(raw) > 0 {
		raw = com.URLSafeBase64(raw, false)
	}
	common.Decrypt(secret, &raw)
	return raw
}

func Encrypt(raw string, secret string) string {
	encrypted := raw
	common.Encrypt(secret, &encrypted)
	if len(encrypted) > 0 {
		encrypted = com.URLSafeBase64(encrypted, true)
	}
	return encrypted
}
