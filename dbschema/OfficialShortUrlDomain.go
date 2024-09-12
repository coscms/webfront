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

type Slice_OfficialShortUrlDomain []*OfficialShortUrlDomain

func (s Slice_OfficialShortUrlDomain) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialShortUrlDomain) RangeRaw(fn func(m *OfficialShortUrlDomain) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialShortUrlDomain) GroupBy(keyField string) map[string][]*OfficialShortUrlDomain {
	r := map[string][]*OfficialShortUrlDomain{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialShortUrlDomain{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialShortUrlDomain) KeyBy(keyField string) map[string]*OfficialShortUrlDomain {
	r := map[string]*OfficialShortUrlDomain{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialShortUrlDomain) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialShortUrlDomain) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialShortUrlDomain) FromList(data interface{}) Slice_OfficialShortUrlDomain {
	values, ok := data.([]*OfficialShortUrlDomain)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialShortUrlDomain{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialShortUrlDomain(ctx echo.Context) *OfficialShortUrlDomain {
	m := &OfficialShortUrlDomain{}
	m.SetContext(ctx)
	return m
}

// OfficialShortUrlDomain 短网址域名
type OfficialShortUrlDomain struct {
	base    factory.Base
	objects []*OfficialShortUrlDomain

	Id        uint64 `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	OwnerId   uint64 `db:"owner_id" bson:"owner_id" comment:"所有者客户ID" json:"owner_id" xml:"owner_id"`
	OwnerType string `db:"owner_type" bson:"owner_type" comment:"所有者类型(customer-前台客户;user-后台用户)" json:"owner_type" xml:"owner_type"`
	Domain    string `db:"domain" bson:"domain" comment:"域名" json:"domain" xml:"domain"`
	UrlCount  uint64 `db:"url_count" bson:"url_count" comment:"网址统计" json:"url_count" xml:"url_count"`
	Disabled  string `db:"disabled" bson:"disabled" comment:"是否(Y/N)禁用" json:"disabled" xml:"disabled"`
	Created   uint   `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Updated   uint   `db:"updated" bson:"updated" comment:"更新时间" json:"updated" xml:"updated"`
}

// - base function

func (a *OfficialShortUrlDomain) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialShortUrlDomain) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialShortUrlDomain) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialShortUrlDomain) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialShortUrlDomain) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialShortUrlDomain) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialShortUrlDomain) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialShortUrlDomain) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialShortUrlDomain) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialShortUrlDomain) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialShortUrlDomain) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialShortUrlDomain) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialShortUrlDomain) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialShortUrlDomain) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialShortUrlDomain) Objects() []*OfficialShortUrlDomain {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialShortUrlDomain) XObjects() Slice_OfficialShortUrlDomain {
	return Slice_OfficialShortUrlDomain(a.Objects())
}

func (a *OfficialShortUrlDomain) NewObjects() factory.Ranger {
	return &Slice_OfficialShortUrlDomain{}
}

func (a *OfficialShortUrlDomain) InitObjects() *[]*OfficialShortUrlDomain {
	a.objects = []*OfficialShortUrlDomain{}
	return &a.objects
}

func (a *OfficialShortUrlDomain) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialShortUrlDomain) Short_() string {
	return "official_short_url_domain"
}

func (a *OfficialShortUrlDomain) Struct_() string {
	return "OfficialShortUrlDomain"
}

func (a *OfficialShortUrlDomain) Name_() string {
	b := a
	if b == nil {
		b = &OfficialShortUrlDomain{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialShortUrlDomain) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialShortUrlDomain) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialShortUrlDomain) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialShortUrlDomain:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialShortUrlDomain(*v))
		case []*OfficialShortUrlDomain:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialShortUrlDomain(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialShortUrlDomain) GroupBy(keyField string, inputRows ...[]*OfficialShortUrlDomain) map[string][]*OfficialShortUrlDomain {
	var rows Slice_OfficialShortUrlDomain
	if len(inputRows) > 0 {
		rows = Slice_OfficialShortUrlDomain(inputRows[0])
	} else {
		rows = Slice_OfficialShortUrlDomain(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialShortUrlDomain) KeyBy(keyField string, inputRows ...[]*OfficialShortUrlDomain) map[string]*OfficialShortUrlDomain {
	var rows Slice_OfficialShortUrlDomain
	if len(inputRows) > 0 {
		rows = Slice_OfficialShortUrlDomain(inputRows[0])
	} else {
		rows = Slice_OfficialShortUrlDomain(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialShortUrlDomain) AsKV(keyField string, valueField string, inputRows ...[]*OfficialShortUrlDomain) param.Store {
	var rows Slice_OfficialShortUrlDomain
	if len(inputRows) > 0 {
		rows = Slice_OfficialShortUrlDomain(inputRows[0])
	} else {
		rows = Slice_OfficialShortUrlDomain(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialShortUrlDomain) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialShortUrlDomain:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialShortUrlDomain(*v))
		case []*OfficialShortUrlDomain:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialShortUrlDomain(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialShortUrlDomain) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.OwnerType) == 0 {
		a.OwnerType = "customer"
	}
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

func (a *OfficialShortUrlDomain) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.OwnerType) == 0 {
		a.OwnerType = "customer"
	}
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

func (a *OfficialShortUrlDomain) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.OwnerType) == 0 {
		a.OwnerType = "customer"
	}
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

func (a *OfficialShortUrlDomain) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.OwnerType) == 0 {
		a.OwnerType = "customer"
	}
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

func (a *OfficialShortUrlDomain) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.OwnerType) == 0 {
		a.OwnerType = "customer"
	}
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

