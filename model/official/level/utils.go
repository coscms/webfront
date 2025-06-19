package level

import (
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/db"
)

func MakeFreeCond(group string, balance float64, accumulated float64, asset string) *db.Compounds {
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{`disabled`: `N`},
		db.Cond{`group`: group},
		db.Cond{`price`: 0},
		db.Cond{`integral_asset`: asset},
	)
	cond.Add(
		db.Or(
			db.And(
				db.Cond{`integral_amount_type`: AmountTypeBalance},
				db.Cond{`integral_min`: db.Lte(balance)},
				db.Cond{`integral_max`: db.Gte(balance)},
			),
			db.And(
				db.Cond{`integral_amount_type`: AmountTypeAccumulated},
				db.Cond{`integral_min`: db.Lte(accumulated)},
				db.Cond{`integral_max`: db.Gte(accumulated)},
			),
		),
	)
	return cond
}

func MakePayCond(group string, balance float64, accumulated float64, asset string) *db.Compounds {
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{`disabled`: `N`},
		db.Cond{`group`: group},
		db.Cond{`price`: db.Gt(0)},
		db.Cond{`integral_asset`: asset},
	)
	cond.Add(
		db.Or(
			db.And(
				db.Cond{`integral_amount_type`: AmountTypeBalance},
				db.Cond{`integral_min`: db.Lte(balance)},
				db.Cond{`integral_max`: db.Gte(balance)},
			),
			db.And(
				db.Cond{`integral_amount_type`: AmountTypeAccumulated},
				db.Cond{`integral_min`: db.Lte(accumulated)},
				db.Cond{`integral_max`: db.Gte(accumulated)},
			),
		),
	)
	return cond
}

func compareGreaterOrEqual(a *dbschema.OfficialCustomerLevel, b *dbschema.OfficialCustomerLevel) bool {
	return a.Score >= b.Score
}

func compareLess(a *dbschema.OfficialCustomerLevel, b *dbschema.OfficialCustomerLevel) bool {
	return a.Score < b.Score
}
