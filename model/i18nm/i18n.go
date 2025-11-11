package i18nm

import (
	"strings"

	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

func GetTranslations(ctx echo.Context, table string, ids []uint64) map[uint64]map[string]string {
	m := map[uint64]map[string]string{}
	rM := dbschema.NewOfficialI18nResource(ctx)
	rM.ListByOffset(nil, nil, 0, -1, `code`, db.Like(table+`.%`))
	rows := rM.Objects()
	rIDs := make([]uint, len(rows))
	rKeys := map[uint]string{}
	for i, v := range rows {
		rIDs[i] = v.Id
		rKeys[v.Id] = strings.SplitN(v.Code, `.`, 2)[1]
	}
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	tM.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`lang`: ctx.Lang().Normalize()},
		db.Cond{`row_id`: db.In(ids)},
		db.Cond{`resource_id`: db.In(rIDs)},
	))
	tRows := tM.Objects()
	for _, v := range tRows {
		if _, ok := m[v.RowId]; !ok {
			m[v.RowId] = map[string]string{}
		}
		m[v.RowId][rKeys[v.ResourceId]] = v.Text
	}
	return m
}

func GetModelTranslations(mdl factory.Model, ids []uint64) map[uint64]map[string]string {
	return GetTranslations(mdl.Context(), mdl.Short_(), ids)
}

func Initialize(ctx echo.Context) error {
	rM := dbschema.NewOfficialI18nResource(ctx)
	var exists bool
	var err error

	for table, fieldInfo := range dbschema.DBI.Fields {
		for field, info := range fieldInfo {
			if !info.Multilingual {
				continue
			}
			exists, err = rM.Exists(nil, `code`, table+`.`+field)
			if err != nil {
				return err
			}
			if exists {
				continue
			}
			rM.Code = table + `.` + field
			_, err = rM.Insert()
			if err != nil {
				return err
			}
			rM.Reset()
		}
	}
	return err
}
