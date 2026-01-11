package upgrade

import (
	"strings"

	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func init() {
	type Role struct {
		Id           uint   `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
		PermCmd      string `db:"perm_cmd" bson:"perm_cmd" comment:"快捷命令权限(多个用“,”隔开)" json:"perm_cmd" xml:"perm_cmd"`
		PermAction   string `db:"perm_action" bson:"perm_action" comment:"操作权限(多个用“,”隔开)" json:"perm_action" xml:"perm_action"`
		PermBehavior string `db:"perm_behavior" bson:"perm_behavior" comment:"行为权限(多个用“,”隔开)" json:"perm_behavior" xml:"perm_behavior"`
	}
	echo.OnCallback(`nging.upgrade.db.before`, func(data echo.Event) error {
		installedSchemaVer := data.Context.Float64(`installedSchemaVer`)
		if installedSchemaVer > 8.6 {
			return nil
		}
		return upgradeArticleTagsData()
	})
}

func upgradeArticleTagsData() error {
	ctx := defaults.NewMockContext()
	m := dbschema.NewOfficialCommonArticle(ctx)
	var err error

	for {
		_, err = m.ListByOffset(nil, func(r db.Result) db.Result {
			return r.Select(`id`, `tags`)
		}, 0, 200, db.Raw("NOT JSON_VALID(`tags`)"))
		if err != nil {
			return err
		}
		rows := m.Objects()
		if len(rows) == 0 {
			break
		}
		for _, row := range rows {
			if row.Tags == `` {
				row.Tags = `[]`
			} else {
				tags := strings.Split(row.Tags, `,`)
				row.Tags, err = com.JSONEncodeToString(tags)
				if err != nil {
					return err
				}
			}
			err = m.UpdateField(nil, `tags`, row.Tags, db.Cond{`id`: row.Id})
			if err != nil {
				return err
			}
		}
	}
	return err
}
