package category

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/coscms/webcore/library/fileupdater/listener"
	"github.com/coscms/webfront/dbschema"
)

func init() {

	// - official_common_navigate
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialCommonNavigate)
		tableID = fmt.Sprint(fm.Id)
		content = fm.Cover
		property = listener.NewPropertyWith(fm, db.Cond{`id`: fm.Id})
		return
	}, false).SetTable(`official_common_navigate`, `cover`).ListenDefault()
}
