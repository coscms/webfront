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

type Slice_OfficialCustomerWalletFlow []*OfficialCustomerWalletFlow

func (s Slice_OfficialCustomerWalletFlow) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerWalletFlow) RangeRaw(fn func(m *OfficialCustomerWalletFlow) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCustomerWalletFlow) GroupBy(keyField string) map[string][]*OfficialCustomerWalletFlow {
	r := map[string][]*OfficialCustomerWalletFlow{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialCustomerWalletFlow{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialCustomerWalletFlow) KeyBy(keyField string) map[string]*OfficialCustomerWalletFlow {
	r := map[string]*OfficialCustomerWalletFlow{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialCustomerWalletFlow) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialCustomerWalletFlow) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialCustomerWalletFlow) FromList(data interface{}) Slice_OfficialCustomerWalletFlow {
	values, ok := data.([]*OfficialCustomerWalletFlow)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialCustomerWalletFlow{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialCustomerWalletFlow(ctx echo.Context) *OfficialCustomerWalletFlow {
	m := &OfficialCustomerWalletFlow{}
	m.SetContext(ctx)
	return m
}

// OfficialCustomerWalletFlow 钱包流水记录
type OfficialCustomerWalletFlow struct {
	base    factory.Base
	objects []*OfficialCustomerWalletFlow

	Id             uint64  `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	CustomerId     uint64  `db:"customer_id" bson:"customer_id" comment:"客户ID" json:"customer_id" xml:"customer_id"`
	AssetType      string  `db:"asset_type" bson:"asset_type" comment:"资产类型" json:"asset_type" xml:"asset_type"`
	AmountType     string  `db:"amount_type" bson:"amount_type" comment:"金额类型(balance-余额;freeze-冻结额)" json:"amount_type" xml:"amount_type"`
	Amount         float64 `db:"amount" bson:"amount" comment:"金额(正数为收入;负数为支出)" json:"amount" xml:"amount"`
	WalletAmount   float64 `db:"wallet_amount" bson:"wallet_amount" comment:"变动后钱包总金额" json:"wallet_amount" xml:"wallet_amount"`
	SourceCustomer uint64  `db:"source_customer" bson:"source_customer" comment:"来自谁" json:"source_customer" xml:"source_customer"`
	SourceType     string  `db:"source_type" bson:"source_type" comment:"来源类型(组)" json:"source_type" xml:"source_type"`
	SourceTable    string  `db:"source_table" bson:"source_table" comment:"来源表(来自物品表)" json:"source_table" xml:"source_table"`
	SourceId       uint64  `db:"source_id" bson:"source_id" comment:"来源ID(来自物品ID)" json:"source_id" xml:"source_id"`
	Number         uint64  `db:"number" bson:"number" comment:"备用编号" json:"number" xml:"number"`
	TradeNo        string  `db:"trade_no" bson:"trade_no" comment:"交易号(来自哪个交易)" json:"trade_no" xml:"trade_no"`
	Status         string  `db:"status" bson:"status" comment:"状态(pending-待确认;confirmed-已确认;canceled-已取消)" json:"status" xml:"status"`
	Description    string  `db:"description" bson:"description" comment:"简短描述" json:"description" xml:"description"`
	Created        uint    `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
}

// - base function

func (a *OfficialCustomerWalletFlow) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialCustomerWalletFlow) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialCustomerWalletFlow) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialCustomerWalletFlow) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialCustomerWalletFlow) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialCustomerWalletFlow) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialCustomerWalletFlow) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialCustomerWalletFlow) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialCustomerWalletFlow) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialCustomerWalletFlow) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialCustomerWalletFlow) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialCustomerWalletFlow) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialCustomerWalletFlow) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialCustomerWalletFlow) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialCustomerWalletFlow) Objects() []*OfficialCustomerWalletFlow {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialCustomerWalletFlow) XObjects() Slice_OfficialCustomerWalletFlow {
	return Slice_OfficialCustomerWalletFlow(a.Objects())
}

func (a *OfficialCustomerWalletFlow) NewObjects() factory.Ranger {
	return &Slice_OfficialCustomerWalletFlow{}
}

func (a *OfficialCustomerWalletFlow) InitObjects() *[]*OfficialCustomerWalletFlow {
	a.objects = []*OfficialCustomerWalletFlow{}
	return &a.objects
}

func (a *OfficialCustomerWalletFlow) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialCustomerWalletFlow) Short_() string {
	return "official_customer_wallet_flow"
}

func (a *OfficialCustomerWalletFlow) Struct_() string {
	return "OfficialCustomerWalletFlow"
}

