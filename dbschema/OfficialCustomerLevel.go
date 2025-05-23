// @generated Do not edit this file, which is automatically generated by the generator.

package dbschema

import (
	"fmt"

	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type Slice_OfficialCustomerLevel []*OfficialCustomerLevel

func (s Slice_OfficialCustomerLevel) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerLevel) RangeRaw(fn func(m *OfficialCustomerLevel) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerLevel) GroupBy(keyField string) map[string][]*OfficialCustomerLevel {
	r := map[string][]*OfficialCustomerLevel{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialCustomerLevel{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialCustomerLevel) KeyBy(keyField string) map[string]*OfficialCustomerLevel {
	r := map[string]*OfficialCustomerLevel{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialCustomerLevel) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialCustomerLevel) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialCustomerLevel) FromList(data interface{}) Slice_OfficialCustomerLevel {
	values, ok := data.([]*OfficialCustomerLevel)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialCustomerLevel{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialCustomerLevel(ctx echo.Context) *OfficialCustomerLevel {
	m := &OfficialCustomerLevel{}
	m.SetContext(ctx)
	return m
}

// OfficialCustomerLevel 客户等级
type OfficialCustomerLevel struct {
	base    factory.Base
	objects []*OfficialCustomerLevel

	Id                 uint    `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	Name               string  `db:"name" bson:"name" comment:"等级名称" json:"name" xml:"name"`
	Short              string  `db:"short" bson:"short" comment:"等级简称" json:"short" xml:"short"`
	Description        string  `db:"description" bson:"description" comment:"等级简介" json:"description" xml:"description"`
	IconImage          string  `db:"icon_image" bson:"icon_image" comment:"图标图片" json:"icon_image" xml:"icon_image"`
	IconClass          string  `db:"icon_class" bson:"icon_class" comment:"图片class名" json:"icon_class" xml:"icon_class"`
	Color              string  `db:"color" bson:"color" comment:"颜色" json:"color" xml:"color"`
	Bgcolor            string  `db:"bgcolor" bson:"bgcolor" comment:"背景色" json:"bgcolor" xml:"bgcolor"`
	Price              float64 `db:"price" bson:"price" comment:"升级价格(0为免费)" json:"price" xml:"price"`
	IntegralAsset      string  `db:"integral_asset" bson:"integral_asset" comment:"当作升级积分的资产" json:"integral_asset" xml:"integral_asset"`
	IntegralAmountType string  `db:"integral_amount_type" bson:"integral_amount_type" comment:"资产金额类型(balance-余额;accumulated-累积额)" json:"integral_amount_type" xml:"integral_amount_type"`
	IntegralMin        float64 `db:"integral_min" bson:"integral_min" comment:"最小积分" json:"integral_min" xml:"integral_min"`
	IntegralMax        float64 `db:"integral_max" bson:"integral_max" comment:"最大积分" json:"integral_max" xml:"integral_max"`
	Created            uint    `db:"created" bson:"created" comment:"添加时间" json:"created" xml:"created"`
	Updated            uint    `db:"updated" bson:"updated" comment:"更新时间" json:"updated" xml:"updated"`
	Score              int     `db:"score" bson:"score" comment:"分值(分值越大等级越高)" json:"score" xml:"score"`
	Disabled           string  `db:"disabled" bson:"disabled" comment:"是否(Y/N)禁用" json:"disabled" xml:"disabled"`
	Extra              string  `db:"extra" bson:"extra" comment:"扩展配置(JSON)" json:"extra" xml:"extra"`
	Group              string  `db:"group" bson:"group" comment:"扩展组(base-基础组,其它名称为扩展组。客户只能有一个基础组等级,可以有多个扩展组等级)" json:"group" xml:"group"`
	RoleIds            string  `db:"role_ids" bson:"role_ids" comment:"角色ID(多个用“,”分隔开)" json:"role_ids" xml:"role_ids"`
}

// - base function

func (a *OfficialCustomerLevel) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialCustomerLevel) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialCustomerLevel) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialCustomerLevel) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialCustomerLevel) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialCustomerLevel) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialCustomerLevel) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialCustomerLevel) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialCustomerLevel) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialCustomerLevel) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialCustomerLevel) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialCustomerLevel) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialCustomerLevel) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialCustomerLevel) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialCustomerLevel) Objects() []*OfficialCustomerLevel {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialCustomerLevel) XObjects() Slice_OfficialCustomerLevel {
	return Slice_OfficialCustomerLevel(a.Objects())
}

func (a *OfficialCustomerLevel) NewObjects() factory.Ranger {
	return &Slice_OfficialCustomerLevel{}
}

func (a *OfficialCustomerLevel) InitObjects() *[]*OfficialCustomerLevel {
	a.objects = []*OfficialCustomerLevel{}
	return &a.objects
}

func (a *OfficialCustomerLevel) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialCustomerLevel) Short_() string {
	return "official_customer_level"
}

func (a *OfficialCustomerLevel) Struct_() string {
	return "OfficialCustomerLevel"
}

func (a *OfficialCustomerLevel) Name_() string {
	b := a
	if b == nil {
		b = &OfficialCustomerLevel{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialCustomerLevel) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialCustomerLevel) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialCustomerLevel) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerLevel:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerLevel(*v))
		case []*OfficialCustomerLevel:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerLevel(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerLevel) GroupBy(keyField string, inputRows ...[]*OfficialCustomerLevel) map[string][]*OfficialCustomerLevel {
	var rows Slice_OfficialCustomerLevel
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerLevel(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerLevel(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialCustomerLevel) KeyBy(keyField string, inputRows ...[]*OfficialCustomerLevel) map[string]*OfficialCustomerLevel {
	var rows Slice_OfficialCustomerLevel
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerLevel(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerLevel(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialCustomerLevel) AsKV(keyField string, valueField string, inputRows ...[]*OfficialCustomerLevel) param.Store {
	var rows Slice_OfficialCustomerLevel
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerLevel(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerLevel(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialCustomerLevel) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerLevel:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerLevel(*v))
		case []*OfficialCustomerLevel:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerLevel(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerLevel) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.IntegralAsset) == 0 {
		a.IntegralAsset = "integral"
	}
	if len(a.IntegralAmountType) == 0 {
		a.IntegralAmountType = "balance"
	}
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.Group) == 0 {
		a.Group = "base"
	}
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

func (a *OfficialCustomerLevel) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.IntegralAsset) == 0 {
		a.IntegralAsset = "integral"
	}
	if len(a.IntegralAmountType) == 0 {
		a.IntegralAmountType = "balance"
	}
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.Group) == 0 {
		a.Group = "base"
	}
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

func (a *OfficialCustomerLevel) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.IntegralAsset) == 0 {
		a.IntegralAsset = "integral"
	}
	if len(a.IntegralAmountType) == 0 {
		a.IntegralAmountType = "balance"
	}
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.Group) == 0 {
		a.Group = "base"
	}
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

func (a *OfficialCustomerLevel) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.IntegralAsset) == 0 {
		a.IntegralAsset = "integral"
	}
	if len(a.IntegralAmountType) == 0 {
		a.IntegralAmountType = "balance"
	}
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.Group) == 0 {
		a.Group = "base"
	}
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

func (a *OfficialCustomerLevel) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.IntegralAsset) == 0 {
		a.IntegralAsset = "integral"
	}
	if len(a.IntegralAmountType) == 0 {
		a.IntegralAmountType = "balance"
	}
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.Group) == 0 {
		a.Group = "base"
	}
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

func (a *OfficialCustomerLevel) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerLevel) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerLevel) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["integral_asset"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["integral_asset"] = "integral"
		}
	}
	if val, ok := kvset["integral_amount_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["integral_amount_type"] = "balance"
		}
	}
	if val, ok := kvset["disabled"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["disabled"] = "N"
		}
	}
	if val, ok := kvset["group"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["group"] = "base"
		}
	}
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(kvset).Update()
	}
	m := *a
	m.FromRow(kvset)
	editColumns := make([]string, 0, len(kvset))
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

func (a *OfficialCustomerLevel) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["integral_asset"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["integral_asset"] = "integral"
		}
	}
	if val, ok := kvset["integral_amount_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["integral_amount_type"] = "balance"
		}
	}
	if val, ok := kvset["disabled"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["disabled"] = "N"
		}
	}
	if val, ok := kvset["group"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["group"] = "base"
		}
	}
	if !a.base.Eventable() {
		return a.Param(mw, args...).SetSend(kvset).Updatex()
	}
	m := *a
	m.FromRow(kvset)
	editColumns := make([]string, 0, len(kvset))
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

func (a *OfficialCustomerLevel) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialCustomerLevel) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		a.Updated = uint(time.Now().Unix())
		if len(a.IntegralAsset) == 0 {
			a.IntegralAsset = "integral"
		}
		if len(a.IntegralAmountType) == 0 {
			a.IntegralAmountType = "balance"
		}
		if len(a.Disabled) == 0 {
			a.Disabled = "N"
		}
		if len(a.Group) == 0 {
			a.Group = "base"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		a.Id = 0
		if len(a.IntegralAsset) == 0 {
			a.IntegralAsset = "integral"
		}
		if len(a.IntegralAmountType) == 0 {
			a.IntegralAmountType = "balance"
		}
		if len(a.Disabled) == 0 {
			a.Disabled = "N"
		}
		if len(a.Group) == 0 {
			a.Group = "base"
		}
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

func (a *OfficialCustomerLevel) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialCustomerLevel) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialCustomerLevel) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialCustomerLevel) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialCustomerLevel) Reset() *OfficialCustomerLevel {
	a.Id = 0
	a.Name = ``
	a.Short = ``
	a.Description = ``
	a.IconImage = ``
	a.IconClass = ``
	a.Color = ``
	a.Bgcolor = ``
	a.Price = 0.0
	a.IntegralAsset = ``
	a.IntegralAmountType = ``
	a.IntegralMin = 0.0
	a.IntegralMax = 0.0
	a.Created = 0
	a.Updated = 0
	a.Score = 0
	a.Disabled = ``
	a.Extra = ``
	a.Group = ``
	a.RoleIds = ``
	return a
}

func (a *OfficialCustomerLevel) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["Name"] = a.Name
		r["Short"] = a.Short
		r["Description"] = a.Description
		r["IconImage"] = a.IconImage
		r["IconClass"] = a.IconClass
		r["Color"] = a.Color
		r["Bgcolor"] = a.Bgcolor
		r["Price"] = a.Price
		r["IntegralAsset"] = a.IntegralAsset
		r["IntegralAmountType"] = a.IntegralAmountType
		r["IntegralMin"] = a.IntegralMin
		r["IntegralMax"] = a.IntegralMax
		r["Created"] = a.Created
		r["Updated"] = a.Updated
		r["Score"] = a.Score
		r["Disabled"] = a.Disabled
		r["Extra"] = a.Extra
		r["Group"] = a.Group
		r["RoleIds"] = a.RoleIds
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "Name":
			r["Name"] = a.Name
		case "Short":
			r["Short"] = a.Short
		case "Description":
			r["Description"] = a.Description
		case "IconImage":
			r["IconImage"] = a.IconImage
		case "IconClass":
			r["IconClass"] = a.IconClass
		case "Color":
			r["Color"] = a.Color
		case "Bgcolor":
			r["Bgcolor"] = a.Bgcolor
		case "Price":
			r["Price"] = a.Price
		case "IntegralAsset":
			r["IntegralAsset"] = a.IntegralAsset
		case "IntegralAmountType":
			r["IntegralAmountType"] = a.IntegralAmountType
		case "IntegralMin":
			r["IntegralMin"] = a.IntegralMin
		case "IntegralMax":
			r["IntegralMax"] = a.IntegralMax
		case "Created":
			r["Created"] = a.Created
		case "Updated":
			r["Updated"] = a.Updated
		case "Score":
			r["Score"] = a.Score
		case "Disabled":
			r["Disabled"] = a.Disabled
		case "Extra":
			r["Extra"] = a.Extra
		case "Group":
			r["Group"] = a.Group
		case "RoleIds":
			r["RoleIds"] = a.RoleIds
		}
	}
	return r
}

func (a *OfficialCustomerLevel) FromRow(row map[string]interface{}) {
	for key, value := range row {
		if _, ok := value.(db.RawValue); ok {
			continue
		}
		switch key {
		case "id":
			a.Id = param.AsUint(value)
		case "name":
			a.Name = param.AsString(value)
		case "short":
			a.Short = param.AsString(value)
		case "description":
			a.Description = param.AsString(value)
		case "icon_image":
			a.IconImage = param.AsString(value)
		case "icon_class":
			a.IconClass = param.AsString(value)
		case "color":
			a.Color = param.AsString(value)
		case "bgcolor":
			a.Bgcolor = param.AsString(value)
		case "price":
			a.Price = param.AsFloat64(value)
		case "integral_asset":
			a.IntegralAsset = param.AsString(value)
		case "integral_amount_type":
			a.IntegralAmountType = param.AsString(value)
		case "integral_min":
			a.IntegralMin = param.AsFloat64(value)
		case "integral_max":
			a.IntegralMax = param.AsFloat64(value)
		case "created":
			a.Created = param.AsUint(value)
		case "updated":
			a.Updated = param.AsUint(value)
		case "score":
			a.Score = param.AsInt(value)
		case "disabled":
			a.Disabled = param.AsString(value)
		case "extra":
			a.Extra = param.AsString(value)
		case "group":
			a.Group = param.AsString(value)
		case "role_ids":
			a.RoleIds = param.AsString(value)
		}
	}
}

func (a *OfficialCustomerLevel) GetField(field string) interface{} {
	switch field {
	case "Id":
		return a.Id
	case "Name":
		return a.Name
	case "Short":
		return a.Short
	case "Description":
		return a.Description
	case "IconImage":
		return a.IconImage
	case "IconClass":
		return a.IconClass
	case "Color":
		return a.Color
	case "Bgcolor":
		return a.Bgcolor
	case "Price":
		return a.Price
	case "IntegralAsset":
		return a.IntegralAsset
	case "IntegralAmountType":
		return a.IntegralAmountType
	case "IntegralMin":
		return a.IntegralMin
	case "IntegralMax":
		return a.IntegralMax
	case "Created":
		return a.Created
	case "Updated":
		return a.Updated
	case "Score":
		return a.Score
	case "Disabled":
		return a.Disabled
	case "Extra":
		return a.Extra
	case "Group":
		return a.Group
	case "RoleIds":
		return a.RoleIds
	default:
		return nil
	}
}

func (a *OfficialCustomerLevel) GetAllFieldNames() []string {
	return []string{
		"Id",
		"Name",
		"Short",
		"Description",
		"IconImage",
		"IconClass",
		"Color",
		"Bgcolor",
		"Price",
		"IntegralAsset",
		"IntegralAmountType",
		"IntegralMin",
		"IntegralMax",
		"Created",
		"Updated",
		"Score",
		"Disabled",
		"Extra",
		"Group",
		"RoleIds",
	}
}

func (a *OfficialCustomerLevel) HasField(field string) bool {
	switch field {
	case "Id":
		return true
	case "Name":
		return true
	case "Short":
		return true
	case "Description":
		return true
	case "IconImage":
		return true
	case "IconClass":
		return true
	case "Color":
		return true
	case "Bgcolor":
		return true
	case "Price":
		return true
	case "IntegralAsset":
		return true
	case "IntegralAmountType":
		return true
	case "IntegralMin":
		return true
	case "IntegralMax":
		return true
	case "Created":
		return true
	case "Updated":
		return true
	case "Score":
		return true
	case "Disabled":
		return true
	case "Extra":
		return true
	case "Group":
		return true
	case "RoleIds":
		return true
	default:
		return false
	}
}

func (a *OfficialCustomerLevel) Set(key interface{}, value ...interface{}) {
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
		case "Name":
			a.Name = param.AsString(vv)
		case "Short":
			a.Short = param.AsString(vv)
		case "Description":
			a.Description = param.AsString(vv)
		case "IconImage":
			a.IconImage = param.AsString(vv)
		case "IconClass":
			a.IconClass = param.AsString(vv)
		case "Color":
			a.Color = param.AsString(vv)
		case "Bgcolor":
			a.Bgcolor = param.AsString(vv)
		case "Price":
			a.Price = param.AsFloat64(vv)
		case "IntegralAsset":
			a.IntegralAsset = param.AsString(vv)
		case "IntegralAmountType":
			a.IntegralAmountType = param.AsString(vv)
		case "IntegralMin":
			a.IntegralMin = param.AsFloat64(vv)
		case "IntegralMax":
			a.IntegralMax = param.AsFloat64(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "Updated":
			a.Updated = param.AsUint(vv)
		case "Score":
			a.Score = param.AsInt(vv)
		case "Disabled":
			a.Disabled = param.AsString(vv)
		case "Extra":
			a.Extra = param.AsString(vv)
		case "Group":
			a.Group = param.AsString(vv)
		case "RoleIds":
			a.RoleIds = param.AsString(vv)
		}
	}
}

func (a *OfficialCustomerLevel) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["name"] = a.Name
		r["short"] = a.Short
		r["description"] = a.Description
		r["icon_image"] = a.IconImage
		r["icon_class"] = a.IconClass
		r["color"] = a.Color
		r["bgcolor"] = a.Bgcolor
		r["price"] = a.Price
		r["integral_asset"] = a.IntegralAsset
		r["integral_amount_type"] = a.IntegralAmountType
		r["integral_min"] = a.IntegralMin
		r["integral_max"] = a.IntegralMax
		r["created"] = a.Created
		r["updated"] = a.Updated
		r["score"] = a.Score
		r["disabled"] = a.Disabled
		r["extra"] = a.Extra
		r["group"] = a.Group
		r["role_ids"] = a.RoleIds
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "name":
			r["name"] = a.Name
		case "short":
			r["short"] = a.Short
		case "description":
			r["description"] = a.Description
		case "icon_image":
			r["icon_image"] = a.IconImage
		case "icon_class":
			r["icon_class"] = a.IconClass
		case "color":
			r["color"] = a.Color
		case "bgcolor":
			r["bgcolor"] = a.Bgcolor
		case "price":
			r["price"] = a.Price
		case "integral_asset":
			r["integral_asset"] = a.IntegralAsset
		case "integral_amount_type":
			r["integral_amount_type"] = a.IntegralAmountType
		case "integral_min":
			r["integral_min"] = a.IntegralMin
		case "integral_max":
			r["integral_max"] = a.IntegralMax
		case "created":
			r["created"] = a.Created
		case "updated":
			r["updated"] = a.Updated
		case "score":
			r["score"] = a.Score
		case "disabled":
			r["disabled"] = a.Disabled
		case "extra":
			r["extra"] = a.Extra
		case "group":
			r["group"] = a.Group
		case "role_ids":
			r["role_ids"] = a.RoleIds
		}
	}
	return r
}

func (a *OfficialCustomerLevel) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPage(a, cond, sorts...)
}

func (a *OfficialCustomerLevel) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageAs(a, recv, cond, sorts...)
}

func (a *OfficialCustomerLevel) ListPageByOffset(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffset(a, cond, sorts...)
}

func (a *OfficialCustomerLevel) ListPageByOffsetAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffsetAs(a, recv, cond, sorts...)
}

func (a *OfficialCustomerLevel) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialCustomerLevel) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
