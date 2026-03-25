package official

import (
	"github.com/coscms/webfront/dbschema"
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
	areaNames = make([]string, len(parents))
	for index, parent := range parents {
		areaNames[index] = parent.Name
	}
	return
}
