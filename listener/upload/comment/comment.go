package comment

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/coscms/webcore/library/fileupdater/listener"
	"github.com/coscms/webfront/dbschema"
)

func init() {
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialCommonComment)
		tableID = fmt.Sprint(fm.Id)
		content = fm.Content
		property = listener.NewPropertyWith(fm, db.Cond{`id`: fm.Id})
		return
	}, true).SetTable(`official_common_comment`, `content`).ListenDefault()
}
