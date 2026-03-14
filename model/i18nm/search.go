package i18nm

import (
	"strings"

	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/echo"
)

func buildDefaultLangSearch(fields []string, alias string, keyword string) db.Compound {
	for i, v := range fields {
		if len(alias) > 0 {
			v = alias + `.` + v
		}
		fields[i] = v
	}
	return mysql.SearchField(`~`+strings.Join(fields, `+`), keyword)
}

// Search performs a keyword search on i18n translations for the specified table and resource IDs.
// Parameters:
//   - ctx: echo context containing request information
//   - table: name of the resource table to search
//   - alias: alias of the resource table to search
//   - keyword: search term to look for in translations
//   - param: factory parameters for building the query
//   - columns: optional columns to select from the resource table
//
// Returns:
//   - error if any occurs during the search operation
func Search(ctx echo.Context, table string, alias string, keyword string, param *factory.Param, columns ...string) error {
	if IsDefaultLang(ctx) {
		if len(columns) == 0 {
			return nil
		}
		param.AddArgs(buildDefaultLangSearch(columns, alias, keyword))
		return nil
	}
	return SearchLanguage(ctx, ctx.Lang().Normalize(), table, alias, keyword, param, columns...)
}

// SearchLanguage
func SearchLanguage(ctx echo.Context, langCode string, table string, alias string, keyword string, param *factory.Param, columns ...string) error {
	fields := make([]string, len(columns))
	copy(fields, columns)
	rows, err := getResources(ctx, table, columns...) // columns 会被更改
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		if len(columns) == 0 {
			return err
		}
		for i, v := range fields {
			if len(alias) > 0 {
				v = alias + `.` + v
			}
			fields[i] = v
		}
		param.AddArgs(buildDefaultLangSearch(fields, alias, keyword))
		return err
	}
	rIDs := make([]uint, len(rows))
	for i, v := range rows {
		rIDs[i] = v.Id
	}
	if len(alias) == 0 {
		alias = dbschema.WithPrefix(table)
	}
	tM := dbschema.NewOfficialI18nTranslation(ctx)
	var trCompound db.Compound
	// cond := db.NewCompounds()
	// cond.Add(db.Cond{`TR.lang`: langCode})
	// cond.Add(db.Cond{`TR.resource_id`: db.In(rIDs)})
	// cond.From(mysql.SearchField(`~TR.text`, keyword))
	// param.AddJoin(`INNER`, tM.Short_(), `TR`, alias+`.id = TR.row_id`)
	// trCompound = cond.And()
	trCompound = db.Raw(
		`EXISTS (SELECT 1 FROM `+tM.Name_()+` AS TR WHERE TR.row_id=`+alias+`.id AND TR.lang=? AND TR.resource_id IN ? AND `+mysql.Match(keyword, true, `TR.text`).String()+`)`,
		langCode,
		rIDs,
	)
	if len(columns) > 0 {
		for i, v := range fields {
			fields[i] = alias + `.` + v
		}
		param.AddArgs(db.Or(
			trCompound,
			mysql.SearchField(`~`+strings.Join(fields, `+`), keyword),
		))
		return err
	}
	param.AddArgs(trCompound)
	return err
}

// SearchModel searches for records in the specified model using the given keyword and parameters.
// It delegates the actual search operation to the Search function with the model's short name.
// ctx: Echo context for the request
// mdl: The model to search in
// alias: alias of the model table
// keyword: The search term
// param: Additional search parameters
// columns: Optional columns to search in
// Returns an error if the search fails
func SearchModel(ctx echo.Context, mdl Model, alias string, keyword string, param *factory.Param, columns ...string) error {
	return Search(ctx, mdl.Short_(), alias, keyword, param, columns...)
}
