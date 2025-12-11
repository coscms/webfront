package initialize

import (
	"github.com/admpub/events"
	"github.com/admpub/log"
	"github.com/coscms/webcore/library/config/startup"
	"github.com/coscms/webfront/model/i18nm"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func init() {
	//onWebInstalled()
	echo.OnCallback(`nging.upgrade.db.after`, func(e events.Event) error {
		initI18Resources()
		return nil
	})
}

func onWebInstalled() {
	startup.OnAfter(`web.installed`, initI18Resources)
}

func initI18Resources() {
	ctx := defaults.NewMockContext()
	err := i18nm.Initialize(ctx)
	if err != nil {
		log.Errorf(`[i18nm.Initialize] %v`, err)
	}
}
