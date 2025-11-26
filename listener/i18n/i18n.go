package i18n

import (
	"fmt"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/fileupdater/listener"
	"github.com/coscms/webfront/dbschema"
)

func getResourceTableAndField(ctx echo.Context, resourceId uint) (string, string) {
	resourceTable := ctx.Internal().String(`i18n_translation_resource_table`)
	resourceField := ctx.Internal().String(`i18n_translation_resource_field`)
	if len(resourceTable) > 0 && len(resourceField) == 0 {
		if rCodes, ok := ctx.Internal().Get(`i18n_translation_resource_codes`).(map[uint]string); ok {
			if code, ok := rCodes[resourceId]; ok {
				parts := strings.SplitN(code, `.`, 2)
				if len(parts) == 2 {
					resourceField = parts[1]
				}
			}
			return resourceTable, resourceField
		}
	}
	if len(resourceTable) == 0 || len(resourceField) == 0 {
		rM := dbschema.NewOfficialI18nResource(ctx)
		if err := rM.Get(nil, `id`, resourceId); err == nil {
			parts := strings.SplitN(rM.Code, `.`, 2)
			if len(parts) == 2 {
				resourceTable = parts[0]
				resourceField = parts[1]
			}
		}
	}
	return resourceTable, resourceField
}

func init() {
	// - official_i18n_translation
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialI18nTranslation)
		tableID = fmt.Sprintf(`%s:%d_%d`, fm.Lang, fm.RowId, fm.ResourceId)
		content = fm.Text
		property = listener.NewPropertyWith(fm, db.And(
			db.Cond{`lang`: fm.Lang},
			db.Cond{`row_id`: fm.RowId},
			db.Cond{`resource_id`: fm.ResourceId},
		))
		if pMap, ok := listener.UpdaterInfos[``]; ok {
			resourceTable, resourceField := getResourceTableAndField(fm.Context(), fm.ResourceId)
			if len(resourceTable) > 0 && len(resourceField) > 0 {
				if updaterInfo, ok := pMap[resourceTable]; ok && updaterInfo != nil {
					if up, ok := updaterInfo[resourceField]; ok {
						property.SetEmbedded(up.Embedded)
						return
					}
				}
			}
		}
		property.SetExit(true)
		return
	}, true).SetTable(`official_i18n_translation`, `text`).ListenDefault()
}
