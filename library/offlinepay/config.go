package offlinepay

import "github.com/coscms/webcore/library/config"

var configDefault = &Config{}

const configKey = `offlinePay`

type MethodOptions struct {
	Disabled *bool   `json:"disabled"`
	Logo     *string `json:"logo"`
}

type Config struct {
	Methods map[string]MethodOptions `json:"methods"`
}

func GetConfig() *Config {
	cfg, ok := config.FromFile().Extend.Get(configKey).(*Config)
	if !ok {
		cfg = configDefault
	}
	return cfg
}
