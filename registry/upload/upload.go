package upload

import (
	"github.com/coscms/webcore/registry/upload"
	"github.com/webx-top/echo"
)

func init() {
	upload.Subdir.Add(`category`, echo.T(`分类`))
	upload.Subdir.Add(`navigate`, echo.T(`导航`))
	upload.Subdir.Add(`friendlink`, echo.T(`友情链接`))
	upload.Subdir.Add(`membership`, echo.T(`会员套餐`))
}
