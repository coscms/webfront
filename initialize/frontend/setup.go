package frontend

import (
	"fmt"
	"os"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/initialize/backend"
	backendLib "github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/config"
)

var MinCustomerID = 1000

func init() {
	backendLib.OnInstalled(func(ctx echo.Context) error {
		sqlStr := fmt.Sprintf("ALTER TABLE `official_customer` AUTO_INCREMENT=%d", MinCustomerID)
		_, err := factory.NewParam().DB().ExecContext(ctx, sqlStr)
		return err
	})
}

func AutoBackendPrefix() {
	if len(config.FromCLI().BackendDomain) == 0 &&
		len(config.FromCLI().FrontendDomain) == 0 &&
		len(os.Getenv(`NGING_BACKTEND_URL_PREFIX`)) == 0 &&
		len(Prefix()) == 0 {
		backend.SetPrefix(`/admin`)
	}
}
