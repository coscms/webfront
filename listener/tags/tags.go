package tags

import (
	"github.com/webx-top/db/lib/factory"

	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/model/i18nm"
)

func init() {

	// - official_common_tags
	dbschema.DBI.On(factory.EventDeleting, func(m factory.Model, _ ...string) error {
		fm := m.(*dbschema.OfficialCommonTags)
		err := i18nm.DeleteModelTranslations(m, uint64(fm.Id))
		return err
	}, `official_common_tags`)
}
