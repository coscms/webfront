package comment

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/coscms/webcore/library/fileupdater/listener"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/model/official"
	"github.com/coscms/webfront/model/official/comment"
)

func init() {
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.OfficialCommonComment)
		tableID = fmt.Sprint(fm.Id)
		content = fm.Content
		property = listener.NewPropertyWith(fm, db.Cond{`id`: fm.Id})
		return
	}, true).SetTable(`official_common_comment`, `content`).ListenDefault()

	dbschema.DBI.On(factory.EventDeleting, func(m factory.Model, editColumns ...string) error {
		fm := m.(*dbschema.OfficialCommonComment)
		if fm.ReplyCommentId > 0 {
			cmtM := comment.NewComment(fm.Context())
			err := cmtM.DecrRepliesBy(fm)
			if err != nil {
				return err
			}
		}
		typeCfg, ok := comment.CommentAllowTypes[fm.TargetType]
		if ok && typeCfg.AfterDelete != nil {
			if err := typeCfg.AfterDelete(fm.Context(), fm); err != nil {
				return err
			}
		}
		flowM := official.NewClickFlow(fm.Context())
		err := flowM.DelByTarget(`comment`, fm.Id)
		return err
	}, `official_common_comment`)
}
