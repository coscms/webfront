package i18nm

import (
	"errors"

	"github.com/coscms/webcore/library/config"
)

// Config 配置
type Config struct {
	Providers           []ProviderConfig `json:"providers"`           // 翻译提供商
	On                  bool             `json:"on"`                  // 是否开启翻译
	AllowForceTranslate bool             `json:"allowForceTranslate"` // 是否允许强制翻译
}

var ErrTranslationOff = errors.New(`translation is turned off`)
var ErrNoTranslationProviders = errors.New(`no translation providers configured`)

// Check validates the translation configuration and returns an error if the configuration is invalid.
// It checks if translation is enabled and if there are any translation providers configured.
func (c *Config) Check() error {
	if !c.On {
		return ErrTranslationOff
	}
	if len(c.Providers) == 0 {
		return ErrNoTranslationProviders
	}
	return nil
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
