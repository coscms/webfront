package customer

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/coscms/webcore/library/fileupdater/listener"
	"github.com/coscms/webfront/dbschema"
)

func init() {
	// - official_customer_level
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialCustomerLevel)
		tableID = fmt.Sprint(fm.Id)
		content = fm.IconImage
		property = listener.NewPropertyWith(fm, db.Cond{`id`: fm.Id})
		return
	}, false).SetTable(`official_customer_level`, `icon_image`).ListenDefault()

	// - official_customer_group_package
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialCustomerGroupPackage)
		tableID = fmt.Sprint(fm.Id)
		content = fm.IconImage
		property = listener.NewPropertyWith(fm, db.Cond{`id`: fm.Id})
		return
	}, false).SetTable(`official_customer_group_package`, `icon_image`).ListenDefault()
}
