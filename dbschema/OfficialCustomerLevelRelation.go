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

type Slice_OfficialCustomerLevelRelation []*OfficialCustomerLevelRelation

func (s Slice_OfficialCustomerLevelRelation) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerLevelRelation) RangeRaw(fn func(m *OfficialCustomerLevelRelation) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerLevelRelation) GroupBy(keyField string) map[string][]*OfficialCustomerLevelRelation {
	r := map[string][]*OfficialCustomerLevelRelation{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialCustomerLevelRelation{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialCustomerLevelRelation) KeyBy(keyField string) map[string]*OfficialCustomerLevelRelation {
	r := map[string]*OfficialCustomerLevelRelation{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialCustomerLevelRelation) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialCustomerLevelRelation) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialCustomerLevelRelation) FromList(data interface{}) Slice_OfficialCustomerLevelRelation {
	values, ok := data.([]*OfficialCustomerLevelRelation)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialCustomerLevelRelation{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialCustomerLevelRelation(ctx echo.Context) *OfficialCustomerLevelRelation {
	m := &OfficialCustomerLevelRelation{}
	m.SetContext(ctx)
	return m
}

// OfficialCustomerLevelRelation 客户等级关联
type OfficialCustomerLevelRelation struct {
	base    factory.Base
	objects []*OfficialCustomerLevelRelation

	Id              uint64 `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	CustomerId      uint64 `db:"customer_id" bson:"customer_id" comment:"客户ID" json:"customer_id" xml:"customer_id"`
	LevelId         uint   `db:"level_id" bson:"level_id" comment:"等级ID" json:"level_id" xml:"level_id"`
	Status          string `db:"status" bson:"status" comment:"状态(actived-有效;expired-已过期)" json:"status" xml:"status"`
	Expired         uint   `db:"expired" bson:"expired" comment:"过期时间(0为永不过期)" json:"expired" xml:"expired"`
	AccumulatedDays uint   `db:"accumulated_days" bson:"accumulated_days" comment:"累计天数" json:"accumulated_days" xml:"accumulated_days"`
	LastRenewalAt   uint   `db:"last_renewal_at" bson:"last_renewal_at" comment:"最近续费时间" json:"last_renewal_at" xml:"last_renewal_at"`
	Created         uint   `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Updated         uint   `db:"updated" bson:"updated" comment:"更新时间" json:"updated" xml:"updated"`
}

// - base function

func (a *OfficialCustomerLevelRelation) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialCustomerLevelRelation) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialCustomerLevelRelation) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialCustomerLevelRelation) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialCustomerLevelRelation) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialCustomerLevelRelation) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialCustomerLevelRelation) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialCustomerLevelRelation) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialCustomerLevelRelation) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialCustomerLevelRelation) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialCustomerLevelRelation) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialCustomerLevelRelation) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialCustomerLevelRelation) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialCustomerLevelRelation) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialCustomerLevelRelation) Objects() []*OfficialCustomerLevelRelation {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialCustomerLevelRelation) XObjects() Slice_OfficialCustomerLevelRelation {
	return Slice_OfficialCustomerLevelRelation(a.Objects())
}

func (a *OfficialCustomerLevelRelation) NewObjects() factory.Ranger {
	return &Slice_OfficialCustomerLevelRelation{}
}

func (a *OfficialCustomerLevelRelation) InitObjects() *[]*OfficialCustomerLevelRelation {
	a.objects = []*OfficialCustomerLevelRelation{}
	return &a.objects
}

func (a *OfficialCustomerLevelRelation) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialCustomerLevelRelation) Short_() string {
	return "official_customer_level_relation"
}

func (a *OfficialCustomerLevelRelation) Struct_() string {
	return "OfficialCustomerLevelRelation"
}

func (a *OfficialCustomerLevelRelation) Name_() string {
	b := a
	if b == nil {
		b = &OfficialCustomerLevelRelation{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialCustomerLevelRelation) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialCustomerLevelRelation) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialCustomerLevelRelation) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerLevelRelation:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerLevelRelation(*v))
		case []*OfficialCustomerLevelRelation:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerLevelRelation(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerLevelRelation) GroupBy(keyField string, inputRows ...[]*OfficialCustomerLevelRelation) map[string][]*OfficialCustomerLevelRelation {
	var rows Slice_OfficialCustomerLevelRelation
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerLevelRelation(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerLevelRelation(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialCustomerLevelRelation) KeyBy(keyField string, inputRows ...[]*OfficialCustomerLevelRelation) map[string]*OfficialCustomerLevelRelation {
	var rows Slice_OfficialCustomerLevelRelation
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerLevelRelation(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerLevelRelation(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialCustomerLevelRelation) AsKV(keyField string, valueField string, inputRows ...[]*OfficialCustomerLevelRelation) param.Store {
	var rows Slice_OfficialCustomerLevelRelation
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerLevelRelation(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerLevelRelation(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialCustomerLevelRelation) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerLevelRelation:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerLevelRelation(*v))
		case []*OfficialCustomerLevelRelation:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerLevelRelation(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerLevelRelation) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.Status) == 0 {
		a.Status = "actived"
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

func (a *OfficialCustomerLevelRelation) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Status) == 0 {
		a.Status = "actived"
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

func (a *OfficialCustomerLevelRelation) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Status) == 0 {
		a.Status = "actived"
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

func (a *OfficialCustomerLevelRelation) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Status) == 0 {
		a.Status = "actived"
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

func (a *OfficialCustomerLevelRelation) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.Status) == 0 {
		a.Status = "actived"
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

func (a *OfficialCustomerLevelRelation) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerLevelRelation) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerLevelRelation) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["status"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["status"] = "actived"
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

func (a *OfficialCustomerLevelRelation) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["status"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["status"] = "actived"
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

func (a *OfficialCustomerLevelRelation) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialCustomerLevelRelation) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		a.Updated = uint(time.Now().Unix())
		if len(a.Status) == 0 {
			a.Status = "actived"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		a.Id = 0
		if len(a.Status) == 0 {
			a.Status = "actived"
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

func (a *OfficialCustomerLevelRelation) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialCustomerLevelRelation) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialCustomerLevelRelation) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialCustomerLevelRelation) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialCustomerLevelRelation) Reset() *OfficialCustomerLevelRelation {
	a.Id = 0
	a.CustomerId = 0
	a.LevelId = 0
	a.Status = ``
	a.Expired = 0
	a.AccumulatedDays = 0
	a.LastRenewalAt = 0
	a.Created = 0
	a.Updated = 0
	return a
}

func (a *OfficialCustomerLevelRelation) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["CustomerId"] = a.CustomerId
		r["LevelId"] = a.LevelId
		r["Status"] = a.Status
		r["Expired"] = a.Expired
		r["AccumulatedDays"] = a.AccumulatedDays
		r["LastRenewalAt"] = a.LastRenewalAt
		r["Created"] = a.Created
		r["Updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "CustomerId":
			r["CustomerId"] = a.CustomerId
		case "LevelId":
			r["LevelId"] = a.LevelId
		case "Status":
			r["Status"] = a.Status
		case "Expired":
			r["Expired"] = a.Expired
		case "AccumulatedDays":
			r["AccumulatedDays"] = a.AccumulatedDays
		case "LastRenewalAt":
			r["LastRenewalAt"] = a.LastRenewalAt
		case "Created":
			r["Created"] = a.Created
		case "Updated":
			r["Updated"] = a.Updated
		}
	}
	return r
}

