package top

import (
	"github.com/webx-top/com"
)

var (
	ErrInvalidDuration = com.ErrInvalidDuration
)

const (
	DurationDay     = com.DurationDay
	DurationMonth   = com.DurationMonth
	DurationWeek    = com.DurationWeek
	DurationYear    = com.DurationYear
	MonthMaxSeconds = com.MonthMaxSeconds // 月份中可能的最大秒数
)

// ParseDuration 解析持续时间(在支持标准库time.ParseDuration的基础上增加了年(y)月(mo)周(w)日(d)的支持)
var ParseDuration = com.ParseTimeDuration

// IsSameDay 是否为同一天
var IsSameDay = com.IsSameDay

// MonthDay 计算某个月的天数
var MonthDay = com.MonthDay

// MonthDayByTime 计算某个月的天数
var MonthDayByTime = com.MonthDayByTime

// IsSameMonth 是否为同一月
var IsSameMonth = com.IsSameMonth

var TodayTimestamp = com.TodayTimestamp

var DayStart = com.DayStart

var DayEnd = com.DayEnd
