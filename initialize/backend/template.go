package backend

import (
	"github.com/coscms/webcore/library/bindata"
	_ "github.com/coscms/webfront/library/setup"
	"github.com/coscms/webfront/library/xtemplate"
)

// TmplPathFixers 后台模板文件路径修正器
var TmplPathFixers = xtemplate.New(xtemplate.KindBackend, bindata.PathAliases, true)