func (a *OfficialCustomerLevelRelation) FromRow(row map[string]interface{}) {
	for key, value := range row {
		if _, ok := value.(db.RawValue); ok {
			continue
		}
		switch key {
		case "id":
			a.Id = param.AsUint64(value)
		case "customer_id":
			a.CustomerId = param.AsUint64(value)
		case "level_id":
			a.LevelId = param.AsUint(value)
		case "status":
			a.Status = param.AsString(value)
		case "expired":
			a.Expired = param.AsUint(value)
		case "accumulated_days":
			a.AccumulatedDays = param.AsUint(value)
		case "last_renewal_at":
			a.LastRenewalAt = param.AsUint(value)
		case "created":
			a.Created = param.AsUint(value)
		case "updated":
			a.Updated = param.AsUint(value)
		}
	}
}

func (a *OfficialCustomerLevelRelation) GetField(field string) interface{} {
	switch field {
	case "Id":
		return a.Id
	case "CustomerId":
		return a.CustomerId
	case "LevelId":
		return a.LevelId
	case "Status":
		return a.Status
	case "Expired":
		return a.Expired
	case "AccumulatedDays":
		return a.AccumulatedDays
	case "LastRenewalAt":
		return a.LastRenewalAt
	case "Created":
		return a.Created
	case "Updated":
		return a.Updated
	default:
		return nil
	}
}

func (a *OfficialCustomerLevelRelation) GetAllFieldNames() []string {
	return []string{
		"Id",
		"CustomerId",
		"LevelId",
		"Status",
		"Expired",
		"AccumulatedDays",
		"LastRenewalAt",
		"Created",
		"Updated",
	}
}

func (a *OfficialCustomerLevelRelation) HasField(field string) bool {
	switch field {
	case "Id":
		return true
	case "CustomerId":
		return true
	case "LevelId":
		return true
	case "Status":
		return true
	case "Expired":
		return true
	case "AccumulatedDays":
		return true
	case "LastRenewalAt":
		return true
	case "Created":
		return true
	case "Updated":
		return true
	default:
		return false
	}
}

func (a *OfficialCustomerLevelRelation) Set(key interface{}, value ...interface{}) {
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
		case "LevelId":
			a.LevelId = param.AsUint(vv)
		case "Status":
			a.Status = param.AsString(vv)
		case "Expired":
			a.Expired = param.AsUint(vv)
		case "AccumulatedDays":
			a.AccumulatedDays = param.AsUint(vv)
		case "LastRenewalAt":
			a.LastRenewalAt = param.AsUint(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "Updated":
			a.Updated = param.AsUint(vv)
		}
	}
}

func (a *OfficialCustomerLevelRelation) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["customer_id"] = a.CustomerId
		r["level_id"] = a.LevelId
		r["status"] = a.Status
		r["expired"] = a.Expired
		r["accumulated_days"] = a.AccumulatedDays
		r["last_renewal_at"] = a.LastRenewalAt
		r["created"] = a.Created
		r["updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "customer_id":
			r["customer_id"] = a.CustomerId
		case "level_id":
			r["level_id"] = a.LevelId
		case "status":
			r["status"] = a.Status
		case "expired":
			r["expired"] = a.Expired
		case "accumulated_days":
			r["accumulated_days"] = a.AccumulatedDays
		case "last_renewal_at":
			r["last_renewal_at"] = a.LastRenewalAt
		case "created":
			r["created"] = a.Created
		case "updated":
			r["updated"] = a.Updated
		}
	}
	return r
}

func (a *OfficialCustomerLevelRelation) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPage(a, cond, sorts...)
}

func (a *OfficialCustomerLevelRelation) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageAs(a, recv, cond, sorts...)
}

func (a *OfficialCustomerLevelRelation) ListPageByOffset(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffset(a, cond, sorts...)
}

func (a *OfficialCustomerLevelRelation) ListPageByOffsetAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffsetAs(a, recv, cond, sorts...)
}

func (a *OfficialCustomerLevelRelation) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialCustomerLevelRelation) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
