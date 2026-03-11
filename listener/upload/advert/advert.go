package advert

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/coscms/webcore/library/fileupdater/listener"
	"github.com/coscms/webfront/dbschema"
)

func init() {

	// - official_ad_item
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialAdItem)
		tableID = fmt.Sprint(fm.Id)
		if fm.Contype == `image` || fm.Contype == `video` || fm.Contype == `audio` {
			content = fm.Content
		}
		property = listener.NewPropertyWith(fm, db.Cond{`id`: fm.Id})
		return
	}, false).SetTable(`official_ad_item`, `content`).ListenDefault()

	// - official_ad_position
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialAdPosition)
		tableID = fmt.Sprint(fm.Id)
		if fm.Contype == `image` || fm.Contype == `video` || fm.Contype == `audio` {
			content = fm.Content
		}
		property = listener.NewPropertyWith(fm, db.Cond{`id`: fm.Id})
		return
	}, false).SetTable(`official_ad_position`, `content`).ListenDefault()

}
