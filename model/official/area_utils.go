package official

import (
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func QueryAreaNames(ctx echo.Context, areaID uint, areaM *Area) (areaNames []string, err error) {
	if areaID <= 0 {
		return
	}
	if areaM == nil {
		areaM = NewArea(ctx)
	}
	var parents []*dbschema.OfficialCommonArea
	parents, err = areaM.Positions(areaID)
	if err != nil {
		return
	}
	areaNames = make([]string, 0, len(parents)+1)
	if len(parents) > 0 {
		areaNames = append(areaNames, QueryCountryName(ctx, parents[0].CountryAbbr))
	}
	for _, parent := range parents {
		areaNames = append(areaNames, parent.Name)
	}
	return
}

func QueryCountryName(ctx echo.Context, countryAbbr string) string {
	countryM := dbschema.NewOfficialCommonAreaCountry(ctx)
	countryM.Get(func(r db.Result) db.Result {
		return r.Select(`name`)
	}, `abbr`, countryAbbr)
	if len(countryM.Name) > 0 {
		return countryM.Name
	} else if countryAbbr == `CN` {
		return `中国`
	} else {
		return countryAbbr
	}
}
