package middleware

import (
	"html/template"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/license"
	"github.com/coscms/webcore/library/nsql"
	"github.com/coscms/webcore/library/ntemplate"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/logic/articlelogic"
	"github.com/coscms/webfront/library/xcommon"
	"github.com/coscms/webfront/model/official"
	modelAdvert "github.com/coscms/webfront/model/official/advert"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"
)

func init() {
	tplfunc.TplFuncMap[`Advert`] = func(idents ...string) interface{} {
		ctx := defaults.AcquireMockContext()
		r := modelAdvert.GetAdvertForHTML(ctx, idents...)
		defaults.ReleaseMockContext(ctx)
		return r
	}
}

var DefaultRenderDataWrapper = func(ctx echo.Context, data interface{}) interface{} {
	return NewRenderData(ctx, data)
}

func NewRenderData(ctx echo.Context, data interface{}) *RenderData {
	if v, ok := data.(*RenderData); ok {
		return v
	}
	return &RenderData{
		ctx:        ctx,
		RenderData: echo.NewRenderData(ctx, data),
	}
}

type RenderData struct {
	ctx echo.Context
	*echo.RenderData
}

func (r *RenderData) Customer() *dbschema.OfficialCustomer {
	return Customer(r.ctx)
}

func (r *RenderData) CustomerDetail() *modelCustomer.CustomerAndGroup {
	return CustomerDetail(r.ctx)
}

func (r *RenderData) Friendlink(limit int, categoryIds ...uint) []*dbschema.OfficialCommonFriendlink {
	m := official.NewFriendlink(r.ctx)
	list, _ := m.ListShowAndRecord(limit, categoryIds...)
	return list
}

func (r *RenderData) FrontendNav(parentIDs ...uint) []*official.NavigateExt {
	return NavigateList(r.ctx, dbschema.NewOfficialCommonNavigate(r.ctx), `default`, parentIDs...)
}

func (r *RenderData) CustomerNav(parentIDs ...uint) []*official.NavigateExt {
	return NavigateList(r.ctx, dbschema.NewOfficialCommonNavigate(r.ctx), `userCenter`, parentIDs...)
}

func (r *RenderData) SQLQuery() *nsql.SQLQuery {
	return nsql.NewSQLQuery(r.ctx)
}

func (r *RenderData) SQLQueryLimit(offset int, limit int, linkID ...int) *nsql.SQLQuery {
	return nsql.NewSQLQueryLimit(r.ctx, offset, limit, linkID...)
}

func (r *RenderData) CaptchaForm(tmpl string, args ...interface{}) template.HTML {
	return captchabiz.CaptchaForm(r.ctx, tmpl, args...)
}

func (r *RenderData) CaptchaFormWithType(typ string, tmpl string, args ...interface{}) template.HTML {
	return captchabiz.CaptchaFormWithType(r.ctx, typ, tmpl, args...)
}

func (r *RenderData) TagList(group ...string) []*dbschema.OfficialCommonTags {
	list, _ := articlelogic.GetTags(r.ctx, group...)
	return list
}

func (r *RenderData) CategoryList(limit int, ctype ...string) []*dbschema.OfficialCommonCategory {
	categories, _ := articlelogic.GetCategories(r.ctx, limit, ctype...)
	return categories
}

func (r *RenderData) SubCategoryList(parentId int, limit int, ctype ...string) []*dbschema.OfficialCommonCategory {
	categories, _ := articlelogic.GetSubCategories(r.ctx, parentId, limit, ctype...)
	return categories
}

func (r *RenderData) SoftwareURL() string {
	if license.SkipLicenseCheck {
		return ``
	}
	return license.ProductURL()
}

func (r *RenderData) SkipLicenseCheck() bool {
	return license.SkipLicenseCheck
}

func (r *RenderData) SoftwareName() string {
	return bootconfig.SoftwareName
}

// Advert 广告
//
//	{{$ads := $.Advert `test`}}
//	{{$ads.Rand.HTML}} {{$ads.First.HTML}} {{$ads.Last.HTML}} {{range $adk,$adv := $ads}}...{{end}}
//
//	{{$ads := $.Advert `test` `test2`}}
//	{{($ads.Place `test`).Rand.HTML}} {{($ads.Place `test`).First.HTML}} {{($ads.Place `test`).Last.HTML}}  {{range $adk,$adv := $ads.Place `test`}}...{{end}}
func (r *RenderData) Advert(idents ...string) interface{} {
	return modelAdvert.GetAdvertForHTML(r.ctx, idents...)
}

func (r *RenderData) ThemeInfo() *ntemplate.ThemeInfo {
	return httpserver.Frontend.Template.ThemeInfo(r.ctx)
}

func (r *RenderData) Price(price float64) float64 {
	conv, ok := r.ctx.Internal().Get(`CurrencyRate`).(FloatConverter)
	if !ok {
		return price
	}
	return conv.Convert(price)
}

func (r *RenderData) PriceFormat(price float64) template.HTML {
	conv, ok := r.ctx.Internal().Get(`CurrencyRate`).(echo.RenderContextWithData)
	if !ok {
		h := xcommon.HTMLCurrency(r.ctx, price, true)
		switch r := h.(type) {
		case template.HTML:
			return r
		case string:
			return template.HTML(r)
		default:
			return template.HTML(param.AsString(r))
		}
	}
	return conv.RenderWithData(r.ctx, price)
}

func (r *RenderData) Currency() string {
	conv, ok := r.ctx.Internal().Get(`CurrencyRate`).(CurrencyGetter)
	if !ok {
		return xcommon.DefaultCurrency()
	}
	return conv.Currency()
}

type FloatConverter interface {
	Convert(float64) float64
}

type CurrencyGetter interface {
	Currency() string
}
