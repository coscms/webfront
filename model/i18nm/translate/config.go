package translate

import "github.com/coscms/webcore/library/config"

// Config 配置
type Config struct {
	Providers           []ProviderConfig `json:"providers"`           // 翻译提供商
	On                  bool             `json:"on"`                  // 是否开启翻译
	AllowForceTranslate bool             `json:"allowForceTranslate"` // 是否允许强制翻译
}

// ProviderConfig 配置
type ProviderConfig struct {
	Provider string            `json:"provider"` // 翻译提供商名称
	Config   map[string]string `json:"config"`   // 翻译提供商配置
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
//					provider  : "tencent"
//					config : {
//						appid  : ""
//						secret : ""
//					}
//				}
//			]
//			on : true
//			allowForceTranslate : true
//		}
//	}
func GetConfig() *Config {
	cfg, ok := config.FromFile().Extend.Get(`translate`).(*Config)
	if !ok {
		return NewConfig()
	}
	return cfg
}
