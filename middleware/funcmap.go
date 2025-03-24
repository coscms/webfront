/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package middleware

import (
	"html/template"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/mwutils"
	"github.com/coscms/webfront/library/xcommon"
	"github.com/coscms/webfront/model/official"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

var TmplFuncGenerator = mwutils.TmplFuncGenerators{
	`Currency`: func(ctx echo.Context) interface{} {
		return func(v float64, withFlags ...bool) interface{} {
			return xcommon.HTMLCurrency(ctx, v, withFlags...)
		}
	},
	`CurrencySymbol`: func(ctx echo.Context) interface{} {
		return func() template.HTML {
			return xcommon.HTMLCurrencySymbol(ctx)
		}
	},
}

func SetFunc(ctx echo.Context) error {
	TmplFuncGenerator.Apply(ctx)
	return nil
}

func CustomerDetail(c echo.Context) *modelCustomer.CustomerAndGroup {
	customerDetail, _ := c.Internal().Get(`customerDetail`).(*modelCustomer.CustomerAndGroup)
	return customerDetail
}

func NavigateList(ctx echo.Context, m *dbschema.OfficialCommonNavigate, navType string, parentIDs ...uint) []*official.NavigateExt {
	internalKey := `navigate.` + navType
	childrenMapping, ok := ctx.Internal().Get(internalKey).(map[uint][]*official.NavigateExt)
	if !ok {
		nav := []*official.NavigateExt{}
		m.ListByOffset(&nav, func(r db.Result) db.Result {
			return r.OrderBy(`level`, `sort`, `id`)
		}, 0, -1, db.And(
			db.Cond{`disabled`: `N`},
			db.Cond{`type`: navType},
		))
		childrenMapping = map[uint][]*official.NavigateExt{}
		for _, _nav := range nav {
			_nav.Init().SetContext(ctx)
			if _, ok := childrenMapping[_nav.ParentId]; !ok {
				childrenMapping[_nav.ParentId] = []*official.NavigateExt{}
			}
			childrenMapping[_nav.ParentId] = append(childrenMapping[_nav.ParentId], _nav)
		}
		for _, _nav := range nav {
			children, ok := childrenMapping[_nav.Id]
			if !ok {
				continue
			}
			_nav.Children = &children
		}
		ctx.Internal().Set(internalKey, childrenMapping)
	}
	if len(parentIDs) > 0 {
		return childrenMapping[parentIDs[0]]
	}
	return childrenMapping[0]
}

func FuncMap() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			SetFunc(c)
			return h.Handle(c)
		})
	}
}
