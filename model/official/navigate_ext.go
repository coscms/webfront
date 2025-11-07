package official

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/null"
	"github.com/coscms/webfront/dbschema"
	"github.com/phuslu/lru"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

var NavigateLinkType = echo.NewKVData()
var navigateIdentRegexps = lru.NewLRUCache[string, *regexp.Regexp](100)

func init() {
	NavigateLinkType.Add(`custom`, echo.T(`自定义链接`))
	NavigateLinkType.Add(`article-category`, echo.T(`文章分类`), echo.KVOptFn(queryArticleCategory))
}

func queryArticleCategory(c context.Context) interface{} {
	ctx := c.(echo.Context)
	m := NewCategory(ctx)
	categories := m.ListAllParentBy(`article`, 0, 2, db.Cond{`show_on_menu`: `Y`})
	var (
		list     []*NavigateExt
		children = map[uint][]*NavigateExt{}
		ext      = ctx.DefaultExtension()
	)
	for _, category := range categories {
		navExt := &NavigateExt{
			OfficialCommonNavigate: &dbschema.OfficialCommonNavigate{
				Id:       category.Id,
				ParentId: category.ParentId,
				Title:    category.Name,
				Url:      ctx.URLFor(`/articles` + ext + `?categoryId=` + fmt.Sprint(category.Id)),
			},
			Extra: echo.H{
				`object`: category,
			},
			isActive: GenActiveDetector(`categoryId`, category.Id),
		}
		navExt.SetInsideURL(navExt.Url)
		navExt.Init().SetContext(ctx)
		if category.ParentId < 1 {
			list = append(list, navExt)
			continue
		}
		if _, ok := children[category.ParentId]; !ok {
			children[category.ParentId] = []*NavigateExt{}
		}
		children[category.ParentId] = append(children[category.ParentId], navExt)
	}
	fillArticleChildrenCategory(&list, children)
	return list
}

func fillArticleChildrenCategory(list *[]*NavigateExt, children map[uint][]*NavigateExt) {
	if len(children) == 0 {
		return
	}
	for _, item := range *list {
		subItems, ok := children[item.Id]
		if !ok {
			continue
		}
		item.Children = &subItems
		delete(children, item.Id)
		fillArticleChildrenCategory(item.Children, children)
	}
}

func NewNavigateExt(nav *dbschema.OfficialCommonNavigate) *NavigateExt {
	return &NavigateExt{OfficialCommonNavigate: nav}
}

func GenActiveDetector(categoryKey string, categoryID uint) func(ctx echo.Context) bool {
	return func(ctx echo.Context) bool {
		return ctx.Formx(categoryKey).Uint() == categoryID
	}
}

type NavigateExt struct {
	*dbschema.OfficialCommonNavigate
	identRegexp *regexp.Regexp
	isInside    null.Bool // 是否是内部链接
	isActive    func(echo.Context) bool
	insideURL   string          // 内部链接的URL
	Children    *[]*NavigateExt //`db:",relation=parent_id:id|eq(has_child,Y)"`
	Extra       echo.H
}

func loader(ctx context.Context, s string) (*regexp.Regexp, error) {
	return regexp.Compile(s)
}

func (f *NavigateExt) initIdentRegexp() error {
	if f.identRegexp != nil {
		return nil
	}
	var err error
	if expr, found := strings.CutPrefix(f.Ident, `regexp:`); found && len(expr) > 0 {
		f.identRegexp, err, _ = navigateIdentRegexps.GetOrLoad(context.Background(), expr, loader)
		if err != nil {
			log.Error(expr+`: `, err)
		}
	}
	return err
}

func (f *NavigateExt) Init() *NavigateExt {
	if len(f.Ident) > 0 {
		f.initIdentRegexp()
	}
	return f
}

func (f *NavigateExt) IsInside() bool {
	if !f.isInside.Valid {
		f.isInside.Bool = strings.HasPrefix(f.Url, `/`)
		f.isInside.Valid = true
	}
	return f.isInside.Bool
}

func (f *NavigateExt) URL() string {
	if len(f.Url) == 0 {
		return ``
	}
	if len(f.insideURL) > 0 {
		return f.insideURL
	}
	if f.IsInside() {
		f.insideURL = f.Context().URLFor(f.Url)
		return f.insideURL
	}
	return f.Url
}

func (f *NavigateExt) SetInsideURL(insideURL string) *NavigateExt {
	f.insideURL = insideURL
	return f
}

func (f *NavigateExt) IsValidURL() bool {
	if len(f.Url) == 0 || strings.HasPrefix(f.Url, `#`) {
		return false
	}
	return true
}

func (f *NavigateExt) IsActive() bool {
	if f.isActive != nil {
		return f.isActive(f.Context())
	}
	if f.IsInside() {
		navURL := f.Url
		if navURL == f.Context().Request().URL().Path() ||
			navURL == f.Context().DispatchPath() {
			return true
		}
	}
	if len(f.Ident) > 0 {
		err := f.initIdentRegexp()
		if err != nil {
			return false
		}
		if f.identRegexp != nil {
			return f.identRegexp.MatchString(f.Context().DispatchPath())
		}
		return strings.HasSuffix(f.Context().DispatchPath(), f.Ident)
	}
	return false
}

func (f *NavigateExt) SetActiveDetector(fn func(echo.Context) bool) *NavigateExt {
	f.isActive = fn
	return f
}

func (f *NavigateExt) SetExtra(extra echo.H) *NavigateExt {
	f.Extra = extra
	return f
}

func (f *NavigateExt) SetExtraKV(k string, v interface{}) *NavigateExt {
	if f.Extra == nil {
		f.Extra = echo.H{}
	}
	f.Extra.Set(k, v)
	return f
}

func (f *NavigateExt) HasChildren() bool {
	f.getChildren()
	return len(*f.Children) > 0
}

func (f *NavigateExt) FetchChildren(forces ...bool) []*NavigateExt {
	f.getChildren(forces...)
	return *f.Children
}

func (f *NavigateExt) ClearChildren() *NavigateExt {
	f.Children = nil
	return f
}

func (f *NavigateExt) getChildren(forces ...bool) *NavigateExt {
	var force bool
	if len(forces) > 0 {
		force = forces[0]
	}
	if !force && f.Children != nil {
		return f
	}
	defer func() {
		if f.Children == nil {
			f.Children = &[]*NavigateExt{}
		}
	}()
	if len(f.LinkType) == 0 {
		return f
	}
	item := NavigateLinkType.GetItem(f.LinkType)
	if item == nil {
		return f
	}
	if item.Fn() == nil {
		return f
	}
	if list, ok := item.Fn()(f.Context()).([]*NavigateExt); ok {
		f.Children = &list
	}
	return f
}
