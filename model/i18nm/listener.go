package i18nm

import (
	"github.com/admpub/log"
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/db/lib/factory"
)

// ListenTable listen table
func ListenTable() {
	for table, fieldInfo := range dbschema.DBI.Fields {
		var hasMultilingual bool
		for _, info := range fieldInfo {
			if !info.Multilingual {
				continue
			}
			hasMultilingual = true
			break
		}
		if !hasMultilingual {
			continue
		}
		dbschema.DBI.On(factory.EventDeleting, func(m factory.Model, _ ...string) error {
			id := GetRowID(m)
			if id == 0 {
				return nil
			}
			err := DeleteModelTranslations(m.Context(), m, id)
			return err
		}, table)
		log.Infof(`[i18nm.ListenTable] %v`, table)
	}
}
