package offlinepay

import "github.com/coscms/webcore/library/config/extend"

func init() {
	extend.Register(configKey, func() interface{} {
		return &Config{
			Methods: map[string]MethodOptions{},
		}
	})
}
