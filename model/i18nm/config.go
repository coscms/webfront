package i18nm

import (
	"errors"
	"slices"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/extend"
	"github.com/coscms/webfront/library/xkv"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var _ extend.SetDefaults = (*Config)(nil)

// Config 配置
type Config struct {
	Providers           []ProviderConfig `json:"providers"`           // 翻译提供商
	On                  bool             `json:"on"`                  // 是否开启翻译
	AllowForceTranslate bool             `json:"allowForceTranslate"` // 是否允许强制翻译
}

func (c *Config) SetDefaults() {
	c.Providers = slices.DeleteFunc(c.Providers, func(pc ProviderConfig) bool {
		return len(pc.Provider) == 0
	})
}

func (c *Config) FromStore(r echo.H) *Config {
	prov := ProviderConfig{
		Provider: r.String(`provider`),
		Config:   com.SplitKVRows(r.String(`config`)),
	}
	c.Providers = []ProviderConfig{prov}
	c.AllowForceTranslate = r.Bool(`allowForceTranslate`)
	c.On = r.Bool(`on`)
	return c
}

func (c *Config) IsValid() bool {
	switch len(c.Providers) {
	case 0:
		return false
	default:
		if len(c.Providers[0].Provider) == 0 {
			return false
		}
	}
	return true
}

func (c *Config) MergeProviders(from *Config) *Config {
	var size int
	validA := c.IsValid()
	if validA {
		size = len(c.Providers)
	}
	validB := from.IsValid()
	if validB {
		size += len(from.Providers)
	}
	cloned := Config{
		Providers:           make([]ProviderConfig, size),
		On:                  c.On,
		AllowForceTranslate: c.AllowForceTranslate,
	}
	if validA {
		n := copy(cloned.Providers, c.Providers)
		if validB {
			copy(cloned.Providers[n:], from.Providers)
		}
	} else if validB {
		copy(cloned.Providers, from.Providers)
	}
	return &cloned
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
func GetConfig(ctx echo.Context) *Config {
	if ctx == nil {
		return getConfig()
	}
	cfg, _ := xkv.GetOnce(ctx, `translate.config`, func() (*Config, error) {
		return getConfig(), nil
	})
	return cfg
}

func getConfig() *Config {
	cfgDB, okDB := config.FromDB(`thirdparty`).Get(`translate`).(*Config)
	cfgFile, okFile := config.FromFile().Extend.Get(`translate`).(*Config)
	if !okFile {
		if okDB {
			return cfgDB
		}
		return NewConfig()
	}
	if okDB {
		return cfgDB.MergeProviders(cfgFile)
	}
	return cfgFile
}
