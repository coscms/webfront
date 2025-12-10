package translate

import "github.com/coscms/webcore/library/config"

// Config 配置
type Config struct {
	Providers []ProviderConfig `json:"providers"`
}

// ProviderConfig 配置
type ProviderConfig struct {
	Provider string            `json:"provider"`
	Config   map[string]string `json:"config"`
}

// NewConfig creates and returns a new Config instance with initialized APIConfig map
func NewConfig() *Config {
	return &Config{
		Providers: []ProviderConfig{},
	}
}

// GetConfig returns the translation configuration from the extended config.
// If no translation config exists, returns a new default configuration.
//
//	extend : {
//		translate : {
//			providers : [
//				{
//					provider  : "tencent",
//					config : {
//						"appid": "",
//						"secret": ""
//					}
//				}
//			]
//		}
//	}
func GetConfig() *Config {
	cfg, ok := config.FromFile().Extend.Get(`translate`).(*Config)
	if !ok {
		return NewConfig()
	}
	return cfg
}
