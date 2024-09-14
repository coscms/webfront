package initialize

import (
	_ "github.com/coscms/webcore"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/registry/navigate"
	_ "github.com/coscms/webfront/initialize/backend"
)

var nav = &navigate.List{}

var Project = navigate.NewProject(`内容管理`, `webx`, `/official/customer/index`, nav)

func init() {
	bootconfig.OfficialHomepage = `https://www.coscms.com`
	navigate.ProjectGet(`nging`).Name = `其它功能`
	navigate.ProjectAdd(1, Project)
}
