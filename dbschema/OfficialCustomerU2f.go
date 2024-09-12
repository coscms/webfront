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

type Slice_OfficialCustomerU2f []*OfficialCustomerU2f

func (s Slice_OfficialCustomerU2f) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerU2f) RangeRaw(fn func(m *OfficialCustomerU2f) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerU2f) GroupBy(keyField string) map[string][]*OfficialCustomerU2f {
	r := map[string][]*OfficialCustomerU2f{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialCustomerU2f{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialCustomerU2f) KeyBy(keyField string) map[string]*OfficialCustomerU2f {
	r := map[string]*OfficialCustomerU2f{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialCustomerU2f) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialCustomerU2f) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialCustomerU2f) FromList(data interface{}) Slice_OfficialCustomerU2f {
	values, ok := data.([]*OfficialCustomerU2f)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialCustomerU2f{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialCustomerU2f(ctx echo.Context) *OfficialCustomerU2f {
	m := &OfficialCustomerU2f{}
	m.SetContext(ctx)
	return m
}

// OfficialCustomerU2f 两步验证
type OfficialCustomerU2f struct {
	base    factory.Base
	objects []*OfficialCustomerU2f

	Id           uint64 `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	CustomerId   uint64 `db:"customer_id" bson:"customer_id" comment:"客户ID" json:"customer_id" xml:"customer_id"`
	Name         string `db:"name" bson:"name" comment:"名称" json:"name" xml:"name"`
	Token        string `db:"token" bson:"token" comment:"签名" json:"token" xml:"token"`
	Type         string `db:"type" bson:"type" comment:"类型" json:"type" xml:"type"`
	Extra        string `db:"extra" bson:"extra" comment:"扩展设置" json:"extra" xml:"extra"`
	Step         uint   `db:"step" bson:"step" comment:"第几步" json:"step" xml:"step"`
	Precondition string `db:"precondition" bson:"precondition" comment:"除了密码登录外的其它前置条件(仅step=2时有效),用半角逗号分隔" json:"precondition" xml:"precondition"`
	Created      uint   `db:"created" bson:"created" comment:"绑定时间" json:"created" xml:"created"`
}

// - base function

func (a *OfficialCustomerU2f) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialCustomerU2f) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialCustomerU2f) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialCustomerU2f) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialCustomerU2f) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialCustomerU2f) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialCustomerU2f) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialCustomerU2f) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialCustomerU2f) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialCustomerU2f) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialCustomerU2f) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialCustomerU2f) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialCustomerU2f) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialCustomerU2f) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialCustomerU2f) Objects() []*OfficialCustomerU2f {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialCustomerU2f) XObjects() Slice_OfficialCustomerU2f {
	return Slice_OfficialCustomerU2f(a.Objects())
}

func (a *OfficialCustomerU2f) NewObjects() factory.Ranger {
	return &Slice_OfficialCustomerU2f{}
}

func (a *OfficialCustomerU2f) InitObjects() *[]*OfficialCustomerU2f {
	a.objects = []*OfficialCustomerU2f{}
	return &a.objects
}

func (a *OfficialCustomerU2f) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialCustomerU2f) Short_() string {
	return "official_customer_u2f"
}

func (a *OfficialCustomerU2f) Struct_() string {
	return "OfficialCustomerU2f"
}

func (a *OfficialCustomerU2f) Name_() string {
	b := a
	if b == nil {
		b = &OfficialCustomerU2f{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialCustomerU2f) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialCustomerU2f) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialCustomerU2f) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerU2f:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerU2f(*v))
		case []*OfficialCustomerU2f:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerU2f(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerU2f) GroupBy(keyField string, inputRows ...[]*OfficialCustomerU2f) map[string][]*OfficialCustomerU2f {
	var rows Slice_OfficialCustomerU2f
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerU2f(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerU2f(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialCustomerU2f) KeyBy(keyField string, inputRows ...[]*OfficialCustomerU2f) map[string]*OfficialCustomerU2f {
	var rows Slice_OfficialCustomerU2f
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerU2f(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerU2f(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialCustomerU2f) AsKV(keyField string, valueField string, inputRows ...[]*OfficialCustomerU2f) param.Store {
	var rows Slice_OfficialCustomerU2f
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerU2f(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerU2f(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialCustomerU2f) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerU2f:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerU2f(*v))
		case []*OfficialCustomerU2f:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerU2f(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerU2f) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.Precondition) == 0 {
		a.Precondition = "password"
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

func (a *OfficialCustomerU2f) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {

	if len(a.Precondition) == 0 {
		a.Precondition = "password"
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

func (a *OfficialCustomerU2f) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

	if len(a.Precondition) == 0 {
		a.Precondition = "password"
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

func (a *OfficialCustomerU2f) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {

	if len(a.Precondition) == 0 {
		a.Precondition = "password"
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

func (a *OfficialCustomerU2f) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {

	if len(a.Precondition) == 0 {
		a.Precondition = "password"
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

func (a *OfficialCustomerU2f) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerU2f) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerU2f) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["precondition"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["precondition"] = "password"
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

func (a *OfficialCustomerU2f) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["precondition"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["precondition"] = "password"
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

func (a *OfficialCustomerU2f) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialCustomerU2f) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		if len(a.Precondition) == 0 {
			a.Precondition = "password"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		a.Id = 0
		if len(a.Precondition) == 0 {
			a.Precondition = "password"
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

func (a *OfficialCustomerU2f) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialCustomerU2f) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialCustomerU2f) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialCustomerU2f) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialCustomerU2f) Reset() *OfficialCustomerU2f {
	a.Id = 0
	a.CustomerId = 0
	a.Name = ``
	a.Token = ``
	a.Type = ``
	a.Extra = ``
	a.Step = 0
	a.Precondition = ``
	a.Created = 0
	return a
}

func (a *OfficialCustomerU2f) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["CustomerId"] = a.CustomerId
		r["Name"] = a.Name
		r["Token"] = a.Token
		r["Type"] = a.Type
		r["Extra"] = a.Extra
		r["Step"] = a.Step
		r["Precondition"] = a.Precondition
		r["Created"] = a.Created
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "CustomerId":
			r["CustomerId"] = a.CustomerId
		case "Name":
			r["Name"] = a.Name
		case "Token":
			r["Token"] = a.Token
		case "Type":
			r["Type"] = a.Type
		case "Extra":
			r["Extra"] = a.Extra
		case "Step":
			r["Step"] = a.Step
		case "Precondition":
			r["Precondition"] = a.Precondition
		case "Created":
			r["Created"] = a.Created
		}
	}
	return r
}

func (a *OfficialCustomerU2f) FromRow(row map[string]interface{}) {
	for key, value := range row {
		switch key {
		case "id":
			a.Id = param.AsUint64(value)
		case "customer_id":
			a.CustomerId = param.AsUint64(value)
		case "name":
			a.Name = param.AsString(value)
		case "token":
			a.Token = param.AsString(value)
		case "type":
			a.Type = param.AsString(value)
		case "extra":
			a.Extra = param.AsString(value)
		case "step":
			a.Step = param.AsUint(value)
		case "precondition":
			a.Precondition = param.AsString(value)
		case "created":
			a.Created = param.AsUint(value)
		}
	}
}

func (a *OfficialCustomerU2f) Set(key interface{}, value ...interface{}) {
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
		case "CustomerId":
			a.CustomerId = param.AsUint64(vv)
		case "Name":
			a.Name = param.AsString(vv)
		case "Token":
			a.Token = param.AsString(vv)
		case "Type":
			a.Type = param.AsString(vv)
		case "Extra":
			a.Extra = param.AsString(vv)
		case "Step":
			a.Step = param.AsUint(vv)
		case "Precondition":
			a.Precondition = param.AsString(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		}
	}
}

func (a *OfficialCustomerU2f) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["customer_id"] = a.CustomerId
		r["name"] = a.Name
		r["token"] = a.Token
		r["type"] = a.Type
		r["extra"] = a.Extra
		r["step"] = a.Step
		r["precondition"] = a.Precondition
		r["created"] = a.Created
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "customer_id":
			r["customer_id"] = a.CustomerId
		case "name":
			r["name"] = a.Name
		case "token":
			r["token"] = a.Token
		case "type":
			r["type"] = a.Type
		case "extra":
			r["extra"] = a.Extra
		case "step":
			r["step"] = a.Step
		case "precondition":
			r["precondition"] = a.Precondition
		case "created":
			r["created"] = a.Created
		}
	}
	return r
}

func (a *OfficialCustomerU2f) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCustomerU2f) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, recv, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCustomerU2f) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialCustomerU2f) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
