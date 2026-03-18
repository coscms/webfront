package xcommon

import (
	"html/template"

	"github.com/admpub/decimal"
	"github.com/webx-top/echo"
)

type PriceFormatter interface {
	PriceFormat(ctx echo.Context, price float64, precision ...int32) string
}

type FloatConverter interface {
	Convert(price float64, precision ...int32) float64
}

type CurrencyGetter interface {
	Currency() string
}

func CurrencyByInternal(ctx echo.Context) string {
	conv, ok := ctx.Internal().Get(`CurrencyRate`).(CurrencyGetter)
	if !ok {
		return DefaultCurrency()
	}
	return conv.Currency()
}

func CalcPriceByInternal(ctx echo.Context, price float64, precision ...int32) float64 {
	conv, ok := ctx.Internal().Get(`CurrencyRate`).(FloatConverter)
	if !ok {
		return price
	}
	return conv.Convert(price, precision...)
}

func PriceFormatByInternal(ctx echo.Context, price float64, precision ...int32) template.HTML {
	conv, ok := ctx.Internal().Get(`CurrencyRate`).(PriceFormatter)
	if !ok {
		if len(precision) == 0 {
			return HTMLCurrency(ctx, price, true, true)
		}
		return HTMLCurrencyWithPrecision(ctx, price, precision[0], true, true)
	}
	return template.HTML(conv.PriceFormat(ctx, price, precision...))
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
