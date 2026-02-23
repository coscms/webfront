package version

import (
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/version"
)

const (
	PkgName = `webfront`
	// 当前应用数据表结构版本号
	dbschema = 2.2
	// 数据表结构最终版本号
	DBSCHEMA = dbschema + version.DBSCHEMA
)

func init() {
	config.Version.SetPkgDBSchemas(PkgName, dbschema)
}
