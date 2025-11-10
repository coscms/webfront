package initialize

import (
	"github.com/admpub/log"
	"github.com/coscms/webcore/library/config/startup"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/webx-top/echo/defaults"
)

func init() {
	startup.OnAfter(`web.installed`, func() {
		ctx := defaults.NewMockContext()
		err := i18nm.Initialize(ctx)
		if err != nil {
			log.Errorf(`[i18nm.Initialize] %v`, err)
		}
	})
}
