package i18nm

import (
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"
)

// Search performs a keyword search on i18n translations for the specified table and resource IDs.
// It joins with the translations table to find matches in the specified language.
// Parameters:
//   - ctx: echo context containing request information
//   - table: name of the resource table to search
//   - keyword: search term to look for in translations
//   - param: factory parameters for building the query
//   - columns: optional columns to select from the resource table
//
// Returns:
//   - error if any occurs during the search operation
func Search(ctx echo.Context, table string, keyword string, param *factory.Param, columns ...string) error {
	rows, err := getResources(ctx, table, columns...)
	if err != nil || len(rows) == 0 {
		return err
	}
	rIDs := make([]uint, len(rows))
	for i, v := range rows {
		rIDs[i] = v.Id
	}
	cond := db.NewCompounds()
	cond.Add(db.Cond{`TR.lang`: ctx.Lang().Normalize()})
	cond.Add(db.Cond{`TR.resource_id`: db.In(rIDs)})
	cond.From(mysql.SearchField(`~TR.text`, keyword))
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	param.AddJoin(`INNER`, tM.Short_(), `TR`, table+`.id = TR.row_id`).AddArgs(cond.And())
	return err
}
