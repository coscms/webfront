package xcommon

import (
	"fmt"
	"html/template"
	"strconv"

	"github.com/admpub/decimal"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// CurrencySymbols 货币符号
var CurrencySymbols = map[string]string{
	`CNY`: `¥`,   //人民币
	`USD`: `$`,   //美元
	`JPY`: `¥`,   //日元
	`EUR`: `€`,   //欧元
	`GBP`: `£`,   //英镑
	`FRF`: `₣`,   //法郎
	`KRW`: `₩`,   //韩元
	`RUB`: `₽`,   //俄罗斯卢布
	`HKD`: `HK$`, //港元
	`AUD`: `A$`,  //澳元
	`CAD`: `C$`,  //加元
	`INR`: `₹`,   //印度卢比
}

// SetCurrencySymbol 登记货币符号
func SetCurrencySymbol(currency string, symbol string) {
	CurrencySymbols[currency] = symbol
}

// DefaultCurrency 默认币种
func DefaultCurrency() string {
	defaultCurrency := config.Setting(`base`).String(`defaultCurrency`)
	if len(defaultCurrency) == 0 {
		defaultCurrency = `CNY`
	}
	return defaultCurrency
}

// DefaultCurrencySymbol 默认币种符号
func DefaultCurrencySymbol() string {
	currencySymbol, ok := CurrencySymbols[DefaultCurrency()]
	if !ok || len(currencySymbol) == 0 {
		currencySymbol = `¥`
	}
	return currencySymbol
}

const Precision = 4 // 小数位数

func GetCurrencySymbol(ctx echo.Context, inputCurrency ...string) string {
	var currency string
	if len(inputCurrency) > 0 && len(inputCurrency[0]) > 0 {
		currency = inputCurrency[0]
	} else {
		currency = ctx.Internal().String(`currency`)
	}
	var currencySymbol string
	if len(currency) > 0 {
		if symbol, ok := CurrencySymbols[currency]; ok {
			currencySymbol = symbol
		} else {
			currencySymbol = currency
		}
	} else {
		currencySymbol = DefaultCurrencySymbol()
	}
	return currencySymbol
}

func GetCurrencyPrecision(ctx echo.Context) int32 {
	v := ctx.Internal().Get(`currencyPrecision`)
	if v == nil {
		return Precision
	}
	switch n := v.(type) {
	case string:
		i, err := strconv.ParseInt(n, 10, 32)
		if err != nil {
			return Precision
		}
		precision := int32(i)
		if precision < 0 {
			precision = Precision
		}
		return precision
	default:
		precision := param.AsInt32(n)
		if precision < 0 {
			precision = Precision
		}
		return precision
	}
}

func SetCurrencyPrecision(ctx echo.Context, precision int32) {
	ctx.Internal().Set(`currencyPrecision`, precision)
}

// HTMLCurrency HTML模板函数：币种
// amount: 金额
// withFlags[0]: 是否带货币符号
// withFlags[1]: 是否清楚小数末尾的0
func HTMLCurrency(ctx echo.Context, amount float64, withFlags ...bool) template.HTML {
	return HTMLCurrencyWithPrecision(ctx, amount, GetCurrencyPrecision(ctx), withFlags...)
}

// CurrencyWithPrecision HTML模板函数：币种
// amount: 金额
// precision: 小数位数
// withFlags[0]: 是否带货币符号
// withFlags[1]: 是否清楚小数末尾的0
func CurrencyWithPrecision(ctx echo.Context, amount float64, precision int32, withFlags ...bool) string {
	return CurrencyWithCurrencyAndPrecision(ctx, amount, ``, precision, withFlags...)
}

// CurrencyWithPrecision HTML模板函数：币种
// amount: 金额
// currency: 币种
// precision: 小数位数
// withFlags[0]: 是否带货币符号
// withFlags[1]: 是否清楚小数末尾的0
func CurrencyWithCurrencyAndPrecision(ctx echo.Context, amount float64, currency string, precision int32, withFlags ...bool) string {
	if len(withFlags) == 0 {
		return fmt.Sprintf(`%.*f`, precision, amount)
	}
	var numberFormatted string
	if len(withFlags) > 1 && withFlags[1] {
		numberFormatted = com.NumberFormat(amount, int(precision))
		numberFormatted = com.NumberTrimZero(numberFormatted)
	} else {
		numberFormatted = fmt.Sprintf(`%.*f`, precision, amount)
	}
	if withFlags[0] {
		currencySymbol := GetCurrencySymbol(ctx, currency)
		return currencySymbol + numberFormatted
	}
	return numberFormatted
}

// HTMLCurrencyWithPrecision HTML模板函数：币种
// amount: 金额
// precision: 小数位数
// withFlags[0]: 是否带货币符号
// withFlags[1]: 是否清楚小数末尾的0
func HTMLCurrencyWithPrecision(ctx echo.Context, amount float64, precision int32, withFlags ...bool) template.HTML {
	return template.HTML(CurrencyWithPrecision(ctx, amount, precision, withFlags...))
}

// HTMLCurrencySymbol HTML模板函数：币种符号
func HTMLCurrencySymbol(ctx echo.Context) template.HTML {
	currencySymbol := GetCurrencySymbol(ctx)
	return template.HTML(currencySymbol)
}

func CalcPrice(price float64, exchangeRate float64, precision ...int32) float64 {
	var _precision int32
	if len(precision) > 0 {
		_precision = precision[0]
	} else {
		_precision = Precision
	}
	priceD := decimal.NewFromFloat(price)
	exchangeRateD := decimal.NewFromFloat(exchangeRate)
	return priceD.Mul(exchangeRateD).Round(_precision).InexactFloat64()
}

func HTMLPriceFormat(ctx echo.Context, price float64, currency string, exchangeRate float64, precision ...int32) template.HTML {
	var _precision int32
	if len(precision) == 0 {
		_precision = GetCurrencyPrecision(ctx)
	} else {
		_precision = precision[0]
	}
	price = CalcPrice(price, exchangeRate, precision...)
	return template.HTML(CurrencyWithCurrencyAndPrecision(ctx, price, currency, _precision, true, true))
}
