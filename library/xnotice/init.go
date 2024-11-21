package xnotice

import (
	"github.com/admpub/log"
	"github.com/webx-top/db"
	"github.com/webx-top/echo/defaults"

	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/initialize/frontend"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func init() {
	frontend.Notify.OnOpen(func(user string) {
		ctx := defaults.NewMockContext()
		custM := dbschema.NewOfficialCustomer(ctx)
		err := custM.UpdateField(nil, `online`, `Y`, db.And(
			db.Cond{`name`: user},
			db.Cond{`online`: `N`},
		))
		if err != nil {
			log.Errorf(`failed to custM.UpdateField(online=Y,name=%q): %v`, user, err)
		}
	})
	frontend.Notify.OnClose(func(user string) {
		ctx := defaults.NewMockContext()
		userM := dbschema.NewOfficialCustomer(ctx)
		err := userM.UpdateField(nil, `online`, `N`, db.And(
			db.Cond{`name`: user},
			db.Cond{`online`: `Y`},
		))
		if err != nil {
			log.Errorf(`failed to custM.UpdateField(online=N,name=%q): %v`, user, err)
		}
	})
	modelCustomer.OnSignOut(onLogout)
}

func onLogout(customer *dbschema.OfficialCustomer) error {
	frontend.Notify.CloseMessage(customer.Name)
	return nil
}
