package listener

import (
	"github.com/coscms/webcore/library/config/startup"
	"github.com/coscms/webfront/model/i18nm"
)

func init() {
	startup.OnAfter(`web.installed`, func() {
		i18nm.ListenTable()
	})
}
