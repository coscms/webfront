package upgrade

import (
	"strings"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/version"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func init() {
	echo.OnCallback(`nging.upgrade.db.before`, func(data echo.Event) error {
		installedPkgSchemaVer := config.GetInstalledPkgSchemaVer(version.PkgName)
		if installedPkgSchemaVer >= 2.0 {
			return nil
		}
		return upgradeArticleTagsData()
	})
}

func upgradeArticleTagsData() error {
	ctx := defaults.NewMockContext()
	m := dbschema.NewOfficialCommonArticle(ctx)
	qmw := func(r db.Result) db.Result {
		return r.Select(`id`, `tags`)
	}
	cond := db.Raw("NOT JSON_VALID(`tags`)")
	var err error
	for {
		_, err = m.ListByOffset(nil, qmw, 0, 200, cond)
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
