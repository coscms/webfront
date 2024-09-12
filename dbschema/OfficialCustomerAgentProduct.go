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

type Slice_OfficialCustomerAgentProduct []*OfficialCustomerAgentProduct

func (s Slice_OfficialCustomerAgentProduct) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerAgentProduct) RangeRaw(fn func(m *OfficialCustomerAgentProduct) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerAgentProduct) GroupBy(keyField string) map[string][]*OfficialCustomerAgentProduct {
	r := map[string][]*OfficialCustomerAgentProduct{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialCustomerAgentProduct{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialCustomerAgentProduct) KeyBy(keyField string) map[string]*OfficialCustomerAgentProduct {
	r := map[string]*OfficialCustomerAgentProduct{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialCustomerAgentProduct) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialCustomerAgentProduct) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialCustomerAgentProduct) FromList(data interface{}) Slice_OfficialCustomerAgentProduct {
	values, ok := data.([]*OfficialCustomerAgentProduct)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialCustomerAgentProduct{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialCustomerAgentProduct(ctx echo.Context) *OfficialCustomerAgentProduct {
	m := &OfficialCustomerAgentProduct{}
	m.SetContext(ctx)
	return m
}

// OfficialCustomerAgentProduct 代理产品列表
type OfficialCustomerAgentProduct struct {
	base    factory.Base
	objects []*OfficialCustomerAgentProduct

	Id           uint64  `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	AgentId      uint64  `db:"agent_id" bson:"agent_id" comment:"代理商UID(customer表id)" json:"agent_id" xml:"agent_id"`
	ProductId    string  `db:"product_id" bson:"product_id" comment:"商品ID" json:"product_id" xml:"product_id"`
	ProductTable string  `db:"product_table" bson:"product_table" comment:"商品表名称(不含official_前缀)" json:"product_table" xml:"product_table"`
	Sold         uint64  `db:"sold" bson:"sold" comment:"销量" json:"sold" xml:"sold"`
	Performance  float64 `db:"performance" bson:"performance" comment:"业绩" json:"performance" xml:"performance"`
	Created      uint    `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Expired      uint    `db:"expired" bson:"expired" comment:"过期时间" json:"expired" xml:"expired"`
	Updated      uint    `db:"updated" bson:"updated" comment:"更新时间" json:"updated" xml:"updated"`
	Disabled     string  `db:"disabled" bson:"disabled" comment:"是否禁用" json:"disabled" xml:"disabled"`
}

// - base function

func (a *OfficialCustomerAgentProduct) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialCustomerAgentProduct) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialCustomerAgentProduct) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialCustomerAgentProduct) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialCustomerAgentProduct) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialCustomerAgentProduct) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialCustomerAgentProduct) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialCustomerAgentProduct) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialCustomerAgentProduct) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialCustomerAgentProduct) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialCustomerAgentProduct) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialCustomerAgentProduct) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialCustomerAgentProduct) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialCustomerAgentProduct) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialCustomerAgentProduct) Objects() []*OfficialCustomerAgentProduct {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialCustomerAgentProduct) XObjects() Slice_OfficialCustomerAgentProduct {
	return Slice_OfficialCustomerAgentProduct(a.Objects())
}

func (a *OfficialCustomerAgentProduct) NewObjects() factory.Ranger {
	return &Slice_OfficialCustomerAgentProduct{}
}

func (a *OfficialCustomerAgentProduct) InitObjects() *[]*OfficialCustomerAgentProduct {
	a.objects = []*OfficialCustomerAgentProduct{}
	return &a.objects
}

func (a *OfficialCustomerAgentProduct) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialCustomerAgentProduct) Short_() string {
	return "official_customer_agent_product"
}

func (a *OfficialCustomerAgentProduct) Struct_() string {
	return "OfficialCustomerAgentProduct"
}

func (a *OfficialCustomerAgentProduct) Name_() string {
	b := a
	if b == nil {
		b = &OfficialCustomerAgentProduct{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialCustomerAgentProduct) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialCustomerAgentProduct) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialCustomerAgentProduct) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerAgentProduct:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerAgentProduct(*v))
		case []*OfficialCustomerAgentProduct:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerAgentProduct(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerAgentProduct) GroupBy(keyField string, inputRows ...[]*OfficialCustomerAgentProduct) map[string][]*OfficialCustomerAgentProduct {
	var rows Slice_OfficialCustomerAgentProduct
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerAgentProduct(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerAgentProduct(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialCustomerAgentProduct) KeyBy(keyField string, inputRows ...[]*OfficialCustomerAgentProduct) map[string]*OfficialCustomerAgentProduct {
	var rows Slice_OfficialCustomerAgentProduct
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerAgentProduct(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerAgentProduct(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialCustomerAgentProduct) AsKV(keyField string, valueField string, inputRows ...[]*OfficialCustomerAgentProduct) param.Store {
	var rows Slice_OfficialCustomerAgentProduct
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerAgentProduct(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerAgentProduct(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialCustomerAgentProduct) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerAgentProduct:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerAgentProduct(*v))
		case []*OfficialCustomerAgentProduct:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerAgentProduct(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerAgentProduct) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if a.base.Eventable() {
		err = DBI.Fire("creating", a, nil)
		if err != nil {
			return
		}
	}
	pk, err = a.Param(nil).SetSend(a).Insert()
	if err == nil && pk != nil {
		if v, y := pk.(uint64); y {
			a.Id = v
		} else if v, y := pk.(int64); y {
			a.Id = uint64(v)
		}
	}
	if err == nil && a.base.Eventable() {
		err = DBI.Fire("created", a, nil)
	}
	return
}

func (a *OfficialCustomerAgentProduct) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
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

func (a *OfficialCustomerAgentProduct) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
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

func (a *OfficialCustomerAgentProduct) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
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

func (a *OfficialCustomerAgentProduct) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
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

func (a *OfficialCustomerAgentProduct) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerAgentProduct) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerAgentProduct) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["disabled"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["disabled"] = "N"
		}
	}
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

func (a *OfficialCustomerAgentProduct) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["disabled"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["disabled"] = "N"
		}
	}
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

