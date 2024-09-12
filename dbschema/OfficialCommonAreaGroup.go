// @generated Do not edit this file, which is automatically generated by the generator.

package dbschema

import (
	"fmt"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type Slice_OfficialCommonAreaGroup []*OfficialCommonAreaGroup

func (s Slice_OfficialCommonAreaGroup) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCommonAreaGroup) RangeRaw(fn func(m *OfficialCommonAreaGroup) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCommonAreaGroup) GroupBy(keyField string) map[string][]*OfficialCommonAreaGroup {
	r := map[string][]*OfficialCommonAreaGroup{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialCommonAreaGroup{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialCommonAreaGroup) KeyBy(keyField string) map[string]*OfficialCommonAreaGroup {
	r := map[string]*OfficialCommonAreaGroup{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialCommonAreaGroup) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialCommonAreaGroup) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialCommonAreaGroup) FromList(data interface{}) Slice_OfficialCommonAreaGroup {
	values, ok := data.([]*OfficialCommonAreaGroup)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialCommonAreaGroup{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialCommonAreaGroup(ctx echo.Context) *OfficialCommonAreaGroup {
	m := &OfficialCommonAreaGroup{}
	m.SetContext(ctx)
	return m
}

// OfficialCommonAreaGroup 地区分组
type OfficialCommonAreaGroup struct {
	base    factory.Base
	objects []*OfficialCommonAreaGroup

	Id          uint   `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	CountryAbbr string `db:"country_abbr" bson:"country_abbr" comment:"国家缩写" json:"country_abbr" xml:"country_abbr"`
	Name        string `db:"name" bson:"name" comment:"组名称" json:"name" xml:"name"`
	Abbr        string `db:"abbr" bson:"abbr" comment:"组缩写" json:"abbr" xml:"abbr"`
	AreaIds     string `db:"area_ids" bson:"area_ids" comment:"根地区ID" json:"area_ids" xml:"area_ids"`
	Sort        int    `db:"sort" bson:"sort" comment:"排序编号" json:"sort" xml:"sort"`
}

// - base function

func (a *OfficialCommonAreaGroup) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialCommonAreaGroup) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialCommonAreaGroup) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialCommonAreaGroup) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialCommonAreaGroup) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialCommonAreaGroup) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialCommonAreaGroup) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialCommonAreaGroup) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialCommonAreaGroup) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialCommonAreaGroup) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialCommonAreaGroup) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialCommonAreaGroup) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialCommonAreaGroup) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialCommonAreaGroup) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialCommonAreaGroup) Objects() []*OfficialCommonAreaGroup {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialCommonAreaGroup) XObjects() Slice_OfficialCommonAreaGroup {
	return Slice_OfficialCommonAreaGroup(a.Objects())
}

func (a *OfficialCommonAreaGroup) NewObjects() factory.Ranger {
	return &Slice_OfficialCommonAreaGroup{}
}

func (a *OfficialCommonAreaGroup) InitObjects() *[]*OfficialCommonAreaGroup {
	a.objects = []*OfficialCommonAreaGroup{}
	return &a.objects
}

func (a *OfficialCommonAreaGroup) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialCommonAreaGroup) Short_() string {
	return "official_common_area_group"
}

func (a *OfficialCommonAreaGroup) Struct_() string {
	return "OfficialCommonAreaGroup"
}

func (a *OfficialCommonAreaGroup) Name_() string {
	b := a
	if b == nil {
		b = &OfficialCommonAreaGroup{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialCommonAreaGroup) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialCommonAreaGroup) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	base := a.base
	if !a.base.Eventable() {
		err = a.Param(mw, args...).SetRecv(a).One()
		a.base = base
		return
	}
	queryParam := a.Param(mw, args...).SetRecv(a)
	if err = DBI.FireReading(a, queryParam); err != nil {
		return
	}
	err = queryParam.One()
	a.base = base
	if err == nil {
		err = DBI.FireReaded(a, queryParam)
	}
	return
}

func (a *OfficialCommonAreaGroup) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = a.InitObjects()
	}
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetPage(page).SetSize(size).SetRecv(recv).List()
	}
	queryParam := a.Param(mw, args...).SetPage(page).SetSize(size).SetRecv(recv)
	if err := DBI.FireReading(a, queryParam); err != nil {
		return nil, err
	}
	cnt, err := queryParam.List()
	if err == nil {
		switch v := recv.(type) {
		case *[]*OfficialCommonAreaGroup:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonAreaGroup(*v))
		case []*OfficialCommonAreaGroup:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonAreaGroup(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCommonAreaGroup) GroupBy(keyField string, inputRows ...[]*OfficialCommonAreaGroup) map[string][]*OfficialCommonAreaGroup {
	var rows Slice_OfficialCommonAreaGroup
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonAreaGroup(inputRows[0])
	} else {
		rows = Slice_OfficialCommonAreaGroup(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialCommonAreaGroup) KeyBy(keyField string, inputRows ...[]*OfficialCommonAreaGroup) map[string]*OfficialCommonAreaGroup {
	var rows Slice_OfficialCommonAreaGroup
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonAreaGroup(inputRows[0])
	} else {
		rows = Slice_OfficialCommonAreaGroup(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialCommonAreaGroup) AsKV(keyField string, valueField string, inputRows ...[]*OfficialCommonAreaGroup) param.Store {
	var rows Slice_OfficialCommonAreaGroup
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonAreaGroup(inputRows[0])
	} else {
		rows = Slice_OfficialCommonAreaGroup(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialCommonAreaGroup) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
	if recv == nil {
		recv = a.InitObjects()
	}
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetOffset(offset).SetSize(size).SetRecv(recv).List()
	}
	queryParam := a.Param(mw, args...).SetOffset(offset).SetSize(size).SetRecv(recv)
	if err := DBI.FireReading(a, queryParam); err != nil {
		return nil, err
	}
	cnt, err := queryParam.List()
	if err == nil {
		switch v := recv.(type) {
		case *[]*OfficialCommonAreaGroup:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonAreaGroup(*v))
		case []*OfficialCommonAreaGroup:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonAreaGroup(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCommonAreaGroup) Insert() (pk interface{}, err error) {
	a.Id = 0
	if a.base.Eventable() {
		err = DBI.Fire("creating", a, nil)
		if err != nil {
			return
		}
	}
	pk, err = a.Param(nil).SetSend(a).Insert()
	if err == nil && pk != nil {
		if v, y := pk.(uint); y {
			a.Id = v
		} else if v, y := pk.(int64); y {
			a.Id = uint(v)
		}
	}
	if err == nil && a.base.Eventable() {
		err = DBI.Fire("created", a, nil)
	}
	return
}

func (a *OfficialCommonAreaGroup) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {

	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(a).Update()
	}
	if err = DBI.Fire("updating", a, mw, args...); err != nil {
		return
	}
	if err = a.Param(mw, args...).SetSend(a).Update(); err != nil {
		return
	}
	return DBI.Fire("updated", a, mw, args...)
}

func (a *OfficialCommonAreaGroup) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(a).Updatex()
	}
	if err = DBI.Fire("updating", a, mw, args...); err != nil {
		return
	}
	if affected, err = a.Param(mw, args...).SetSend(a).Updatex(); err != nil {
		return
	}
	err = DBI.Fire("updated", a, mw, args...)
	return
}

func (a *OfficialCommonAreaGroup) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {

	if !a.base.Eventable() {
		return a.Param(mw, args...).UpdateByStruct(a, fields...)
	}
	editColumns := make([]string, len(fields))
	for index, field := range fields {
		editColumns[index] = com.SnakeCase(field)
	}
	if err = DBI.FireUpdate("updating", a, editColumns, mw, args...); err != nil {
		return
	}
	if err = a.Param(mw, args...).UpdateByStruct(a, fields...); err != nil {
		return
	}
	err = DBI.FireUpdate("updated", a, editColumns, mw, args...)
	return
}

func (a *OfficialCommonAreaGroup) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {

	if !a.base.Eventable() {
		return a.Param(mw, args...).UpdatexByStruct(a, fields...)
	}
	editColumns := make([]string, len(fields))
	for index, field := range fields {
		editColumns[index] = com.SnakeCase(field)
	}
	if err = DBI.FireUpdate("updating", a, editColumns, mw, args...); err != nil {
		return
	}
	if affected, err = a.Param(mw, args...).UpdatexByStruct(a, fields...); err != nil {
		return
	}
	err = DBI.FireUpdate("updated", a, editColumns, mw, args...)
	return
}

func (a *OfficialCommonAreaGroup) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCommonAreaGroup) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCommonAreaGroup) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(kvset).Update()
	}
	m := *a
	m.FromRow(kvset)
	var editColumns []string
	for column := range kvset {
		editColumns = append(editColumns, column)
	}
	if err = DBI.FireUpdate("updating", &m, editColumns, mw, args...); err != nil {
		return
	}
	if err = a.Param(mw, args...).SetSend(kvset).Update(); err != nil {
		return
	}
	return DBI.FireUpdate("updated", &m, editColumns, mw, args...)
}

func (a *OfficialCommonAreaGroup) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(kvset).Updatex()
	}
	m := *a
	m.FromRow(kvset)
	var editColumns []string
	for column := range kvset {
		editColumns = append(editColumns, column)
	}
	if err = DBI.FireUpdate("updating", &m, editColumns, mw, args...); err != nil {
		return
	}
	if affected, err = a.Param(mw, args...).SetSend(kvset).Updatex(); err != nil {
		return
	}
	err = DBI.FireUpdate("updated", &m, editColumns, mw, args...)
	return
}

func (a *OfficialCommonAreaGroup) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(keysValues).Update()
	}
	m := *a
	m.FromRow(keysValues.Map())
	if err = DBI.FireUpdate("updating", &m, keysValues.Keys(), mw, args...); err != nil {
		return
	}
	if err = a.Param(mw, args...).SetSend(keysValues).Update(); err != nil {
		return
	}
	return DBI.FireUpdate("updated", &m, keysValues.Keys(), mw, args...)
}

func (a *OfficialCommonAreaGroup) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Id = 0
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("creating", a, nil)
	})
	if err == nil && pk != nil {
		if v, y := pk.(uint); y {
			a.Id = v
		} else if v, y := pk.(int64); y {
			a.Id = uint(v)
		}
	}
	if err == nil && a.base.Eventable() {
		if pk == nil {
			err = DBI.Fire("updated", a, mw, args...)
		} else {
			err = DBI.Fire("created", a, nil)
		}
	}
	return
}

func (a *OfficialCommonAreaGroup) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

	if !a.base.Eventable() {
		return a.Param(mw, args...).Delete()
	}
	if err = DBI.Fire("deleting", a, mw, args...); err != nil {
		return
	}
	if err = a.Param(mw, args...).Delete(); err != nil {
		return
	}
	return DBI.Fire("deleted", a, mw, args...)
}

func (a *OfficialCommonAreaGroup) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

	if !a.base.Eventable() {
		return a.Param(mw, args...).Deletex()
	}
	if err = DBI.Fire("deleting", a, mw, args...); err != nil {
		return
	}
	if affected, err = a.Param(mw, args...).Deletex(); err != nil {
		return
	}
	err = DBI.Fire("deleted", a, mw, args...)
	return
}

func (a *OfficialCommonAreaGroup) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialCommonAreaGroup) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialCommonAreaGroup) Reset() *OfficialCommonAreaGroup {
	a.Id = 0
	a.CountryAbbr = ``
	a.Name = ``
	a.Abbr = ``
	a.AreaIds = ``
	a.Sort = 0
	return a
}

func (a *OfficialCommonAreaGroup) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["CountryAbbr"] = a.CountryAbbr
		r["Name"] = a.Name
		r["Abbr"] = a.Abbr
		r["AreaIds"] = a.AreaIds
		r["Sort"] = a.Sort
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "CountryAbbr":
			r["CountryAbbr"] = a.CountryAbbr
		case "Name":
			r["Name"] = a.Name
		case "Abbr":
			r["Abbr"] = a.Abbr
		case "AreaIds":
			r["AreaIds"] = a.AreaIds
		case "Sort":
			r["Sort"] = a.Sort
		}
	}
	return r
}

func (a *OfficialCommonAreaGroup) FromRow(row map[string]interface{}) {
	for key, value := range row {
		switch key {
		case "id":
			a.Id = param.AsUint(value)
		case "country_abbr":
			a.CountryAbbr = param.AsString(value)
		case "name":
			a.Name = param.AsString(value)
		case "abbr":
			a.Abbr = param.AsString(value)
		case "area_ids":
			a.AreaIds = param.AsString(value)
		case "sort":
			a.Sort = param.AsInt(value)
		}
	}
}

func (a *OfficialCommonAreaGroup) Set(key interface{}, value ...interface{}) {
	switch k := key.(type) {
	case map[string]interface{}:
		for kk, vv := range k {
			a.Set(kk, vv)
		}
	default:
		var (
			kk string
			vv interface{}
		)
		if k, y := key.(string); y {
			kk = k
		} else {
			kk = fmt.Sprint(key)
		}
		if len(value) > 0 {
			vv = value[0]
		}
		switch kk {
		case "Id":
			a.Id = param.AsUint(vv)
		case "CountryAbbr":
			a.CountryAbbr = param.AsString(vv)
		case "Name":
			a.Name = param.AsString(vv)
		case "Abbr":
			a.Abbr = param.AsString(vv)
		case "AreaIds":
			a.AreaIds = param.AsString(vv)
		case "Sort":
			a.Sort = param.AsInt(vv)
		}
	}
}

func (a *OfficialCommonAreaGroup) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["country_abbr"] = a.CountryAbbr
		r["name"] = a.Name
		r["abbr"] = a.Abbr
		r["area_ids"] = a.AreaIds
		r["sort"] = a.Sort
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "country_abbr":
			r["country_abbr"] = a.CountryAbbr
		case "name":
			r["name"] = a.Name
		case "abbr":
			r["abbr"] = a.Abbr
		case "area_ids":
			r["area_ids"] = a.AreaIds
		case "sort":
			r["sort"] = a.Sort
		}
	}
	return r
}

func (a *OfficialCommonAreaGroup) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCommonAreaGroup) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, recv, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCommonAreaGroup) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialCommonAreaGroup) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
