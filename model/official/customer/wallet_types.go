package customer

import (
	"html/template"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var (
	// AssetTypes 资产类型
	AssetTypes = echo.NewKVData()
	// AmountTypes 金额类型
	AmountTypes = echo.NewKVData()
	// 资金流水记录状态
	FlowStatus = echo.NewKVData()
)

// 资产类型
const (
	AssetTypeMoney      = `money`      // 人民币
	AssetTypeIntegral   = `integral`   // 积分
	AssetTypeCredit     = `credit`     // 信用分
	AssetTypePoint      = `point`      // 点数
	AssetTypeExperience = `experience` // 经验
)

// 金额类型
const (
	AmountTypeBalance = `balance` // 余额
	AmountTypeFreeze  = `freeze`  // 冻结金额
)

// 资金流水记录状态
const (
	//状态(pending-待确认;confirmed-已确认;canceled-已取消)
	FlowStatusPending   = `pending`   // 待确认
	FlowStatusConfirmed = `confirmed` // 已确认
	FlowStatusCanceled  = `canceled`  // 已取消
)

func init() {
	// 注册资产类型
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypeMoney,
		V: `人民币`,
		H: echo.H{
			`icon`:   `icon ion-cash`,
			`bg`:     `warning`,
			`symbol`: `&yen;`,
		},
	})
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypeIntegral,
		V: `积分`,
		H: echo.H{
			`icon`: `icon ion-android-cart`,
			`bg`:   `indigo-light`,
		},
	})
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypeCredit,
		V: `信用分`,
		H: echo.H{
			`icon`:              `icon ion-heart`,
			`bg`:                `pink-light`,
			`comment`:           `满分10`,
			`ignoreAccumulated`: true, // 不支持累计历史值
		},
	})
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypePoint,
		V: `点数`,
		H: echo.H{
			`icon`: `iconfont icon-dengji-11`,
			`bg`:   `gray`,
		},
	})
	AssetTypes.AddItem(&echo.KV{
		K: AssetTypeExperience,
		V: `经验`,
		H: echo.H{
			`icon`: `iconfont icon-youxiu`,
			`bg`:   `info`,
		},
	})

	// 注册金额类型
	AmountTypes.Add(AmountTypeBalance, `余额`)
	AmountTypes.Add(AmountTypeFreeze, `冻结`)

	// 注册资金流水数据的状态
	FlowStatus.Add(FlowStatusPending, `待确认`)
	FlowStatus.Add(FlowStatusConfirmed, `已确认`)
	FlowStatus.Add(FlowStatusCanceled, `已取消`)
}

func AssetTypeList() []*echo.KV {
	copied := []*echo.KV{}
	for _, assetType := range AssetTypes.Slice() {
		if assetType.H.Bool(`disabled`) {
			continue
		}
		copied = append(copied, assetType)
	}
	return copied
}

// AssetTypeIsIgnoreAccumulated 判断某种资产类型是否不支持累计历史值
func AssetTypeIsIgnoreAccumulated(assetType string) bool {
	item := AssetTypes.GetItem(assetType)
	if item == nil {
		return false
	}
	return item.H.Bool(`ignoreAccumulated`)
}

func assetSymbol(item *echo.KV) string {
	var symbol string
	if item.H != nil && item.H.Has(`symbol`) {
		symbol = item.H.String(`symbol`)
	}
	return symbol
}

// MakeAssetAmountFormatter 构造指定类型资产的金额格式化函数
func MakeAssetAmountFormatter(ctx echo.Context, assetType string) func(amount float64) template.HTML {
	item := AssetTypes.GetItem(assetType)
	if item == nil {
		return func(amount float64) template.HTML {
			return template.HTML(com.String(amount))
		}
	}
	if item.X == nil {
		symbol := assetSymbol(item)
		return func(amount float64) template.HTML {
			return template.HTML(symbol + com.String(amount))
		}
	}
	rd, ok := item.X.(echo.RenderContextWithData)
	if !ok {
		symbol := assetSymbol(item)
		return func(amount float64) template.HTML {
			return template.HTML(symbol + com.String(amount))
		}
	}
	return func(amount float64) template.HTML {
		return rd.RenderWithData(ctx, amount)
	}
}

// MakeAnyAssetAmountFormatter 构造任意类型资产的金额格式化函数
func MakeAnyAssetAmountFormatter(ctx echo.Context) func(assetType string, amount float64) template.HTML {
	return func(assetType string, amount float64) template.HTML {
		item := AssetTypes.GetItem(assetType)
		if item == nil {
			return template.HTML(com.String(amount))
		}
		if item.X == nil {
			symbol := assetSymbol(item)
			return template.HTML(symbol + com.String(amount))
		}
		rd, ok := item.X.(echo.RenderContextWithData)
		if !ok {
			symbol := assetSymbol(item)
			return template.HTML(symbol + com.String(amount))
		}
		return rd.RenderWithData(ctx, amount)
	}
}
