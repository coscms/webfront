package apiutils

import (
	"github.com/coscms/sdk/sdk_options"
)

type Type = sdk_options.Type

const (
	// TypeOAuth 社区登录类型
	TypeOAuth Type = sdk_options.TypeOAuth
	// TypePayment 支付类型
	TypePayment Type = sdk_options.TypePayment
)
