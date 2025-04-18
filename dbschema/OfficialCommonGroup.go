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

type Slice_OfficialCommonGroup []*OfficialCommonGroup

func (s Slice_OfficialCommonGroup) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCommonGroup) RangeRaw(fn func(m *OfficialCommonGroup) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCommonGroup) GroupBy(keyField string) map[string][]*OfficialCommonGroup {
	r := map[string][]*OfficialCommonGroup{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialCommonGroup{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialCommonGroup) KeyBy(keyField string) map[string]*OfficialCommonGroup {
	r := map[string]*OfficialCommonGroup{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialCommonGroup) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialCommonGroup) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialCommonGroup) FromList(data interface{}) Slice_OfficialCommonGroup {
	values, ok := data.([]*OfficialCommonGroup)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialCommonGroup{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialCommonGroup(ctx echo.Context) *OfficialCommonGroup {
	m := &OfficialCommonGroup{}
	m.SetContext(ctx)
	return m
}

// OfficialCommonGroup 分组
type OfficialCommonGroup struct {
	base    factory.Base
	objects []*OfficialCommonGroup

	Id          uint   `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	ParentId    uint   `db:"parent_id" bson:"parent_id" comment:"上级ID" json:"parent_id" xml:"parent_id"`
	Uid         uint   `db:"uid" bson:"uid" comment:"用户ID" json:"uid" xml:"uid"`
	Name        string `db:"name" bson:"name" comment:"组名" json:"name" xml:"name"`
	Type        string `db:"type" bson:"type" comment:"类型(customer-客户组;cert-证书组;order-订单组;product-产品组;attr-产品属性组;openapp-开放平台应用;api-外部接口组)" json:"type" xml:"type"`
	Description string `db:"description" bson:"description" comment:"说明" json:"description" xml:"description"`
	Created     uint   `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
}

// - base function

func (a *OfficialCommonGroup) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialCommonGroup) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialCommonGroup) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialCommonGroup) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialCommonGroup) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialCommonGroup) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialCommonGroup) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialCommonGroup) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialCommonGroup) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialCommonGroup) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialCommonGroup) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialCommonGroup) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialCommonGroup) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialCommonGroup) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialCommonGroup) Objects() []*OfficialCommonGroup {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialCommonGroup) XObjects() Slice_OfficialCommonGroup {
	return Slice_OfficialCommonGroup(a.Objects())
}

func (a *OfficialCommonGroup) NewObjects() factory.Ranger {
	return &Slice_OfficialCommonGroup{}
}

func (a *OfficialCommonGroup) InitObjects() *[]*OfficialCommonGroup {
	a.objects = []*OfficialCommonGroup{}
	return &a.objects
}

func (a *OfficialCommonGroup) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialCommonGroup) Short_() string {
	return "official_common_group"
}

func (a *OfficialCommonGroup) Struct_() string {
	return "OfficialCommonGroup"
}