func (a *OfficialCustomerAgentProduct) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialCustomerAgentProduct) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		a.Updated = uint(time.Now().Unix())
		if len(a.Disabled) == 0 {
			a.Disabled = "N"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		a.Id = 0
		if len(a.Disabled) == 0 {
			a.Disabled = "N"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("creating", a, nil)
	})
	if err == nil && pk != nil {
		if v, y := pk.(uint64); y {
			a.Id = v
		} else if v, y := pk.(int64); y {
			a.Id = uint64(v)
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

func (a *OfficialCustomerAgentProduct) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialCustomerAgentProduct) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialCustomerAgentProduct) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialCustomerAgentProduct) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialCustomerAgentProduct) Reset() *OfficialCustomerAgentProduct {
	a.Id = 0
	a.AgentId = 0
	a.ProductId = ``
	a.ProductTable = ``
	a.Sold = 0
	a.Performance = 0.0
	a.Created = 0
	a.Expired = 0
	a.Updated = 0
	a.Disabled = ``
	return a
}

func (a *OfficialCustomerAgentProduct) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["AgentId"] = a.AgentId
		r["ProductId"] = a.ProductId
		r["ProductTable"] = a.ProductTable
		r["Sold"] = a.Sold
		r["Performance"] = a.Performance
		r["Created"] = a.Created
		r["Expired"] = a.Expired
		r["Updated"] = a.Updated
		r["Disabled"] = a.Disabled
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "AgentId":
			r["AgentId"] = a.AgentId
		case "ProductId":
			r["ProductId"] = a.ProductId
		case "ProductTable":
			r["ProductTable"] = a.ProductTable
		case "Sold":
			r["Sold"] = a.Sold
		case "Performance":
			r["Performance"] = a.Performance
		case "Created":
			r["Created"] = a.Created
		case "Expired":
			r["Expired"] = a.Expired
		case "Updated":
			r["Updated"] = a.Updated
		case "Disabled":
			r["Disabled"] = a.Disabled
		}
	}
	return r
}

func (a *OfficialCustomerAgentProduct) FromRow(row map[string]interface{}) {
	for key, value := range row {
		switch key {
		case "id":
			a.Id = param.AsUint64(value)
		case "agent_id":
			a.AgentId = param.AsUint64(value)
		case "product_id":
			a.ProductId = param.AsString(value)
		case "product_table":
			a.ProductTable = param.AsString(value)
		case "sold":
			a.Sold = param.AsUint64(value)
		case "performance":
			a.Performance = param.AsFloat64(value)
		case "created":
			a.Created = param.AsUint(value)
		case "expired":
			a.Expired = param.AsUint(value)
		case "updated":
			a.Updated = param.AsUint(value)
		case "disabled":
			a.Disabled = param.AsString(value)
		}
	}
}

func (a *OfficialCustomerAgentProduct) Set(key interface{}, value ...interface{}) {
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
			a.Id = param.AsUint64(vv)
		case "AgentId":
			a.AgentId = param.AsUint64(vv)
		case "ProductId":
			a.ProductId = param.AsString(vv)
		case "ProductTable":
			a.ProductTable = param.AsString(vv)
		case "Sold":
			a.Sold = param.AsUint64(vv)
		case "Performance":
			a.Performance = param.AsFloat64(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "Expired":
			a.Expired = param.AsUint(vv)
		case "Updated":
			a.Updated = param.AsUint(vv)
		case "Disabled":
			a.Disabled = param.AsString(vv)
		}
	}
}

func (a *OfficialCustomerAgentProduct) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["agent_id"] = a.AgentId
		r["product_id"] = a.ProductId
		r["product_table"] = a.ProductTable
		r["sold"] = a.Sold
		r["performance"] = a.Performance
		r["created"] = a.Created
		r["expired"] = a.Expired
		r["updated"] = a.Updated
		r["disabled"] = a.Disabled
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "agent_id":
			r["agent_id"] = a.AgentId
		case "product_id":
			r["product_id"] = a.ProductId
		case "product_table":
			r["product_table"] = a.ProductTable
		case "sold":
			r["sold"] = a.Sold
		case "performance":
			r["performance"] = a.Performance
		case "created":
			r["created"] = a.Created
		case "expired":
			r["expired"] = a.Expired
		case "updated":
			r["updated"] = a.Updated
		case "disabled":
			r["disabled"] = a.Disabled
		}
	}
	return r
}

func (a *OfficialCustomerAgentProduct) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCustomerAgentProduct) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, recv, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCustomerAgentProduct) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialCustomerAgentProduct) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
