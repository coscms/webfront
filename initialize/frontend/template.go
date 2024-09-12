package frontend

import (
	"github.com/coscms/webcore/library/ntemplate"
	_ "github.com/coscms/webfront/library/formbuilder"
	"github.com/coscms/webfront/library/xtemplate"
)

// TmplPathFixers 前台模板文件路径修正器
var TmplPathFixers = xtemplate.New(xtemplate.KindFrontend, ntemplate.NewPathAliases(), true)
