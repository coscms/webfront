package initialize

import (
	"github.com/coscms/webcore/library/config"

	"github.com/coscms/webfront/library/search/segment"
)

// init registers segment configuration initializer with the global config system.
// It applies segment-specific configuration when the config system initializes.
func init() {
	config.AddConfigInitor(func(c *config.Config) {
		segment.ApplySegmentConfig(c)
	})
}
