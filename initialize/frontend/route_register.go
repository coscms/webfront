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

package frontend

import (
	"github.com/coscms/webfront/registry/route"
)

var (
	IRegister          = route.IRegister
	MakeHandler        = route.MakeHandler
	WithMeta           = route.MetaHandler
	WithMetaAndRequest = route.MetaHandlerWithRequest
	WithRequest        = route.HandlerWithRequest
	Pre                = route.Pre
	Use                = route.Use
	UseToGroup         = route.UseToGroup
	Register           = route.Register
	RegisterToGroup    = route.RegisterToGroup
	Host               = route.Host
	Apply              = route.Apply
	Prefix             = route.Prefix
)
