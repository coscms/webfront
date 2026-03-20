package customer

import (
	"github.com/coscms/webcore/library/perm"
	"github.com/coscms/webfront/library/xcommon"
	"github.com/coscms/webfront/library/xrole"
	"github.com/webx-top/echo"
)

const (
	OfflinePayBehaviorName = `offlinePay`
)

func init() {
	xrole.Behaviors.Register(OfflinePayBehaviorName, echo.T(`线下转账设置`),
		perm.BehaviorOptFormHelpBlock(echo.T(`配置线下转账提交频率。maxPerDay - 表示每天的最大发布数量(<=0代表禁止发布); maxPending - 表示待审核文章上限(<=0代表不限)`)),
		perm.BehaviorOptValue(&xcommon.ConfigCustomerAdd{}),
		perm.BehaviorOptValueInitor(func() interface{} {
			return &xcommon.ConfigCustomerAdd{}
		}),
		perm.BehaviorOptValueType(`json`),
	)
}
