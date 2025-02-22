package initialize

import (
	_ "github.com/coscms/webcore"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/navigate"
	_ "github.com/coscms/webfront/library/minify"
	"github.com/webx-top/echo"
)

var nav = &navigate.List{}

var Project = navigate.NewProject(echo.T(`#short#内容管理`), `webx`, `/official/customer/index`, nav)

func init() {
	bootconfig.OfficialHomepage = `https://www.coscms.com`
	httpserver.Backend.Navigate.Projects.Get(`nging`).Name = echo.T(`#short#其它功能`)
	httpserver.Backend.Navigate.AddProject(1, Project)
}
