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

type Slice_OfficialCustomerRole []*OfficialCustomerRole

func (s Slice_OfficialCustomerRole) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerRole) RangeRaw(fn func(m *OfficialCustomerRole) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerRole) GroupBy(keyField string) map[string][]*OfficialCustomerRole {
	r := map[string][]*OfficialCustomerRole{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialCustomerRole{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialCustomerRole) KeyBy(keyField string) map[string]*OfficialCustomerRole {
	r := map[string]*OfficialCustomerRole{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialCustomerRole) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialCustomerRole) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialCustomerRole) FromList(data interface{}) Slice_OfficialCustomerRole {
	values, ok := data.([]*OfficialCustomerRole)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialCustomerRole{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialCustomerRole(ctx echo.Context) *OfficialCustomerRole {
	m := &OfficialCustomerRole{}
	m.SetContext(ctx)
	return m
}

// OfficialCustomerRole 客户角色
type OfficialCustomerRole struct {
	base    factory.Base
	objects []*OfficialCustomerRole

	Id          uint   `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	Name        string `db:"name" bson:"name" comment:"名称" json:"name" xml:"name"`
	Description string `db:"description" bson:"description" comment:"说明" json:"description" xml:"description"`
	Created     uint   `db:"created" bson:"created" comment:"添加时间" json:"created" xml:"created"`
	Updated     uint   `db:"updated" bson:"updated" comment:"修改时间" json:"updated" xml:"updated"`
	Disabled    string `db:"disabled" bson:"disabled" comment:"是否禁用" json:"disabled" xml:"disabled"`
	IsDefault   string `db:"is_default" bson:"is_default" comment:"是否作为新用户注册时的默认角色" json:"is_default" xml:"is_default"`
	ParentId    uint   `db:"parent_id" bson:"parent_id" comment:"父级ID" json:"parent_id" xml:"parent_id"`
}

// - base function

func (a *OfficialCustomerRole) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialCustomerRole) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialCustomerRole) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialCustomerRole) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialCustomerRole) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialCustomerRole) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialCustomerRole) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialCustomerRole) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialCustomerRole) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialCustomerRole) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialCustomerRole) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialCustomerRole) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialCustomerRole) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialCustomerRole) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialCustomerRole) Objects() []*OfficialCustomerRole {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialCustomerRole) XObjects() Slice_OfficialCustomerRole {
	return Slice_OfficialCustomerRole(a.Objects())
}

func (a *OfficialCustomerRole) NewObjects() factory.Ranger {
	return &Slice_OfficialCustomerRole{}
}

func (a *OfficialCustomerRole) InitObjects() *[]*OfficialCustomerRole {
	a.objects = []*OfficialCustomerRole{}
	return &a.objects
}

func (a *OfficialCustomerRole) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialCustomerRole) Short_() string {
	return "official_customer_role"
}

func (a *OfficialCustomerRole) Struct_() string {
	return "OfficialCustomerRole"
}

func (a *OfficialCustomerRole) Name_() string {
	b := a
	if b == nil {
		b = &OfficialCustomerRole{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialCustomerRole) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialCustomerRole) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialCustomerRole) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerRole:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerRole(*v))
		case []*OfficialCustomerRole:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerRole(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerRole) GroupBy(keyField string, inputRows ...[]*OfficialCustomerRole) map[string][]*OfficialCustomerRole {
	var rows Slice_OfficialCustomerRole
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerRole(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerRole(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialCustomerRole) KeyBy(keyField string, inputRows ...[]*OfficialCustomerRole) map[string]*OfficialCustomerRole {
	var rows Slice_OfficialCustomerRole
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerRole(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerRole(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialCustomerRole) AsKV(keyField string, valueField string, inputRows ...[]*OfficialCustomerRole) param.Store {
	var rows Slice_OfficialCustomerRole
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerRole(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerRole(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialCustomerRole) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerRole:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerRole(*v))
		case []*OfficialCustomerRole:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerRole(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerRole) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.IsDefault) == 0 {
		a.IsDefault = "N"
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

func (a *OfficialCustomerRole) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.IsDefault) == 0 {
		a.IsDefault = "N"
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

func (a *OfficialCustomerRole) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.IsDefault) == 0 {
		a.IsDefault = "N"
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

func (a *OfficialCustomerRole) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.IsDefault) == 0 {
		a.IsDefault = "N"
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

func (a *OfficialCustomerRole) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Disabled) == 0 {
		a.Disabled = "N"
	}
	if len(a.IsDefault) == 0 {
		a.IsDefault = "N"
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

func (a *OfficialCustomerRole) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerRole) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerRole) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["disabled"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["disabled"] = "N"
		}
	}
	if val, ok := kvset["is_default"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["is_default"] = "N"
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

func (a *OfficialCustomerRole) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["disabled"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["disabled"] = "N"
		}
	}
	if val, ok := kvset["is_default"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["is_default"] = "N"
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

func (a *OfficialCustomerRole) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialCustomerRole) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		a.Updated = uint(time.Now().Unix())
		if len(a.Disabled) == 0 {
			a.Disabled = "N"
		}
		if len(a.IsDefault) == 0 {
			a.IsDefault = "N"
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
		if len(a.IsDefault) == 0 {
			a.IsDefault = "N"
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

func (a *OfficialCustomerRole) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialCustomerRole) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialCustomerRole) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialCustomerRole) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialCustomerRole) Reset() *OfficialCustomerRole {
	a.Id = 0
	a.Name = ``
	a.Description = ``
	a.Created = 0
	a.Updated = 0
	a.Disabled = ``
	a.IsDefault = ``
	a.ParentId = 0
	return a
}

func (a *OfficialCustomerRole) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["Name"] = a.Name
		r["Description"] = a.Description
		r["Created"] = a.Created
		r["Updated"] = a.Updated
		r["Disabled"] = a.Disabled
		r["IsDefault"] = a.IsDefault
		r["ParentId"] = a.ParentId
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "Name":
			r["Name"] = a.Name
		case "Description":
			r["Description"] = a.Description
		case "Created":
			r["Created"] = a.Created
		case "Updated":
			r["Updated"] = a.Updated
		case "Disabled":
			r["Disabled"] = a.Disabled
		case "IsDefault":
			r["IsDefault"] = a.IsDefault
		case "ParentId":
			r["ParentId"] = a.ParentId
		}
	}
	return r
}

func (a *OfficialCustomerRole) FromRow(row map[string]interface{}) {
	for key, value := range row {
		if _, ok := value.(db.RawValue); ok {
			continue
		}
		switch key {
		case "id":
			a.Id = param.AsUint(value)
		case "name":
			a.Name = param.AsString(value)
		case "description":
			a.Description = param.AsString(value)
		case "created":
			a.Created = param.AsUint(value)
		case "updated":
			a.Updated = param.AsUint(value)
		case "disabled":
			a.Disabled = param.AsString(value)
		case "is_default":
			a.IsDefault = param.AsString(value)
		case "parent_id":
			a.ParentId = param.AsUint(value)
		}
	}
}

func (a *OfficialCustomerRole) GetField(field string) interface{} {
	switch field {
	case "Id":
		return a.Id
	case "Name":
		return a.Name
	case "Description":
		return a.Description
	case "Created":
		return a.Created
	case "Updated":
		return a.Updated
	case "Disabled":
		return a.Disabled
	case "IsDefault":
		return a.IsDefault
	case "ParentId":
		return a.ParentId
	default:
		return nil
	}
}

func (a *OfficialCustomerRole) GetAllFieldNames() []string {
	return []string{
		"Id",
		"Name",
		"Description",
		"Created",
		"Updated",
		"Disabled",
		"IsDefault",
		"ParentId",
	}
}

func (a *OfficialCustomerRole) HasField(field string) bool {
	switch field {
	case "Id":
		return true
	case "Name":
		return true
	case "Description":
		return true
	case "Created":
		return true
	case "Updated":
		return true
	case "Disabled":
		return true
	case "IsDefault":
		return true
	case "ParentId":
		return true
	default:
		return false
	}
}

func (a *OfficialCustomerRole) Set(key interface{}, value ...interface{}) {
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
		case "Description":
			a.Description = param.AsString(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "Updated":
			a.Updated = param.AsUint(vv)
		case "Disabled":
			a.Disabled = param.AsString(vv)
		case "IsDefault":
			a.IsDefault = param.AsString(vv)
		case "ParentId":
			a.ParentId = param.AsUint(vv)
		}
	}
}

func (a *OfficialCustomerRole) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["name"] = a.Name
		r["description"] = a.Description
		r["created"] = a.Created
		r["updated"] = a.Updated
		r["disabled"] = a.Disabled
		r["is_default"] = a.IsDefault
		r["parent_id"] = a.ParentId
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "name":
			r["name"] = a.Name
		case "description":
			r["description"] = a.Description
		case "created":
			r["created"] = a.Created
		case "updated":
			r["updated"] = a.Updated
		case "disabled":
			r["disabled"] = a.Disabled
		case "is_default":
			r["is_default"] = a.IsDefault
		case "parent_id":
			r["parent_id"] = a.ParentId
		}
	}
	return r
}

func (a *OfficialCustomerRole) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPage(a, cond, sorts...)
}

func (a *OfficialCustomerRole) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageAs(a, recv, cond, sorts...)
}

func (a *OfficialCustomerRole) ListPageByOffset(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffset(a, cond, sorts...)
}

func (a *OfficialCustomerRole) ListPageByOffsetAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffsetAs(a, recv, cond, sorts...)
}

func (a *OfficialCustomerRole) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialCustomerRole) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