func (a *OfficialCommonGroup) Name_() string {
	b := a
	if b == nil {
		b = &OfficialCommonGroup{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialCommonGroup) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialCommonGroup) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialCommonGroup) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCommonGroup:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonGroup(*v))
		case []*OfficialCommonGroup:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonGroup(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCommonGroup) GroupBy(keyField string, inputRows ...[]*OfficialCommonGroup) map[string][]*OfficialCommonGroup {
	var rows Slice_OfficialCommonGroup
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonGroup(inputRows[0])
	} else {
		rows = Slice_OfficialCommonGroup(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialCommonGroup) KeyBy(keyField string, inputRows ...[]*OfficialCommonGroup) map[string]*OfficialCommonGroup {
	var rows Slice_OfficialCommonGroup
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonGroup(inputRows[0])
	} else {
		rows = Slice_OfficialCommonGroup(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialCommonGroup) AsKV(keyField string, valueField string, inputRows ...[]*OfficialCommonGroup) param.Store {
	var rows Slice_OfficialCommonGroup
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonGroup(inputRows[0])
	} else {
		rows = Slice_OfficialCommonGroup(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialCommonGroup) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCommonGroup:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonGroup(*v))
		case []*OfficialCommonGroup:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonGroup(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCommonGroup) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.Type) == 0 {
		a.Type = "customer"
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

func (a *OfficialCommonGroup) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {

	if len(a.Type) == 0 {
		a.Type = "customer"
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

func (a *OfficialCommonGroup) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

	if len(a.Type) == 0 {
		a.Type = "customer"
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

func (a *OfficialCommonGroup) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {

	if len(a.Type) == 0 {
		a.Type = "customer"
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

func (a *OfficialCommonGroup) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {

	if len(a.Type) == 0 {
		a.Type = "customer"
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

func (a *OfficialCommonGroup) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCommonGroup) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCommonGroup) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["type"] = "customer"
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

func (a *OfficialCommonGroup) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["type"] = "customer"
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

func (a *OfficialCommonGroup) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialCommonGroup) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		if len(a.Type) == 0 {
			a.Type = "customer"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		a.Id = 0
		if len(a.Type) == 0 {
			a.Type = "customer"
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

func (a *OfficialCommonGroup) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialCommonGroup) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialCommonGroup) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialCommonGroup) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialCommonGroup) Reset() *OfficialCommonGroup {
	a.Id = 0
	a.ParentId = 0
	a.Uid = 0
	a.Name = ``
	a.Type = ``
	a.Description = ``
	a.Created = 0
	return a
}

func (a *OfficialCommonGroup) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["ParentId"] = a.ParentId
		r["Uid"] = a.Uid
		r["Name"] = a.Name
		r["Type"] = a.Type
		r["Description"] = a.Description
		r["Created"] = a.Created
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "ParentId":
			r["ParentId"] = a.ParentId
		case "Uid":
			r["Uid"] = a.Uid
		case "Name":
			r["Name"] = a.Name
		case "Type":
			r["Type"] = a.Type
		case "Description":
			r["Description"] = a.Description
		case "Created":
			r["Created"] = a.Created
		}
	}
	return r
}

func (a *OfficialCommonGroup) FromRow(row map[string]interface{}) {
	for key, value := range row {
		if _, ok := value.(db.RawValue); ok {
			continue
		}
		switch key {
		case "id":
			a.Id = param.AsUint(value)
		case "parent_id":
			a.ParentId = param.AsUint(value)
		case "uid":
			a.Uid = param.AsUint(value)
		case "name":
			a.Name = param.AsString(value)
		case "type":
			a.Type = param.AsString(value)
		case "description":
			a.Description = param.AsString(value)
		case "created":
			a.Created = param.AsUint(value)
		}
	}
}

func (a *OfficialCommonGroup) GetField(field string) interface{} {
	switch field {
	case "Id":
		return a.Id
	case "ParentId":
		return a.ParentId
	case "Uid":
		return a.Uid
	case "Name":
		return a.Name
	case "Type":
		return a.Type
	case "Description":
		return a.Description
	case "Created":
		return a.Created
	default:
		return nil
	}
}

func (a *OfficialCommonGroup) GetAllFieldNames() []string {
	return []string{
		"Id",
		"ParentId",
		"Uid",
		"Name",
		"Type",
		"Description",
		"Created",
	}
}

func (a *OfficialCommonGroup) HasField(field string) bool {
	switch field {
	case "Id":
		return true
	case "ParentId":
		return true
	case "Uid":
		return true
	case "Name":
		return true
	case "Type":
		return true
	case "Description":
		return true
	case "Created":
		return true
	default:
		return false
	}
}

func (a *OfficialCommonGroup) Set(key interface{}, value ...interface{}) {
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
		case "ParentId":
			a.ParentId = param.AsUint(vv)
		case "Uid":
			a.Uid = param.AsUint(vv)
		case "Name":
			a.Name = param.AsString(vv)
		case "Type":
			a.Type = param.AsString(vv)
		case "Description":
			a.Description = param.AsString(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		}
	}
}

func (a *OfficialCommonGroup) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["parent_id"] = a.ParentId
		r["uid"] = a.Uid
		r["name"] = a.Name
		r["type"] = a.Type
		r["description"] = a.Description
		r["created"] = a.Created
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "parent_id":
			r["parent_id"] = a.ParentId
		case "uid":
			r["uid"] = a.Uid
		case "name":
			r["name"] = a.Name
		case "type":
			r["type"] = a.Type
		case "description":
			r["description"] = a.Description
		case "created":
			r["created"] = a.Created
		}
	}
	return r
}

func (a *OfficialCommonGroup) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPage(a, cond, sorts...)
}

func (a *OfficialCommonGroup) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageAs(a, recv, cond, sorts...)
}

func (a *OfficialCommonGroup) ListPageByOffset(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffset(a, cond, sorts...)
}

func (a *OfficialCommonGroup) ListPageByOffsetAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffsetAs(a, recv, cond, sorts...)
}

func (a *OfficialCommonGroup) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialCommonGroup) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
