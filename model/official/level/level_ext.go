package level

import (
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/echo"
)

const (
	// 等级状态
	LevelStatusActived = `actived`
	LevelStatusExpired = `expired`

	// 金额类型
	AmountTypeBalance     = `balance`
	AmountTypeAccumulated = `accumulated`
)

var LevelStatuses = echo.NewKVData().Add(LevelStatusActived, echo.T(`有效`)).Add(LevelStatusExpired, echo.T(`过期`))
var AmountTypes = echo.NewKVData().Add(AmountTypeBalance, echo.T(`余额`)).Add(AmountTypeAccumulated, echo.T(`累积总收入`))

type LevelGroup struct {
	Group string
	Title string
	List  []*dbschema.OfficialCustomerLevel
}

type RelationExt struct {
	*dbschema.OfficialCustomerLevelRelation
	Level *dbschema.OfficialCustomerLevel `db:"-,relation=id:level_id|gtZero"`
}

func (r *RelationExt) Name_() string {
	if r.OfficialCustomerLevelRelation == nil {
		r.OfficialCustomerLevelRelation = &dbschema.OfficialCustomerLevelRelation{}
	}
	return r.OfficialCustomerLevelRelation.Name_()
}
