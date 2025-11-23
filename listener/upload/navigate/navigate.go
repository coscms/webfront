package category

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/coscms/webcore/library/fileupdater/listener"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/model/i18nm"
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

	dbschema.DBI.On(factory.EventDeleting, func(m factory.Model, _ ...string) error {
		fm := m.(*dbschema.OfficialCommonNavigate)
		err := i18nm.DeleteModelTranslations(m, uint64(fm.Id))
		return err
	}, `official_common_navigate`)
}
