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

type Slice_OfficialCustomerWallet []*OfficialCustomerWallet

func (s Slice_OfficialCustomerWallet) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerWallet) RangeRaw(fn func(m *OfficialCustomerWallet) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerWallet) GroupBy(keyField string) map[string][]*OfficialCustomerWallet {
	r := map[string][]*OfficialCustomerWallet{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialCustomerWallet{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialCustomerWallet) KeyBy(keyField string) map[string]*OfficialCustomerWallet {
	r := map[string]*OfficialCustomerWallet{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialCustomerWallet) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialCustomerWallet) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialCustomerWallet) FromList(data interface{}) Slice_OfficialCustomerWallet {
	values, ok := data.([]*OfficialCustomerWallet)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialCustomerWallet{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialCustomerWallet(ctx echo.Context) *OfficialCustomerWallet {
	m := &OfficialCustomerWallet{}
	m.SetContext(ctx)
	return m
}

// OfficialCustomerWallet 钱包
type OfficialCustomerWallet struct {
	base    factory.Base
	objects []*OfficialCustomerWallet

	CustomerId  uint64  `db:"customer_id,pk" bson:"customer_id" comment:"客户ID" json:"customer_id" xml:"customer_id"`
	AssetType   string  `db:"asset_type,pk" bson:"asset_type" comment:"资产类型(money-钱;point-点数;credit-信用分;integral-积分;gold-金币;silver-银币;copper-铜币;experience-经验)" json:"asset_type" xml:"asset_type"`
	Balance     float64 `db:"balance" bson:"balance" comment:"余额" json:"balance" xml:"balance"`
	Freeze      float64 `db:"freeze" bson:"freeze" comment:"冻结金额" json:"freeze" xml:"freeze"`
	Accumulated float64 `db:"accumulated" bson:"accumulated" comment:"累计总金额" json:"accumulated" xml:"accumulated"`
	Created     uint    `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Updated     uint    `db:"updated" bson:"updated" comment:"更新时间" json:"updated" xml:"updated"`
}

// - base function

func (a *OfficialCustomerWallet) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialCustomerWallet) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialCustomerWallet) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialCustomerWallet) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialCustomerWallet) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialCustomerWallet) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialCustomerWallet) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialCustomerWallet) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialCustomerWallet) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialCustomerWallet) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialCustomerWallet) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialCustomerWallet) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialCustomerWallet) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialCustomerWallet) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialCustomerWallet) Objects() []*OfficialCustomerWallet {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialCustomerWallet) XObjects() Slice_OfficialCustomerWallet {
	return Slice_OfficialCustomerWallet(a.Objects())
}

func (a *OfficialCustomerWallet) NewObjects() factory.Ranger {
	return &Slice_OfficialCustomerWallet{}
}

func (a *OfficialCustomerWallet) InitObjects() *[]*OfficialCustomerWallet {
	a.objects = []*OfficialCustomerWallet{}
	return &a.objects
}

func (a *OfficialCustomerWallet) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialCustomerWallet) Short_() string {
	return "official_customer_wallet"
}

func (a *OfficialCustomerWallet) Struct_() string {
	return "OfficialCustomerWallet"
}

func (a *OfficialCustomerWallet) Name_() string {
	b := a
	if b == nil {
		b = &OfficialCustomerWallet{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialCustomerWallet) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialCustomerWallet) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialCustomerWallet) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerWallet:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerWallet(*v))
		case []*OfficialCustomerWallet:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerWallet(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerWallet) GroupBy(keyField string, inputRows ...[]*OfficialCustomerWallet) map[string][]*OfficialCustomerWallet {
	var rows Slice_OfficialCustomerWallet
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerWallet(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerWallet(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialCustomerWallet) KeyBy(keyField string, inputRows ...[]*OfficialCustomerWallet) map[string]*OfficialCustomerWallet {
	var rows Slice_OfficialCustomerWallet
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerWallet(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerWallet(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialCustomerWallet) AsKV(keyField string, valueField string, inputRows ...[]*OfficialCustomerWallet) param.Store {
	var rows Slice_OfficialCustomerWallet
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerWallet(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerWallet(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialCustomerWallet) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerWallet:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerWallet(*v))
		case []*OfficialCustomerWallet:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerWallet(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerWallet) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	if len(a.AssetType) == 0 {
		a.AssetType = "money"
	}
	if a.base.Eventable() {
		err = DBI.Fire("creating", a, nil)
		if err != nil {
			return
		}
	}
	pk, err = a.Param(nil).SetSend(a).Insert()

	if err == nil && a.base.Eventable() {
		err = DBI.Fire("created", a, nil)
	}
	return
}

func (a *OfficialCustomerWallet) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.AssetType) == 0 {
		a.AssetType = "money"
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

func (a *OfficialCustomerWallet) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.AssetType) == 0 {
		a.AssetType = "money"
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

func (a *OfficialCustomerWallet) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.AssetType) == 0 {
		a.AssetType = "money"
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

func (a *OfficialCustomerWallet) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.AssetType) == 0 {
		a.AssetType = "money"
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

func (a *OfficialCustomerWallet) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerWallet) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerWallet) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["asset_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["asset_type"] = "money"
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

func (a *OfficialCustomerWallet) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["asset_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["asset_type"] = "money"
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

func (a *OfficialCustomerWallet) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialCustomerWallet) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		a.Updated = uint(time.Now().Unix())
		if len(a.AssetType) == 0 {
			a.AssetType = "money"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		if len(a.AssetType) == 0 {
			a.AssetType = "money"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("creating", a, nil)
	})

	if err == nil && a.base.Eventable() {
		if pk == nil {
			err = DBI.Fire("updated", a, mw, args...)
		} else {
			err = DBI.Fire("created", a, nil)
		}
	}
	return
}

func (a *OfficialCustomerWallet) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialCustomerWallet) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialCustomerWallet) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialCustomerWallet) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialCustomerWallet) Reset() *OfficialCustomerWallet {
	a.CustomerId = 0
	a.AssetType = ``
	a.Balance = 0.0
	a.Freeze = 0.0
	a.Accumulated = 0.0
	a.Created = 0
	a.Updated = 0
	return a
}

func (a *OfficialCustomerWallet) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["CustomerId"] = a.CustomerId
		r["AssetType"] = a.AssetType
		r["Balance"] = a.Balance
		r["Freeze"] = a.Freeze
		r["Accumulated"] = a.Accumulated
		r["Created"] = a.Created
		r["Updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "CustomerId":
			r["CustomerId"] = a.CustomerId
		case "AssetType":
			r["AssetType"] = a.AssetType
		case "Balance":
			r["Balance"] = a.Balance
		case "Freeze":
			r["Freeze"] = a.Freeze
		case "Accumulated":
			r["Accumulated"] = a.Accumulated
		case "Created":
			r["Created"] = a.Created
		case "Updated":
			r["Updated"] = a.Updated
		}
	}
	return r
}

func (a *OfficialCustomerWallet) FromRow(row map[string]interface{}) {
	for key, value := range row {
		switch key {
		case "customer_id":
			a.CustomerId = param.AsUint64(value)
		case "asset_type":
			a.AssetType = param.AsString(value)
		case "balance":
			a.Balance = param.AsFloat64(value)
		case "freeze":
			a.Freeze = param.AsFloat64(value)
		case "accumulated":
			a.Accumulated = param.AsFloat64(value)
		case "created":
			a.Created = param.AsUint(value)
		case "updated":
			a.Updated = param.AsUint(value)
		}
	}
}

func (a *OfficialCustomerWallet) Set(key interface{}, value ...interface{}) {
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
		case "CustomerId":
			a.CustomerId = param.AsUint64(vv)
		case "AssetType":
			a.AssetType = param.AsString(vv)
		case "Balance":
			a.Balance = param.AsFloat64(vv)
		case "Freeze":
			a.Freeze = param.AsFloat64(vv)
		case "Accumulated":
			a.Accumulated = param.AsFloat64(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "Updated":
			a.Updated = param.AsUint(vv)
		}
	}
}

func (a *OfficialCustomerWallet) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["customer_id"] = a.CustomerId
		r["asset_type"] = a.AssetType
		r["balance"] = a.Balance
		r["freeze"] = a.Freeze
		r["accumulated"] = a.Accumulated
		r["created"] = a.Created
		r["updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "customer_id":
			r["customer_id"] = a.CustomerId
		case "asset_type":
			r["asset_type"] = a.AssetType
		case "balance":
			r["balance"] = a.Balance
		case "freeze":
			r["freeze"] = a.Freeze
		case "accumulated":
			r["accumulated"] = a.Accumulated
		case "created":
			r["created"] = a.Created
		case "updated":
			r["updated"] = a.Updated
		}
	}
	return r
}

func (a *OfficialCustomerWallet) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCustomerWallet) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, recv, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCustomerWallet) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialCustomerWallet) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
