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

type Slice_OfficialCommonMessage []*OfficialCommonMessage

func (s Slice_OfficialCommonMessage) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCommonMessage) RangeRaw(fn func(m *OfficialCommonMessage) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCommonMessage) GroupBy(keyField string) map[string][]*OfficialCommonMessage {
	r := map[string][]*OfficialCommonMessage{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialCommonMessage{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialCommonMessage) KeyBy(keyField string) map[string]*OfficialCommonMessage {
	r := map[string]*OfficialCommonMessage{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialCommonMessage) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialCommonMessage) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialCommonMessage) FromList(data interface{}) Slice_OfficialCommonMessage {
	values, ok := data.([]*OfficialCommonMessage)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialCommonMessage{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialCommonMessage(ctx echo.Context) *OfficialCommonMessage {
	m := &OfficialCommonMessage{}
	m.SetContext(ctx)
	return m
}

// OfficialCommonMessage 站内信
type OfficialCommonMessage struct {
	base    factory.Base
	objects []*OfficialCommonMessage

	Id              uint64 `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	Type            string `db:"type" bson:"type" comment:"消息类型" json:"type" xml:"type"`
	CustomerA       uint64 `db:"customer_a" bson:"customer_a" comment:"发信人ID(0为系统消息)" json:"customer_a" xml:"customer_a"`
	CustomerB       uint64 `db:"customer_b" bson:"customer_b" comment:"收信人ID" json:"customer_b" xml:"customer_b"`
	CustomerGroupId uint   `db:"customer_group_id" bson:"customer_group_id" comment:"客户组消息" json:"customer_group_id" xml:"customer_group_id"`
	UserA           uint   `db:"user_a" bson:"user_a" comment:"发信人ID(后台用户ID，用于系统消息)" json:"user_a" xml:"user_a"`
	UserB           uint   `db:"user_b" bson:"user_b" comment:"收信人ID(后台用户ID，用于后台消息)" json:"user_b" xml:"user_b"`
	UserRoleId      uint   `db:"user_role_id" bson:"user_role_id" comment:"后台角色消息" json:"user_role_id" xml:"user_role_id"`
	Title           string `db:"title" bson:"title" comment:"消息标题" json:"title" xml:"title"`
	Content         string `db:"content" bson:"content" comment:"消息内容" json:"content" xml:"content"`
	Contype         string `db:"contype" bson:"contype" comment:"内容类型" json:"contype" xml:"contype"`
	Encrypted       string `db:"encrypted" bson:"encrypted" comment:"是否为加密消息" json:"encrypted" xml:"encrypted"`
	Password        string `db:"password" bson:"password" comment:"密码" json:"password" xml:"password"`
	Created         uint   `db:"created" bson:"created" comment:"发送时间" json:"created" xml:"created"`
	Url             string `db:"url" bson:"url" comment:"网址" json:"url" xml:"url"`
	RootId          uint64 `db:"root_id" bson:"root_id" comment:"根ID" json:"root_id" xml:"root_id"`
	ReplyId         uint64 `db:"reply_id" bson:"reply_id" comment:"回复ID" json:"reply_id" xml:"reply_id"`
	HasNewReply     uint   `db:"has_new_reply" bson:"has_new_reply" comment:"是否(1/0)有新回复" json:"has_new_reply" xml:"has_new_reply"`
	ViewProgress    uint   `db:"view_progress" bson:"view_progress" comment:"查看总进度(100为100%)" json:"view_progress" xml:"view_progress"`
}

// - base function

func (a *OfficialCommonMessage) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialCommonMessage) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialCommonMessage) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialCommonMessage) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialCommonMessage) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialCommonMessage) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialCommonMessage) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialCommonMessage) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialCommonMessage) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialCommonMessage) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialCommonMessage) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialCommonMessage) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialCommonMessage) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialCommonMessage) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialCommonMessage) Objects() []*OfficialCommonMessage {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialCommonMessage) XObjects() Slice_OfficialCommonMessage {
	return Slice_OfficialCommonMessage(a.Objects())
}

func (a *OfficialCommonMessage) NewObjects() factory.Ranger {
	return &Slice_OfficialCommonMessage{}
}

func (a *OfficialCommonMessage) InitObjects() *[]*OfficialCommonMessage {
	a.objects = []*OfficialCommonMessage{}
	return &a.objects
}

func (a *OfficialCommonMessage) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialCommonMessage) Short_() string {
	return "official_common_message"
}

func (a *OfficialCommonMessage) Struct_() string {
	return "OfficialCommonMessage"
}

func (a *OfficialCommonMessage) Name_() string {
	b := a
	if b == nil {
		b = &OfficialCommonMessage{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialCommonMessage) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialCommonMessage) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialCommonMessage) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCommonMessage:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonMessage(*v))
		case []*OfficialCommonMessage:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonMessage(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCommonMessage) GroupBy(keyField string, inputRows ...[]*OfficialCommonMessage) map[string][]*OfficialCommonMessage {
	var rows Slice_OfficialCommonMessage
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonMessage(inputRows[0])
	} else {
		rows = Slice_OfficialCommonMessage(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialCommonMessage) KeyBy(keyField string, inputRows ...[]*OfficialCommonMessage) map[string]*OfficialCommonMessage {
	var rows Slice_OfficialCommonMessage
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonMessage(inputRows[0])
	} else {
		rows = Slice_OfficialCommonMessage(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialCommonMessage) AsKV(keyField string, valueField string, inputRows ...[]*OfficialCommonMessage) param.Store {
	var rows Slice_OfficialCommonMessage
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonMessage(inputRows[0])
	} else {
		rows = Slice_OfficialCommonMessage(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialCommonMessage) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCommonMessage:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonMessage(*v))
		case []*OfficialCommonMessage:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonMessage(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCommonMessage) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.Contype) == 0 {
		a.Contype = "text"
	}
	if len(a.Encrypted) == 0 {
		a.Encrypted = "N"
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

func (a *OfficialCommonMessage) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {

	if len(a.Contype) == 0 {
		a.Contype = "text"
	}
	if len(a.Encrypted) == 0 {
		a.Encrypted = "N"
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

func (a *OfficialCommonMessage) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

	if len(a.Contype) == 0 {
		a.Contype = "text"
	}
	if len(a.Encrypted) == 0 {
		a.Encrypted = "N"
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

func (a *OfficialCommonMessage) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {

	if len(a.Contype) == 0 {
		a.Contype = "text"
	}
	if len(a.Encrypted) == 0 {
		a.Encrypted = "N"
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

func (a *OfficialCommonMessage) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {

	if len(a.Contype) == 0 {
		a.Contype = "text"
	}
	if len(a.Encrypted) == 0 {
		a.Encrypted = "N"
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

func (a *OfficialCommonMessage) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCommonMessage) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCommonMessage) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["contype"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["contype"] = "text"
		}
	}
	if val, ok := kvset["encrypted"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["encrypted"] = "N"
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

func (a *OfficialCommonMessage) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["contype"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["contype"] = "text"
		}
	}
	if val, ok := kvset["encrypted"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["encrypted"] = "N"
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

func (a *OfficialCommonMessage) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialCommonMessage) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		if len(a.Contype) == 0 {
			a.Contype = "text"
		}
		if len(a.Encrypted) == 0 {
			a.Encrypted = "N"
		}
		if !a.base.Eventable() {
			return nil
		}
		return DBI.Fire("updating", a, mw, args...)
	}, func() error {
		a.Created = uint(time.Now().Unix())
		a.Id = 0
		if len(a.Contype) == 0 {
			a.Contype = "text"
		}
		if len(a.Encrypted) == 0 {
			a.Encrypted = "N"
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

func (a *OfficialCommonMessage) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialCommonMessage) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialCommonMessage) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialCommonMessage) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialCommonMessage) Reset() *OfficialCommonMessage {
	a.Id = 0
	a.Type = ``
	a.CustomerA = 0
	a.CustomerB = 0
	a.CustomerGroupId = 0
	a.UserA = 0
	a.UserB = 0
	a.UserRoleId = 0
	a.Title = ``
	a.Content = ``
	a.Contype = ``
	a.Encrypted = ``
	a.Password = ``
	a.Created = 0
	a.Url = ``
	a.RootId = 0
	a.ReplyId = 0
	a.HasNewReply = 0
	a.ViewProgress = 0
	return a
}

func (a *OfficialCommonMessage) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["Type"] = a.Type
		r["CustomerA"] = a.CustomerA
		r["CustomerB"] = a.CustomerB
		r["CustomerGroupId"] = a.CustomerGroupId
		r["UserA"] = a.UserA
		r["UserB"] = a.UserB
		r["UserRoleId"] = a.UserRoleId
		r["Title"] = a.Title
		r["Content"] = a.Content
		r["Contype"] = a.Contype
		r["Encrypted"] = a.Encrypted
		r["Password"] = a.Password
		r["Created"] = a.Created
		r["Url"] = a.Url
		r["RootId"] = a.RootId
		r["ReplyId"] = a.ReplyId
		r["HasNewReply"] = a.HasNewReply
		r["ViewProgress"] = a.ViewProgress
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "Type":
			r["Type"] = a.Type
		case "CustomerA":
			r["CustomerA"] = a.CustomerA
		case "CustomerB":
			r["CustomerB"] = a.CustomerB
		case "CustomerGroupId":
			r["CustomerGroupId"] = a.CustomerGroupId
		case "UserA":
			r["UserA"] = a.UserA
		case "UserB":
			r["UserB"] = a.UserB
		case "UserRoleId":
			r["UserRoleId"] = a.UserRoleId
		case "Title":
			r["Title"] = a.Title
		case "Content":
			r["Content"] = a.Content
		case "Contype":
			r["Contype"] = a.Contype
		case "Encrypted":
			r["Encrypted"] = a.Encrypted
		case "Password":
			r["Password"] = a.Password
		case "Created":
			r["Created"] = a.Created
		case "Url":
			r["Url"] = a.Url
		case "RootId":
			r["RootId"] = a.RootId
		case "ReplyId":
			r["ReplyId"] = a.ReplyId
		case "HasNewReply":
			r["HasNewReply"] = a.HasNewReply
		case "ViewProgress":
			r["ViewProgress"] = a.ViewProgress
		}
	}
	return r
}

func (a *OfficialCommonMessage) FromRow(row map[string]interface{}) {
	for key, value := range row {
		switch key {
		case "id":
			a.Id = param.AsUint64(value)
		case "type":
			a.Type = param.AsString(value)
		case "customer_a":
			a.CustomerA = param.AsUint64(value)
		case "customer_b":
			a.CustomerB = param.AsUint64(value)
		case "customer_group_id":
			a.CustomerGroupId = param.AsUint(value)
		case "user_a":
			a.UserA = param.AsUint(value)
		case "user_b":
			a.UserB = param.AsUint(value)
		case "user_role_id":
			a.UserRoleId = param.AsUint(value)
		case "title":
			a.Title = param.AsString(value)
		case "content":
			a.Content = param.AsString(value)
		case "contype":
			a.Contype = param.AsString(value)
		case "encrypted":
			a.Encrypted = param.AsString(value)
		case "password":
			a.Password = param.AsString(value)
		case "created":
			a.Created = param.AsUint(value)
		case "url":
			a.Url = param.AsString(value)
		case "root_id":
			a.RootId = param.AsUint64(value)
		case "reply_id":
			a.ReplyId = param.AsUint64(value)
		case "has_new_reply":
			a.HasNewReply = param.AsUint(value)
		case "view_progress":
			a.ViewProgress = param.AsUint(value)
		}
	}
}

func (a *OfficialCommonMessage) Set(key interface{}, value ...interface{}) {
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
		case "Type":
			a.Type = param.AsString(vv)
		case "CustomerA":
			a.CustomerA = param.AsUint64(vv)
		case "CustomerB":
			a.CustomerB = param.AsUint64(vv)
		case "CustomerGroupId":
			a.CustomerGroupId = param.AsUint(vv)
		case "UserA":
			a.UserA = param.AsUint(vv)
		case "UserB":
			a.UserB = param.AsUint(vv)
		case "UserRoleId":
			a.UserRoleId = param.AsUint(vv)
		case "Title":
			a.Title = param.AsString(vv)
		case "Content":
			a.Content = param.AsString(vv)
		case "Contype":
			a.Contype = param.AsString(vv)
		case "Encrypted":
			a.Encrypted = param.AsString(vv)
		case "Password":
			a.Password = param.AsString(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "Url":
			a.Url = param.AsString(vv)
		case "RootId":
			a.RootId = param.AsUint64(vv)
		case "ReplyId":
			a.ReplyId = param.AsUint64(vv)
		case "HasNewReply":
			a.HasNewReply = param.AsUint(vv)
		case "ViewProgress":
			a.ViewProgress = param.AsUint(vv)
		}
	}
}

func (a *OfficialCommonMessage) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["type"] = a.Type
		r["customer_a"] = a.CustomerA
		r["customer_b"] = a.CustomerB
		r["customer_group_id"] = a.CustomerGroupId
		r["user_a"] = a.UserA
		r["user_b"] = a.UserB
		r["user_role_id"] = a.UserRoleId
		r["title"] = a.Title
		r["content"] = a.Content
		r["contype"] = a.Contype
		r["encrypted"] = a.Encrypted
		r["password"] = a.Password
		r["created"] = a.Created
		r["url"] = a.Url
		r["root_id"] = a.RootId
		r["reply_id"] = a.ReplyId
		r["has_new_reply"] = a.HasNewReply
		r["view_progress"] = a.ViewProgress
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "type":
			r["type"] = a.Type
		case "customer_a":
			r["customer_a"] = a.CustomerA
		case "customer_b":
			r["customer_b"] = a.CustomerB
		case "customer_group_id":
			r["customer_group_id"] = a.CustomerGroupId
		case "user_a":
			r["user_a"] = a.UserA
		case "user_b":
			r["user_b"] = a.UserB
		case "user_role_id":
			r["user_role_id"] = a.UserRoleId
		case "title":
			r["title"] = a.Title
		case "content":
			r["content"] = a.Content
		case "contype":
			r["contype"] = a.Contype
		case "encrypted":
			r["encrypted"] = a.Encrypted
		case "password":
			r["password"] = a.Password
		case "created":
			r["created"] = a.Created
		case "url":
			r["url"] = a.Url
		case "root_id":
			r["root_id"] = a.RootId
		case "reply_id":
			r["reply_id"] = a.ReplyId
		case "has_new_reply":
			r["has_new_reply"] = a.HasNewReply
		case "view_progress":
			r["view_progress"] = a.ViewProgress
		}
	}
	return r
}

func (a *OfficialCommonMessage) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCommonMessage) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, recv, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCommonMessage) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialCommonMessage) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
