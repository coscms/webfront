package translate

import "github.com/coscms/webcore/library/config"

// Config 配置
type Config struct {
	Provider  string                       `json:"provider"`
	APIConfig map[string]map[string]string `json:"apiConfig"`
}

// NewConfig creates and returns a new Config instance with initialized APIConfig map
func NewConfig() *Config {
	return &Config{
		APIConfig: map[string]map[string]string{},
	}
}

// GetConfig returns the translation configuration from the extended config.
// If no translation config exists, returns a new default configuration.
func GetConfig() *Config {
	cfg, ok := config.FromFile().Extend.Get(`translate`).(*Config)
	if !ok {
		return NewConfig()
	}
	return cfg
}