func (a *OfficialCustomerWalletFlow) Name_() string {
	b := a
	if b == nil {
		b = &OfficialCustomerWalletFlow{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialCustomerWalletFlow) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialCustomerWalletFlow) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialCustomerWalletFlow) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerWalletFlow:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerWalletFlow(*v))
		case []*OfficialCustomerWalletFlow:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerWalletFlow(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerWalletFlow) GroupBy(keyField string, inputRows ...[]*OfficialCustomerWalletFlow) map[string][]*OfficialCustomerWalletFlow {
	var rows Slice_OfficialCustomerWalletFlow
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerWalletFlow(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerWalletFlow(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialCustomerWalletFlow) KeyBy(keyField string, inputRows ...[]*OfficialCustomerWalletFlow) map[string]*OfficialCustomerWalletFlow {
	var rows Slice_OfficialCustomerWalletFlow
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerWalletFlow(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerWalletFlow(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialCustomerWalletFlow) AsKV(keyField string, valueField string, inputRows ...[]*OfficialCustomerWalletFlow) param.Store {
	var rows Slice_OfficialCustomerWalletFlow
	if len(inputRows) > 0 {
		rows = Slice_OfficialCustomerWalletFlow(inputRows[0])
	} else {
		rows = Slice_OfficialCustomerWalletFlow(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialCustomerWalletFlow) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCustomerWalletFlow:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerWalletFlow(*v))
		case []*OfficialCustomerWalletFlow:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCustomerWalletFlow(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCustomerWalletFlow) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.AmountType) == 0 {
		a.AmountType = "balance"
	}
	if len(a.Status) == 0 {
		a.Status = "confirmed"
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

func (a *OfficialCustomerWalletFlow) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {

	if len(a.AmountType) == 0 {
		a.AmountType = "balance"
	}
	if len(a.Status) == 0 {
		a.Status = "confirmed"
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

func (a *OfficialCustomerWalletFlow) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

	if len(a.AmountType) == 0 {
		a.AmountType = "balance"
	}
	if len(a.Status) == 0 {
		a.Status = "confirmed"
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

func (a *OfficialCustomerWalletFlow) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {

	if len(a.AmountType) == 0 {
		a.AmountType = "balance"
	}
	if len(a.Status) == 0 {
		a.Status = "confirmed"
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

func (a *OfficialCustomerWalletFlow) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {

	if len(a.AmountType) == 0 {
		a.AmountType = "balance"
	}
	if len(a.Status) == 0 {
		a.Status = "confirmed"
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

func (a *OfficialCustomerWalletFlow) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerWalletFlow) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCustomerWalletFlow) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["amount_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["amount_type"] = "balance"
		}
	}
	if val, ok := kvset["status"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["status"] = "confirmed"
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

func (a *OfficialCustomerWalletFlow) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["amount_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["amount_type"] = "balance"
		}
	}
	if val, ok := kvset["status"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["status"] = "confirmed"
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

func (a *OfficialCustomerWalletFlow) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialCustomerWalletFlow) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		if len(a.AmountType) == 0 {
			a.AmountType = "balance"
		}
		if len(a.Status) == 0 {
			a.Status = "confirmed"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		a.Id = 0
		if len(a.AmountType) == 0 {
			a.AmountType = "balance"
		}
		if len(a.Status) == 0 {
			a.Status = "confirmed"
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

func (a *OfficialCustomerWalletFlow) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialCustomerWalletFlow) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialCustomerWalletFlow) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialCustomerWalletFlow) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialCustomerWalletFlow) Reset() *OfficialCustomerWalletFlow {
	a.Id = 0
	a.CustomerId = 0
	a.AssetType = ``
	a.AmountType = ``
	a.Amount = 0.0
	a.WalletAmount = 0.0
	a.SourceCustomer = 0
	a.SourceType = ``
	a.SourceTable = ``
	a.SourceId = 0
	a.Number = 0
	a.TradeNo = ``
	a.Status = ``
	a.Description = ``
	a.Created = 0
	return a
}

func (a *OfficialCustomerWalletFlow) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["CustomerId"] = a.CustomerId
		r["AssetType"] = a.AssetType
		r["AmountType"] = a.AmountType
		r["Amount"] = a.Amount
		r["WalletAmount"] = a.WalletAmount
		r["SourceCustomer"] = a.SourceCustomer
		r["SourceType"] = a.SourceType
		r["SourceTable"] = a.SourceTable
		r["SourceId"] = a.SourceId
		r["Number"] = a.Number
		r["TradeNo"] = a.TradeNo
		r["Status"] = a.Status
		r["Description"] = a.Description
		r["Created"] = a.Created
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "CustomerId":
			r["CustomerId"] = a.CustomerId
		case "AssetType":
			r["AssetType"] = a.AssetType
		case "AmountType":
			r["AmountType"] = a.AmountType
		case "Amount":
			r["Amount"] = a.Amount
		case "WalletAmount":
			r["WalletAmount"] = a.WalletAmount
		case "SourceCustomer":
			r["SourceCustomer"] = a.SourceCustomer
		case "SourceType":
			r["SourceType"] = a.SourceType
		case "SourceTable":
			r["SourceTable"] = a.SourceTable
		case "SourceId":
			r["SourceId"] = a.SourceId
		case "Number":
			r["Number"] = a.Number
		case "TradeNo":
			r["TradeNo"] = a.TradeNo
		case "Status":
			r["Status"] = a.Status
		case "Description":
			r["Description"] = a.Description
		case "Created":
			r["Created"] = a.Created
		}
	}
	return r
}

func (a *OfficialCustomerWalletFlow) FromRow(row map[string]interface{}) {
	for key, value := range row {
		switch key {
		case "id":
			a.Id = param.AsUint64(value)
		case "customer_id":
			a.CustomerId = param.AsUint64(value)
		case "asset_type":
			a.AssetType = param.AsString(value)
		case "amount_type":
			a.AmountType = param.AsString(value)
		case "amount":
			a.Amount = param.AsFloat64(value)
		case "wallet_amount":
			a.WalletAmount = param.AsFloat64(value)
		case "source_customer":
			a.SourceCustomer = param.AsUint64(value)
		case "source_type":
			a.SourceType = param.AsString(value)
		case "source_table":
			a.SourceTable = param.AsString(value)
		case "source_id":
			a.SourceId = param.AsUint64(value)
		case "number":
			a.Number = param.AsUint64(value)
		case "trade_no":
			a.TradeNo = param.AsString(value)
		case "status":
			a.Status = param.AsString(value)
		case "description":
			a.Description = param.AsString(value)
		case "created":
			a.Created = param.AsUint(value)
		}
	}
}

func (a *OfficialCustomerWalletFlow) Set(key interface{}, value ...interface{}) {
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
		case "AssetType":
			a.AssetType = param.AsString(vv)
		case "AmountType":
			a.AmountType = param.AsString(vv)
		case "Amount":
			a.Amount = param.AsFloat64(vv)
		case "WalletAmount":
			a.WalletAmount = param.AsFloat64(vv)
		case "SourceCustomer":
			a.SourceCustomer = param.AsUint64(vv)
		case "SourceType":
			a.SourceType = param.AsString(vv)
		case "SourceTable":
			a.SourceTable = param.AsString(vv)
		case "SourceId":
			a.SourceId = param.AsUint64(vv)
		case "Number":
			a.Number = param.AsUint64(vv)
		case "TradeNo":
			a.TradeNo = param.AsString(vv)
		case "Status":
			a.Status = param.AsString(vv)
		case "Description":
			a.Description = param.AsString(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		}
	}
}

func (a *OfficialCustomerWalletFlow) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["customer_id"] = a.CustomerId
		r["asset_type"] = a.AssetType
		r["amount_type"] = a.AmountType
		r["amount"] = a.Amount
		r["wallet_amount"] = a.WalletAmount
		r["source_customer"] = a.SourceCustomer
		r["source_type"] = a.SourceType
		r["source_table"] = a.SourceTable
		r["source_id"] = a.SourceId
		r["number"] = a.Number
		r["trade_no"] = a.TradeNo
		r["status"] = a.Status
		r["description"] = a.Description
		r["created"] = a.Created
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "customer_id":
			r["customer_id"] = a.CustomerId
		case "asset_type":
			r["asset_type"] = a.AssetType
		case "amount_type":
			r["amount_type"] = a.AmountType
		case "amount":
			r["amount"] = a.Amount
		case "wallet_amount":
			r["wallet_amount"] = a.WalletAmount
		case "source_customer":
			r["source_customer"] = a.SourceCustomer
		case "source_type":
			r["source_type"] = a.SourceType
		case "source_table":
			r["source_table"] = a.SourceTable
		case "source_id":
			r["source_id"] = a.SourceId
		case "number":
			r["number"] = a.Number
		case "trade_no":
			r["trade_no"] = a.TradeNo
		case "status":
			r["status"] = a.Status
		case "description":
			r["description"] = a.Description
		case "created":
			r["created"] = a.Created
		}
	}
	return r
}

func (a *OfficialCustomerWalletFlow) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCustomerWalletFlow) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, recv, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCustomerWalletFlow) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialCustomerWalletFlow) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
