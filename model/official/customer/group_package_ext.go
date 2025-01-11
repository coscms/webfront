package customer

import (
	"github.com/webx-top/echo"
)

const (
	GroupPackageTimeDay     = `day`
	GroupPackageTimeWeek    = `week`
	GroupPackageTimeMonth   = `month`
	GroupPackageTimeYear    = `year`
	GroupPackageTimeForever = `forever`
)

var GroupPackageTimeUnits = echo.NewKVData().
	Add(GroupPackageTimeDay, echo.T(`天`)).   // echo.T(`%d天`)
	Add(GroupPackageTimeWeek, echo.T(`周`)).  // echo.T(`%d周`)
	Add(GroupPackageTimeMonth, echo.T(`月`)). // echo.T(`%d月`) echo.T(`%d个月`)
	Add(GroupPackageTimeYear, echo.T(`年`)).  // echo.T(`%d年`)
	Add(GroupPackageTimeForever, echo.T(`永久`))

func i18nUnit(c echo.Context, unit string) string {
	return c.T(GroupPackageTimeUnits.Get(unit))
}

func GroupPackageTimeUnitSuffix(c echo.Context, n uint, unit string) string {
	switch unit {
	case GroupPackageTimeMonth:
		if n > 1 {
			return `/ ` + c.T(`%d个月`, n)
		}
		return `/ ` + i18nUnit(c, unit)
	case GroupPackageTimeForever:
		return i18nUnit(c, unit)
	default:
		if n > 1 {
			return `/ ` + c.T(`%d`+unit, n)
		}
		return `/ ` + i18nUnit(c, unit)
	}
}