func (a *OfficialShortUrlDomain) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialShortUrlDomain) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialShortUrlDomain) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["owner_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["owner_type"] = "customer"
		}
	}
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

func (a *OfficialShortUrlDomain) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["owner_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["owner_type"] = "customer"
		}
	}
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

func (a *OfficialShortUrlDomain) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialShortUrlDomain) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		a.Updated = uint(time.Now().Unix())
		if len(a.OwnerType) == 0 {
			a.OwnerType = "customer"
		}
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
		if len(a.OwnerType) == 0 {
			a.OwnerType = "customer"
		}
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

func (a *OfficialShortUrlDomain) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialShortUrlDomain) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialShortUrlDomain) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialShortUrlDomain) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialShortUrlDomain) Reset() *OfficialShortUrlDomain {
	a.Id = 0
	a.OwnerId = 0
	a.OwnerType = ``
	a.Domain = ``
	a.UrlCount = 0
	a.Disabled = ``
	a.Created = 0
	a.Updated = 0
	return a
}

func (a *OfficialShortUrlDomain) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["OwnerId"] = a.OwnerId
		r["OwnerType"] = a.OwnerType
		r["Domain"] = a.Domain
		r["UrlCount"] = a.UrlCount
		r["Disabled"] = a.Disabled
		r["Created"] = a.Created
		r["Updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "OwnerId":
			r["OwnerId"] = a.OwnerId
		case "OwnerType":
			r["OwnerType"] = a.OwnerType
		case "Domain":
			r["Domain"] = a.Domain
		case "UrlCount":
			r["UrlCount"] = a.UrlCount
		case "Disabled":
			r["Disabled"] = a.Disabled
		case "Created":
			r["Created"] = a.Created
		case "Updated":
			r["Updated"] = a.Updated
		}
	}
	return r
}

func (a *OfficialShortUrlDomain) FromRow(row map[string]interface{}) {
	for key, value := range row {
		switch key {
		case "id":
			a.Id = param.AsUint64(value)
		case "owner_id":
			a.OwnerId = param.AsUint64(value)
		case "owner_type":
			a.OwnerType = param.AsString(value)
		case "domain":
			a.Domain = param.AsString(value)
		case "url_count":
			a.UrlCount = param.AsUint64(value)
		case "disabled":
			a.Disabled = param.AsString(value)
		case "created":
			a.Created = param.AsUint(value)
		case "updated":
			a.Updated = param.AsUint(value)
		}
	}
}

func (a *OfficialShortUrlDomain) Set(key interface{}, value ...interface{}) {
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
		case "OwnerId":
			a.OwnerId = param.AsUint64(vv)
		case "OwnerType":
			a.OwnerType = param.AsString(vv)
		case "Domain":
			a.Domain = param.AsString(vv)
		case "UrlCount":
			a.UrlCount = param.AsUint64(vv)
		case "Disabled":
			a.Disabled = param.AsString(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "Updated":
			a.Updated = param.AsUint(vv)
		}
	}
}

func (a *OfficialShortUrlDomain) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["owner_id"] = a.OwnerId
		r["owner_type"] = a.OwnerType
		r["domain"] = a.Domain
		r["url_count"] = a.UrlCount
		r["disabled"] = a.Disabled
		r["created"] = a.Created
		r["updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "owner_id":
			r["owner_id"] = a.OwnerId
		case "owner_type":
			r["owner_type"] = a.OwnerType
		case "domain":
			r["domain"] = a.Domain
		case "url_count":
			r["url_count"] = a.UrlCount
		case "disabled":
			r["disabled"] = a.Disabled
		case "created":
			r["created"] = a.Created
		case "updated":
			r["updated"] = a.Updated
		}
	}
	return r
}

func (a *OfficialShortUrlDomain) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialShortUrlDomain) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, recv, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialShortUrlDomain) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialShortUrlDomain) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
