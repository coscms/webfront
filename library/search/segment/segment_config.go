package segment

import (
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/extend"
)

func init() {
	extend.Register(`segment`, func() interface{} {
		return &Config{}
	})
	config.OnKeySetSettings(`thirdparty.segment`, func(diff config.Diff) error {
		if !diff.IsDiff {
			return nil
		}
		ApplySegmentConfig(nil)
		return nil
	})
}

var _ extend.Reloader = (*Config)(nil)
var _ config.ReloadByConfig = (*Config)(nil)

type Config struct {
	Engine string `json:"engine"`
	ApiURL string `json:"apiURL"`
	ApiKey string `json:"apiKey"`
}

func (c *Config) Reload() error {
	ApplySegmentConfig(nil)
	return nil
}

func (c *Config) ReloadByConfig(newCfg *config.Config, args ...string) error {
	ApplySegmentConfig(newCfg)
	return nil
}

func GetConfig(c *config.Config) (cfg *Config, ok bool) {
	cfg, ok = config.FromDB(`thirdparty`).Get(`segment`).(*Config)
	if ok && cfg.Engine != `` {
		return
	}
	if c == nil {
		c = config.FromFile()
	}
	cfg, ok = c.Extend.Get(`segment`).(*Config)
	return
}
