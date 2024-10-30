package initialize

import (
	"github.com/coscms/webcore/library/config"

	"github.com/coscms/webfront/library/search/segment"
)

func init() {
	config.AddConfigInitor(func(c *config.Config) {
		segment.ApplySegmentConfig(c)
	})
}
