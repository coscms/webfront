package apiutils

import (
	"github.com/coscms/sdk/sdk_options"
)

type Type = sdk_options.Type

const (
	// TypeOauth 社区登录类型
	TypeOauth Type = sdk_options.TypeOauth
	// TypePayment 支付类型
	TypePayment Type = sdk_options.TypePayment
)
