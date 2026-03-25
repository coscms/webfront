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
		countryM := dbschema.NewOfficialCommonAreaCountry(ctx)
		countryM.Get(func(r db.Result) db.Result {
			return r.Select(`name`)
		}, `abbr`, parents[0].CountryAbbr)
		if len(countryM.Name) > 0 {
			areaNames = append(areaNames, countryM.Name)
		} else if parents[0].CountryAbbr == `CN` {
			areaNames = append(areaNames, `中国`)
		} else {
			areaNames = append(areaNames, parents[0].CountryAbbr)
		}
	}
	for _, parent := range parents {
		areaNames = append(areaNames, parent.Name)
	}
	return
}
