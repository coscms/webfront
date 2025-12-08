package translate

import "github.com/coscms/webcore/library/config"

type Config struct {
	Provider  string                       `json:"provider"`
	APIConfig map[string]map[string]string `json:"apiConfig"`
}

func NewConfig() *Config {
	return &Config{
		APIConfig: map[string]map[string]string{},
	}
}

func GetConfig() *Config {
	cfg, ok := config.FromFile().Extend.Get(`translate`).(*Config)
	if !ok {
		return NewConfig()
	}
	return cfg
}
