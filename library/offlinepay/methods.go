package offlinepay

import "github.com/webx-top/echo"

/*
//
中国银行 - BOC
农业银行 - ABC
工商银行 - ICBC
建设银行 - CCB
交通银行 - BOCOM
招商银行 - CMB
浦发银行 - SPDB
光大银行 - CEB
华夏银行 - HXB
民生银行 - CMBC

恒丰银行 - HFB
汉口银行 - HKB
//
*/
const (
	OfflinePayMethodEBankAlipay   = `ebank.alipay`   // 支付宝
	OfflinePayMethodEBankWechat   = `ebank.wechat`   // 微信
	OfflinePayMethodEBankPaypal   = `ebank.paypal`   // Paypal
	OfflinePayMethodEBankUnionPay = `ebank.unionpay` // 云闪付
	OfflinePayMethodEBankBestPay  = `ebank.bestpay`  // 翼支付
	OfflinePayMethodBankICBC      = `bank.icbc`      // 工商银行
	OfflinePayMethodBankABC       = `bank.abc`       // 农业银行
	OfflinePayMethodBankBOC       = `bank.boc`       // 中国银行
	OfflinePayMethodBankCCB       = `bank.ccb`       // 建设银行
	OfflinePayMethodBankBOCOM     = `bank.bocom`     // 交通银行
	OfflinePayMethodBankCMB       = `bank.cmb`       // 招商银行
	OfflinePayMethodBankSPDB      = `bank.spdb`      // 浦发银行
	OfflinePayMethodBankCEB       = `bank.ceb`       // 光大银行
	OfflinePayMethodBankHXB       = `bank.hxb`       // 华夏银行
	OfflinePayMethodBankCMBC      = `bank.cmbc`      // 民生银行
	OfflinePayMethodBankHFB       = `bank.hfb`       // 恒丰银行
	OfflinePayMethodBankHKB       = `bank.hkb`       // 汉口银行
)

type MethodX struct {
	Disabled bool
	Logo     string
}

func (a *MethodX) CopyFrom(s MethodOptions) {
	if s.Disabled != nil {
		a.Disabled = *s.Disabled
	}
	if s.Logo != nil {
		a.Logo = *s.Logo
	}
}

var methodsRegistry = echo.NewKVxData[MethodX, any]().
	Add(OfflinePayMethodEBankAlipay, echo.T(`支付宝`)).
	Add(OfflinePayMethodEBankWechat, echo.T(`微信`)).
	Add(OfflinePayMethodEBankUnionPay, echo.T(`云闪付`)).
	Add(OfflinePayMethodEBankBestPay, echo.T(`翼支付`)).
	Add(OfflinePayMethodEBankPaypal, echo.T(`Paypal`)).
	Add(OfflinePayMethodBankICBC, echo.T(`工商银行`)).
	Add(OfflinePayMethodBankABC, echo.T(`农业银行`)).
	Add(OfflinePayMethodBankBOC, echo.T(`中国银行`)).
	Add(OfflinePayMethodBankCCB, echo.T(`建设银行`)).
	Add(OfflinePayMethodBankBOCOM, echo.T(`交通银行`)).
	Add(OfflinePayMethodBankCMB, echo.T(`招商银行`)).
	Add(OfflinePayMethodBankSPDB, echo.T(`浦发银行`)).
	Add(OfflinePayMethodBankCEB, echo.T(`光大银行`)).
	Add(OfflinePayMethodBankHXB, echo.T(`华夏银行`)).
	Add(OfflinePayMethodBankCMBC, echo.T(`民生银行`)).
	Add(OfflinePayMethodBankHFB, echo.T(`恒丰银行`)).
	Add(OfflinePayMethodBankHKB, echo.T(`汉口银行`))

func RegisterMethod(k, v string, cfg MethodX) {
	methodsRegistry.Add(k, v, echo.KVxOptX[MethodX, any](cfg))
}

func GetMethods(opts map[string]MethodOptions) []*echo.KVx[MethodX, any] {
	sl := methodsRegistry.Slice()
	rs := make([]*echo.KVx[MethodX, any], 0, len(sl))
	if len(opts) == 0 {
		opts = GetConfig().Methods
	}
	for _, item := range sl {
		cloned := item.Clone()
		if opts != nil {
			if opt, ok := opts[item.K]; ok {
				cloned.X.CopyFrom(opt)
			}
		}
		if cloned.X.Disabled {
			continue
		}
		rs = append(rs, &cloned)
	}
	return rs
}

func GetMethod(name string, opts map[string]MethodOptions) *echo.KVx[MethodX, any] {
	item := methodsRegistry.GetItem(name)
	if item == nil {
		return nil
	}
	cloned := item.Clone()
	if len(opts) == 0 {
		opts = GetConfig().Methods
	}
	if opts != nil {
		if opt, ok := opts[item.K]; ok {
			cloned.X.CopyFrom(opt)
		}
	}
	if cloned.X.Disabled {
		return nil
	}
	return &cloned
}
