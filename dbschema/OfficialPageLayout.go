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

type Slice_OfficialPageLayout []*OfficialPageLayout

func (s Slice_OfficialPageLayout) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialPageLayout) RangeRaw(fn func(m *OfficialPageLayout) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialPageLayout) GroupBy(keyField string) map[string][]*OfficialPageLayout {
	r := map[string][]*OfficialPageLayout{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialPageLayout{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialPageLayout) KeyBy(keyField string) map[string]*OfficialPageLayout {
	r := map[string]*OfficialPageLayout{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialPageLayout) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialPageLayout) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialPageLayout) FromList(data interface{}) Slice_OfficialPageLayout {
	values, ok := data.([]*OfficialPageLayout)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialPageLayout{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialPageLayout(ctx echo.Context) *OfficialPageLayout {
	m := &OfficialPageLayout{}
	m.SetContext(ctx)
	return m
}

// OfficialPageLayout 页面布局所含区块
type OfficialPageLayout struct {
	base    factory.Base
	objects []*OfficialPageLayout

	Id       uint   `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	BlockId  uint   `db:"block_id" bson:"block_id" comment:"区块ID" json:"block_id" xml:"block_id"`
	PageId   uint   `db:"page_id" bson:"page_id" comment:"页面ID" json:"page_id" xml:"page_id"`
	Configs  string `db:"configs" bson:"configs" comment:"区块在布局中的配置" json:"configs" xml:"configs"`
	Sort     int    `db:"sort" bson:"sort" comment:"排序" json:"sort" xml:"sort"`
	Disabled string `db:"disabled" bson:"disabled" comment:"是否禁用" json:"disabled" xml:"disabled"`
	Created  uint   `db:"created" bson:"created" comment:"添加时间" json:"created" xml:"created"`
	Updated  uint   `db:"updated" bson:"updated" comment:"修改时间" json:"updated" xml:"updated"`
}

// - base function

func (a *OfficialPageLayout) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialPageLayout) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialPageLayout) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialPageLayout) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialPageLayout) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialPageLayout) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialPageLayout) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialPageLayout) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialPageLayout) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialPageLayout) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialPageLayout) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialPageLayout) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialPageLayout) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialPageLayout) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialPageLayout) Objects() []*OfficialPageLayout {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialPageLayout) XObjects() Slice_OfficialPageLayout {
	return Slice_OfficialPageLayout(a.Objects())
}

func (a *OfficialPageLayout) NewObjects() factory.Ranger {
	return &Slice_OfficialPageLayout{}
}

func (a *OfficialPageLayout) InitObjects() *[]*OfficialPageLayout {
	a.objects = []*OfficialPageLayout{}
	return &a.objects
}

func (a *OfficialPageLayout) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialPageLayout) Short_() string {
	return "official_page_layout"
}

func (a *OfficialPageLayout) Struct_() string {
	return "OfficialPageLayout"
}

func (a *OfficialPageLayout) Name_() string {
	b := a
	if b == nil {
		b = &OfficialPageLayout{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialPageLayout) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialPageLayout) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialPageLayout) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialPageLayout:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialPageLayout(*v))
		case []*OfficialPageLayout:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialPageLayout(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialPageLayout) GroupBy(keyField string, inputRows ...[]*OfficialPageLayout) map[string][]*OfficialPageLayout {
	var rows Slice_OfficialPageLayout
	if len(inputRows) > 0 {
		rows = Slice_OfficialPageLayout(inputRows[0])
	} else {
		rows = Slice_OfficialPageLayout(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialPageLayout) KeyBy(keyField string, inputRows ...[]*OfficialPageLayout) map[string]*OfficialPageLayout {
	var rows Slice_OfficialPageLayout
	if len(inputRows) > 0 {
		rows = Slice_OfficialPageLayout(inputRows[0])
	} else {
		rows = Slice_OfficialPageLayout(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialPageLayout) AsKV(keyField string, valueField string, inputRows ...[]*OfficialPageLayout) param.Store {
	var rows Slice_OfficialPageLayout
	if len(inputRows) > 0 {
		rows = Slice_OfficialPageLayout(inputRows[0])
	} else {
		rows = Slice_OfficialPageLayout(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialPageLayout) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialPageLayout:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialPageLayout(*v))
		case []*OfficialPageLayout:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialPageLayout(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialPageLayout) Insert() (pk interface{}, err error) {
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

func (a *OfficialPageLayout) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialPageLayout) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {
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

func (a *OfficialPageLayout) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {
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

func (a *OfficialPageLayout) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {
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

func (a *OfficialPageLayout) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialPageLayout) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialPageLayout) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

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

func (a *OfficialPageLayout) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialPageLayout) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialPageLayout) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
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

func (a *OfficialPageLayout) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialPageLayout) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialPageLayout) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialPageLayout) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialPageLayout) Reset() *OfficialPageLayout {
	a.Id = 0
	a.BlockId = 0
	a.PageId = 0
	a.Configs = ``
	a.Sort = 0
	a.Disabled = ``
	a.Created = 0
	a.Updated = 0
	return a
}

func (a *OfficialPageLayout) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["BlockId"] = a.BlockId
		r["PageId"] = a.PageId
		r["Configs"] = a.Configs
		r["Sort"] = a.Sort
		r["Disabled"] = a.Disabled
		r["Created"] = a.Created
		r["Updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "BlockId":
			r["BlockId"] = a.BlockId
		case "PageId":
			r["PageId"] = a.PageId
		case "Configs":
			r["Configs"] = a.Configs
		case "Sort":
			r["Sort"] = a.Sort
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

func (a *OfficialPageLayout) FromRow(row map[string]interface{}) {
	for key, value := range row {
		if _, ok := value.(db.RawValue); ok {
			continue
		}
		switch key {
		case "id":
			a.Id = param.AsUint(value)
		case "block_id":
			a.BlockId = param.AsUint(value)
		case "page_id":
			a.PageId = param.AsUint(value)
		case "configs":
			a.Configs = param.AsString(value)
		case "sort":
			a.Sort = param.AsInt(value)
		case "disabled":
			a.Disabled = param.AsString(value)
		case "created":
			a.Created = param.AsUint(value)
		case "updated":
			a.Updated = param.AsUint(value)
		}
	}
}

func (a *OfficialPageLayout) GetField(field string) interface{} {
	switch field {
	case "Id":
		return a.Id
	case "BlockId":
		return a.BlockId
	case "PageId":
		return a.PageId
	case "Configs":
		return a.Configs
	case "Sort":
		return a.Sort
	case "Disabled":
		return a.Disabled
	case "Created":
		return a.Created
	case "Updated":
		return a.Updated
	default:
		return nil
	}
}

func (a *OfficialPageLayout) GetAllFieldNames() []string {
	return []string{
		"Id",
		"BlockId",
		"PageId",
		"Configs",
		"Sort",
		"Disabled",
		"Created",
		"Updated",
	}
}

func (a *OfficialPageLayout) HasField(field string) bool {
	switch field {
	case "Id":
		return true
	case "BlockId":
		return true
	case "PageId":
		return true
	case "Configs":
		return true
	case "Sort":
		return true
	case "Disabled":
		return true
	case "Created":
		return true
	case "Updated":
		return true
	default:
		return false
	}
}

func (a *OfficialPageLayout) Set(key interface{}, value ...interface{}) {
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
		case "BlockId":
			a.BlockId = param.AsUint(vv)
		case "PageId":
			a.PageId = param.AsUint(vv)
		case "Configs":
			a.Configs = param.AsString(vv)
		case "Sort":
			a.Sort = param.AsInt(vv)
		case "Disabled":
			a.Disabled = param.AsString(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "Updated":
			a.Updated = param.AsUint(vv)
		}
	}
}

func (a *OfficialPageLayout) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["block_id"] = a.BlockId
		r["page_id"] = a.PageId
		r["configs"] = a.Configs
		r["sort"] = a.Sort
		r["disabled"] = a.Disabled
		r["created"] = a.Created
		r["updated"] = a.Updated
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "block_id":
			r["block_id"] = a.BlockId
		case "page_id":
			r["page_id"] = a.PageId
		case "configs":
			r["configs"] = a.Configs
		case "sort":
			r["sort"] = a.Sort
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

func (a *OfficialPageLayout) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPage(a, cond, sorts...)
}

func (a *OfficialPageLayout) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageAs(a, recv, cond, sorts...)
}

func (a *OfficialPageLayout) ListPageByOffset(cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffset(a, cond, sorts...)
}

func (a *OfficialPageLayout) ListPageByOffsetAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	return pagination.ListPageByOffsetAs(a, recv, cond, sorts...)
}

func (a *OfficialPageLayout) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialPageLayout) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
