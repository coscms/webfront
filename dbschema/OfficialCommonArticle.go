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

type Slice_OfficialCommonArticle []*OfficialCommonArticle

func (s Slice_OfficialCommonArticle) Range(fn func(m factory.Model) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCommonArticle) RangeRaw(fn func(m *OfficialCommonArticle) error) error {
	for _, v := range s {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (s Slice_OfficialCommonArticle) GroupBy(keyField string) map[string][]*OfficialCommonArticle {
	r := map[string][]*OfficialCommonArticle{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		if _, y := r[vkey]; !y {
			r[vkey] = []*OfficialCommonArticle{}
		}
		r[vkey] = append(r[vkey], row)
	}
	return r
}

func (s Slice_OfficialCommonArticle) KeyBy(keyField string) map[string]*OfficialCommonArticle {
	r := map[string]*OfficialCommonArticle{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = row
	}
	return r
}

func (s Slice_OfficialCommonArticle) AsKV(keyField string, valueField string) param.Store {
	r := param.Store{}
	for _, row := range s {
		dmap := row.AsMap()
		vkey := fmt.Sprint(dmap[keyField])
		r[vkey] = dmap[valueField]
	}
	return r
}

func (s Slice_OfficialCommonArticle) Transform(transfers map[string]param.Transfer) []param.Store {
	r := make([]param.Store, len(s))
	for idx, row := range s {
		r[idx] = row.AsMap().Transform(transfers)
	}
	return r
}

func (s Slice_OfficialCommonArticle) FromList(data interface{}) Slice_OfficialCommonArticle {
	values, ok := data.([]*OfficialCommonArticle)
	if !ok {
		for _, value := range data.([]interface{}) {
			row := &OfficialCommonArticle{}
			row.FromRow(value.(map[string]interface{}))
			s = append(s, row)
		}
		return s
	}
	s = append(s, values...)

	return s
}

func NewOfficialCommonArticle(ctx echo.Context) *OfficialCommonArticle {
	m := &OfficialCommonArticle{}
	m.SetContext(ctx)
	return m
}

// OfficialCommonArticle 官方新闻
type OfficialCommonArticle struct {
	base    factory.Base
	objects []*OfficialCommonArticle

	Id                 uint64  `db:"id,omitempty,pk" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	Category1          uint    `db:"category1" bson:"category1" comment:"顶级分类ID" json:"category1" xml:"category1"`
	Category2          uint    `db:"category2" bson:"category2" comment:"二级分类ID" json:"category2" xml:"category2"`
	Category3          uint    `db:"category3" bson:"category3" comment:"三级分类ID" json:"category3" xml:"category3"`
	CategoryId         uint    `db:"category_id" bson:"category_id" comment:"最底层分类ID" json:"category_id" xml:"category_id"`
	SourceId           string  `db:"source_id" bson:"source_id" comment:"来源ID(空代表不限)" json:"source_id" xml:"source_id"`
	SourceTable        string  `db:"source_table" bson:"source_table" comment:"来源表(不含official_前缀)" json:"source_table" xml:"source_table"`
	OwnerId            uint64  `db:"owner_id" bson:"owner_id" comment:"新闻发布者" json:"owner_id" xml:"owner_id"`
	OwnerType          string  `db:"owner_type" bson:"owner_type" comment:"所有者类型(customer-前台客户;user-后台用户)" json:"owner_type" xml:"owner_type"`
	Title              string  `db:"title" bson:"title" comment:"新闻标题" json:"title" xml:"title"`
	Keywords           string  `db:"keywords" bson:"keywords" comment:"关键词" json:"keywords" xml:"keywords"`
	Image              string  `db:"image" bson:"image" comment:"缩略图" json:"image" xml:"image"`
	ImageOriginal      string  `db:"image_original" bson:"image_original" comment:"原始图" json:"image_original" xml:"image_original"`
	Summary            string  `db:"summary" bson:"summary" comment:"摘要" json:"summary" xml:"summary"`
	Content            string  `db:"content" bson:"content" comment:"内容" json:"content" xml:"content"`
	Contype            string  `db:"contype" bson:"contype" comment:"内容类型" json:"contype" xml:"contype"`
	Created            uint    `db:"created" bson:"created" comment:"创建时间" json:"created" xml:"created"`
	Updated            uint    `db:"updated" bson:"updated" comment:"修改时间" json:"updated" xml:"updated"`
	Display            string  `db:"display" bson:"display" comment:"是否显示" json:"display" xml:"display"`
	Template           string  `db:"template" bson:"template" comment:"模版" json:"template" xml:"template"`
	Comments           uint64  `db:"comments" bson:"comments" comment:"评论数量" json:"comments" xml:"comments"`
	CloseComment       string  `db:"close_comment" bson:"close_comment" comment:"关闭评论" json:"close_comment" xml:"close_comment"`
	CommentAutoDisplay string  `db:"comment_auto_display" bson:"comment_auto_display" comment:"自动显示评论" json:"comment_auto_display" xml:"comment_auto_display"`
	CommentAllowUser   string  `db:"comment_allow_user" bson:"comment_allow_user" comment:"允许评论的用户(all-所有人;buyer-当前商品买家;author-当前文章作者;admin-管理员;allAgent-所有代理;curAgent-当前产品代理;none-无人;designated-指定人员)" json:"comment_allow_user" xml:"comment_allow_user"`
	Likes              uint64  `db:"likes" bson:"likes" comment:"好评数量" json:"likes" xml:"likes"`
	Hates              uint64  `db:"hates" bson:"hates" comment:"差评数量" json:"hates" xml:"hates"`
	Views              uint64  `db:"views" bson:"views" comment:"浏览次数" json:"views" xml:"views"`
	Tags               string  `db:"tags" bson:"tags" comment:"标签" json:"tags" xml:"tags"`
	Price              float64 `db:"price" bson:"price" comment:"价格" json:"price" xml:"price"`
	Slugify            string  `db:"slugify" bson:"slugify" comment:"SEO-friendly URLs with Slugify" json:"slugify" xml:"slugify"`
}

// - base function

func (a *OfficialCommonArticle) Trans() factory.Transactioner {
	return a.base.Trans()
}

func (a *OfficialCommonArticle) Use(trans factory.Transactioner) factory.Model {
	a.base.Use(trans)
	return a
}

func (a *OfficialCommonArticle) SetContext(ctx echo.Context) factory.Model {
	a.base.SetContext(ctx)
	return a
}

func (a *OfficialCommonArticle) EventON(on ...bool) factory.Model {
	a.base.EventON(on...)
	return a
}

func (a *OfficialCommonArticle) EventOFF(off ...bool) factory.Model {
	a.base.EventOFF(off...)
	return a
}

func (a *OfficialCommonArticle) Context() echo.Context {
	return a.base.Context()
}

func (a *OfficialCommonArticle) SetConnID(connID int) factory.Model {
	a.base.SetConnID(connID)
	return a
}

func (a *OfficialCommonArticle) ConnID() int {
	return a.base.ConnID()
}

func (a *OfficialCommonArticle) SetNamer(namer func(factory.Model) string) factory.Model {
	a.base.SetNamer(namer)
	return a
}

func (a *OfficialCommonArticle) Namer() func(factory.Model) string {
	return a.base.Namer()
}

func (a *OfficialCommonArticle) SetParam(param *factory.Param) factory.Model {
	a.base.SetParam(param)
	return a
}

func (a *OfficialCommonArticle) Param(mw func(db.Result) db.Result, args ...interface{}) *factory.Param {
	if a.base.Param() == nil {
		return a.NewParam().SetMiddleware(mw).SetArgs(args...)
	}
	return a.base.Param().SetMiddleware(mw).SetArgs(args...)
}

func (a *OfficialCommonArticle) New(structName string, connID ...int) factory.Model {
	return a.base.New(structName, connID...)
}

func (a *OfficialCommonArticle) Base_() factory.Baser {
	return &a.base
}

// - current function

func (a *OfficialCommonArticle) Objects() []*OfficialCommonArticle {
	if a.objects == nil {
		return nil
	}
	return a.objects[:]
}

func (a *OfficialCommonArticle) XObjects() Slice_OfficialCommonArticle {
	return Slice_OfficialCommonArticle(a.Objects())
}

func (a *OfficialCommonArticle) NewObjects() factory.Ranger {
	return &Slice_OfficialCommonArticle{}
}

func (a *OfficialCommonArticle) InitObjects() *[]*OfficialCommonArticle {
	a.objects = []*OfficialCommonArticle{}
	return &a.objects
}

func (a *OfficialCommonArticle) NewParam() *factory.Param {
	return factory.NewParam(factory.DefaultFactory).SetIndex(a.base.ConnID()).SetTrans(a.base.Trans()).SetCollection(a.Name_()).SetModel(a)
}

func (a *OfficialCommonArticle) Short_() string {
	return "official_common_article"
}

func (a *OfficialCommonArticle) Struct_() string {
	return "OfficialCommonArticle"
}

func (a *OfficialCommonArticle) Name_() string {
	b := a
	if b == nil {
		b = &OfficialCommonArticle{}
	}
	if b.base.Namer() != nil {
		return WithPrefix(b.base.Namer()(b))
	}
	return WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (a *OfficialCommonArticle) CPAFrom(source factory.Model) factory.Model {
	a.SetContext(source.Context())
	a.SetConnID(source.ConnID())
	a.SetNamer(source.Namer())
	return a
}

func (a *OfficialCommonArticle) Get(mw func(db.Result) db.Result, args ...interface{}) (err error) {
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

func (a *OfficialCommonArticle) List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCommonArticle:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonArticle(*v))
		case []*OfficialCommonArticle:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonArticle(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCommonArticle) GroupBy(keyField string, inputRows ...[]*OfficialCommonArticle) map[string][]*OfficialCommonArticle {
	var rows Slice_OfficialCommonArticle
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonArticle(inputRows[0])
	} else {
		rows = Slice_OfficialCommonArticle(a.Objects())
	}
	return rows.GroupBy(keyField)
}

func (a *OfficialCommonArticle) KeyBy(keyField string, inputRows ...[]*OfficialCommonArticle) map[string]*OfficialCommonArticle {
	var rows Slice_OfficialCommonArticle
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonArticle(inputRows[0])
	} else {
		rows = Slice_OfficialCommonArticle(a.Objects())
	}
	return rows.KeyBy(keyField)
}

func (a *OfficialCommonArticle) AsKV(keyField string, valueField string, inputRows ...[]*OfficialCommonArticle) param.Store {
	var rows Slice_OfficialCommonArticle
	if len(inputRows) > 0 {
		rows = Slice_OfficialCommonArticle(inputRows[0])
	} else {
		rows = Slice_OfficialCommonArticle(a.Objects())
	}
	return rows.AsKV(keyField, valueField)
}

func (a *OfficialCommonArticle) ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error) {
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
		case *[]*OfficialCommonArticle:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonArticle(*v))
		case []*OfficialCommonArticle:
			err = DBI.FireReaded(a, queryParam, Slice_OfficialCommonArticle(v))
		case factory.Ranger:
			err = DBI.FireReaded(a, queryParam, v)
		}
	}
	return cnt, err
}

func (a *OfficialCommonArticle) Insert() (pk interface{}, err error) {
	a.Created = uint(time.Now().Unix())
	a.Id = 0
	if len(a.OwnerType) == 0 {
		a.OwnerType = "customer"
	}
	if len(a.Contype) == 0 {
		a.Contype = "markdown"
	}
	if len(a.Display) == 0 {
		a.Display = "Y"
	}
	if len(a.CloseComment) == 0 {
		a.CloseComment = "N"
	}
	if len(a.CommentAutoDisplay) == 0 {
		a.CommentAutoDisplay = "N"
	}
	if len(a.CommentAllowUser) == 0 {
		a.CommentAllowUser = "all"
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

func (a *OfficialCommonArticle) Update(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.OwnerType) == 0 {
		a.OwnerType = "customer"
	}
	if len(a.Contype) == 0 {
		a.Contype = "markdown"
	}
	if len(a.Display) == 0 {
		a.Display = "Y"
	}
	if len(a.CloseComment) == 0 {
		a.CloseComment = "N"
	}
	if len(a.CommentAutoDisplay) == 0 {
		a.CommentAutoDisplay = "N"
	}
	if len(a.CommentAllowUser) == 0 {
		a.CommentAllowUser = "all"
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

func (a *OfficialCommonArticle) Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.OwnerType) == 0 {
		a.OwnerType = "customer"
	}
	if len(a.Contype) == 0 {
		a.Contype = "markdown"
	}
	if len(a.Display) == 0 {
		a.Display = "Y"
	}
	if len(a.CloseComment) == 0 {
		a.CloseComment = "N"
	}
	if len(a.CommentAutoDisplay) == 0 {
		a.CommentAutoDisplay = "N"
	}
	if len(a.CommentAllowUser) == 0 {
		a.CommentAllowUser = "all"
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

func (a *OfficialCommonArticle) UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.OwnerType) == 0 {
		a.OwnerType = "customer"
	}
	if len(a.Contype) == 0 {
		a.Contype = "markdown"
	}
	if len(a.Display) == 0 {
		a.Display = "Y"
	}
	if len(a.CloseComment) == 0 {
		a.CloseComment = "N"
	}
	if len(a.CommentAutoDisplay) == 0 {
		a.CommentAutoDisplay = "N"
	}
	if len(a.CommentAllowUser) == 0 {
		a.CommentAllowUser = "all"
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

func (a *OfficialCommonArticle) UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error) {
	a.Updated = uint(time.Now().Unix())
	if len(a.OwnerType) == 0 {
		a.OwnerType = "customer"
	}
	if len(a.Contype) == 0 {
		a.Contype = "markdown"
	}
	if len(a.Display) == 0 {
		a.Display = "Y"
	}
	if len(a.CloseComment) == 0 {
		a.CloseComment = "N"
	}
	if len(a.CommentAutoDisplay) == 0 {
		a.CommentAutoDisplay = "N"
	}
	if len(a.CommentAllowUser) == 0 {
		a.CommentAllowUser = "all"
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

func (a *OfficialCommonArticle) UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (err error) {
	return a.UpdateFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCommonArticle) UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (affected int64, err error) {
	return a.UpdatexFields(mw, map[string]interface{}{
		field: value,
	}, args...)
}

func (a *OfficialCommonArticle) UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (err error) {

	if val, ok := kvset["owner_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["owner_type"] = "customer"
		}
	}
	if val, ok := kvset["contype"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["contype"] = "markdown"
		}
	}
	if val, ok := kvset["display"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["display"] = "Y"
		}
	}
	if val, ok := kvset["close_comment"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["close_comment"] = "N"
		}
	}
	if val, ok := kvset["comment_auto_display"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["comment_auto_display"] = "N"
		}
	}
	if val, ok := kvset["comment_allow_user"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["comment_allow_user"] = "all"
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

func (a *OfficialCommonArticle) UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (affected int64, err error) {

	if val, ok := kvset["owner_type"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["owner_type"] = "customer"
		}
	}
	if val, ok := kvset["contype"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["contype"] = "markdown"
		}
	}
	if val, ok := kvset["display"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["display"] = "Y"
		}
	}
	if val, ok := kvset["close_comment"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["close_comment"] = "N"
		}
	}
	if val, ok := kvset["comment_auto_display"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["comment_auto_display"] = "N"
		}
	}
	if val, ok := kvset["comment_allow_user"]; ok && val != nil {
		if v, ok := val.(string); ok && len(v) == 0 {
			kvset["comment_allow_user"] = "all"
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

func (a *OfficialCommonArticle) UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) (err error) {
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

func (a *OfficialCommonArticle) Upsert(mw func(db.Result) db.Result, args ...interface{}) (pk interface{}, err error) {
	pk, err = a.Param(mw, args...).SetSend(a).Upsert(func() error {
		a.Updated = uint(time.Now().Unix())
		if len(a.OwnerType) == 0 {
			a.OwnerType = "customer"
		}
		if len(a.Contype) == 0 {
			a.Contype = "markdown"
		}
		if len(a.Display) == 0 {
			a.Display = "Y"
		}
		if len(a.CloseComment) == 0 {
			a.CloseComment = "N"
		}
		if len(a.CommentAutoDisplay) == 0 {
			a.CommentAutoDisplay = "N"
		}
		if len(a.CommentAllowUser) == 0 {
			a.CommentAllowUser = "all"
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
		if len(a.Contype) == 0 {
			a.Contype = "markdown"
		}
		if len(a.Display) == 0 {
			a.Display = "Y"
		}
		if len(a.CloseComment) == 0 {
			a.CloseComment = "N"
		}
		if len(a.CommentAutoDisplay) == 0 {
			a.CommentAutoDisplay = "N"
		}
		if len(a.CommentAllowUser) == 0 {
			a.CommentAllowUser = "all"
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

func (a *OfficialCommonArticle) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {

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

func (a *OfficialCommonArticle) Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error) {

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

func (a *OfficialCommonArticle) Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error) {
	return a.Param(mw, args...).Count()
}

func (a *OfficialCommonArticle) Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error) {
	return a.Param(mw, args...).Exists()
}

func (a *OfficialCommonArticle) Reset() *OfficialCommonArticle {
	a.Id = 0
	a.Category1 = 0
	a.Category2 = 0
	a.Category3 = 0
	a.CategoryId = 0
	a.SourceId = ``
	a.SourceTable = ``
	a.OwnerId = 0
	a.OwnerType = ``
	a.Title = ``
	a.Keywords = ``
	a.Image = ``
	a.ImageOriginal = ``
	a.Summary = ``
	a.Content = ``
	a.Contype = ``
	a.Created = 0
	a.Updated = 0
	a.Display = ``
	a.Template = ``
	a.Comments = 0
	a.CloseComment = ``
	a.CommentAutoDisplay = ``
	a.CommentAllowUser = ``
	a.Likes = 0
	a.Hates = 0
	a.Views = 0
	a.Tags = ``
	a.Price = 0.0
	a.Slugify = ``
	return a
}

func (a *OfficialCommonArticle) AsMap(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["Id"] = a.Id
		r["Category1"] = a.Category1
		r["Category2"] = a.Category2
		r["Category3"] = a.Category3
		r["CategoryId"] = a.CategoryId
		r["SourceId"] = a.SourceId
		r["SourceTable"] = a.SourceTable
		r["OwnerId"] = a.OwnerId
		r["OwnerType"] = a.OwnerType
		r["Title"] = a.Title
		r["Keywords"] = a.Keywords
		r["Image"] = a.Image
		r["ImageOriginal"] = a.ImageOriginal
		r["Summary"] = a.Summary
		r["Content"] = a.Content
		r["Contype"] = a.Contype
		r["Created"] = a.Created
		r["Updated"] = a.Updated
		r["Display"] = a.Display
		r["Template"] = a.Template
		r["Comments"] = a.Comments
		r["CloseComment"] = a.CloseComment
		r["CommentAutoDisplay"] = a.CommentAutoDisplay
		r["CommentAllowUser"] = a.CommentAllowUser
		r["Likes"] = a.Likes
		r["Hates"] = a.Hates
		r["Views"] = a.Views
		r["Tags"] = a.Tags
		r["Price"] = a.Price
		r["Slugify"] = a.Slugify
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "Id":
			r["Id"] = a.Id
		case "Category1":
			r["Category1"] = a.Category1
		case "Category2":
			r["Category2"] = a.Category2
		case "Category3":
			r["Category3"] = a.Category3
		case "CategoryId":
			r["CategoryId"] = a.CategoryId
		case "SourceId":
			r["SourceId"] = a.SourceId
		case "SourceTable":
			r["SourceTable"] = a.SourceTable
		case "OwnerId":
			r["OwnerId"] = a.OwnerId
		case "OwnerType":
			r["OwnerType"] = a.OwnerType
		case "Title":
			r["Title"] = a.Title
		case "Keywords":
			r["Keywords"] = a.Keywords
		case "Image":
			r["Image"] = a.Image
		case "ImageOriginal":
			r["ImageOriginal"] = a.ImageOriginal
		case "Summary":
			r["Summary"] = a.Summary
		case "Content":
			r["Content"] = a.Content
		case "Contype":
			r["Contype"] = a.Contype
		case "Created":
			r["Created"] = a.Created
		case "Updated":
			r["Updated"] = a.Updated
		case "Display":
			r["Display"] = a.Display
		case "Template":
			r["Template"] = a.Template
		case "Comments":
			r["Comments"] = a.Comments
		case "CloseComment":
			r["CloseComment"] = a.CloseComment
		case "CommentAutoDisplay":
			r["CommentAutoDisplay"] = a.CommentAutoDisplay
		case "CommentAllowUser":
			r["CommentAllowUser"] = a.CommentAllowUser
		case "Likes":
			r["Likes"] = a.Likes
		case "Hates":
			r["Hates"] = a.Hates
		case "Views":
			r["Views"] = a.Views
		case "Tags":
			r["Tags"] = a.Tags
		case "Price":
			r["Price"] = a.Price
		case "Slugify":
			r["Slugify"] = a.Slugify
		}
	}
	return r
}

func (a *OfficialCommonArticle) FromRow(row map[string]interface{}) {
	for key, value := range row {
		switch key {
		case "id":
			a.Id = param.AsUint64(value)
		case "category1":
			a.Category1 = param.AsUint(value)
		case "category2":
			a.Category2 = param.AsUint(value)
		case "category3":
			a.Category3 = param.AsUint(value)
		case "category_id":
			a.CategoryId = param.AsUint(value)
		case "source_id":
			a.SourceId = param.AsString(value)
		case "source_table":
			a.SourceTable = param.AsString(value)
		case "owner_id":
			a.OwnerId = param.AsUint64(value)
		case "owner_type":
			a.OwnerType = param.AsString(value)
		case "title":
			a.Title = param.AsString(value)
		case "keywords":
			a.Keywords = param.AsString(value)
		case "image":
			a.Image = param.AsString(value)
		case "image_original":
			a.ImageOriginal = param.AsString(value)
		case "summary":
			a.Summary = param.AsString(value)
		case "content":
			a.Content = param.AsString(value)
		case "contype":
			a.Contype = param.AsString(value)
		case "created":
			a.Created = param.AsUint(value)
		case "updated":
			a.Updated = param.AsUint(value)
		case "display":
			a.Display = param.AsString(value)
		case "template":
			a.Template = param.AsString(value)
		case "comments":
			a.Comments = param.AsUint64(value)
		case "close_comment":
			a.CloseComment = param.AsString(value)
		case "comment_auto_display":
			a.CommentAutoDisplay = param.AsString(value)
		case "comment_allow_user":
			a.CommentAllowUser = param.AsString(value)
		case "likes":
			a.Likes = param.AsUint64(value)
		case "hates":
			a.Hates = param.AsUint64(value)
		case "views":
			a.Views = param.AsUint64(value)
		case "tags":
			a.Tags = param.AsString(value)
		case "price":
			a.Price = param.AsFloat64(value)
		case "slugify":
			a.Slugify = param.AsString(value)
		}
	}
}

func (a *OfficialCommonArticle) Set(key interface{}, value ...interface{}) {
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
		case "Category1":
			a.Category1 = param.AsUint(vv)
		case "Category2":
			a.Category2 = param.AsUint(vv)
		case "Category3":
			a.Category3 = param.AsUint(vv)
		case "CategoryId":
			a.CategoryId = param.AsUint(vv)
		case "SourceId":
			a.SourceId = param.AsString(vv)
		case "SourceTable":
			a.SourceTable = param.AsString(vv)
		case "OwnerId":
			a.OwnerId = param.AsUint64(vv)
		case "OwnerType":
			a.OwnerType = param.AsString(vv)
		case "Title":
			a.Title = param.AsString(vv)
		case "Keywords":
			a.Keywords = param.AsString(vv)
		case "Image":
			a.Image = param.AsString(vv)
		case "ImageOriginal":
			a.ImageOriginal = param.AsString(vv)
		case "Summary":
			a.Summary = param.AsString(vv)
		case "Content":
			a.Content = param.AsString(vv)
		case "Contype":
			a.Contype = param.AsString(vv)
		case "Created":
			a.Created = param.AsUint(vv)
		case "Updated":
			a.Updated = param.AsUint(vv)
		case "Display":
			a.Display = param.AsString(vv)
		case "Template":
			a.Template = param.AsString(vv)
		case "Comments":
			a.Comments = param.AsUint64(vv)
		case "CloseComment":
			a.CloseComment = param.AsString(vv)
		case "CommentAutoDisplay":
			a.CommentAutoDisplay = param.AsString(vv)
		case "CommentAllowUser":
			a.CommentAllowUser = param.AsString(vv)
		case "Likes":
			a.Likes = param.AsUint64(vv)
		case "Hates":
			a.Hates = param.AsUint64(vv)
		case "Views":
			a.Views = param.AsUint64(vv)
		case "Tags":
			a.Tags = param.AsString(vv)
		case "Price":
			a.Price = param.AsFloat64(vv)
		case "Slugify":
			a.Slugify = param.AsString(vv)
		}
	}
}

func (a *OfficialCommonArticle) AsRow(onlyFields ...string) param.Store {
	r := param.Store{}
	if len(onlyFields) == 0 {
		r["id"] = a.Id
		r["category1"] = a.Category1
		r["category2"] = a.Category2
		r["category3"] = a.Category3
		r["category_id"] = a.CategoryId
		r["source_id"] = a.SourceId
		r["source_table"] = a.SourceTable
		r["owner_id"] = a.OwnerId
		r["owner_type"] = a.OwnerType
		r["title"] = a.Title
		r["keywords"] = a.Keywords
		r["image"] = a.Image
		r["image_original"] = a.ImageOriginal
		r["summary"] = a.Summary
		r["content"] = a.Content
		r["contype"] = a.Contype
		r["created"] = a.Created
		r["updated"] = a.Updated
		r["display"] = a.Display
		r["template"] = a.Template
		r["comments"] = a.Comments
		r["close_comment"] = a.CloseComment
		r["comment_auto_display"] = a.CommentAutoDisplay
		r["comment_allow_user"] = a.CommentAllowUser
		r["likes"] = a.Likes
		r["hates"] = a.Hates
		r["views"] = a.Views
		r["tags"] = a.Tags
		r["price"] = a.Price
		r["slugify"] = a.Slugify
		return r
	}
	for _, field := range onlyFields {
		switch field {
		case "id":
			r["id"] = a.Id
		case "category1":
			r["category1"] = a.Category1
		case "category2":
			r["category2"] = a.Category2
		case "category3":
			r["category3"] = a.Category3
		case "category_id":
			r["category_id"] = a.CategoryId
		case "source_id":
			r["source_id"] = a.SourceId
		case "source_table":
			r["source_table"] = a.SourceTable
		case "owner_id":
			r["owner_id"] = a.OwnerId
		case "owner_type":
			r["owner_type"] = a.OwnerType
		case "title":
			r["title"] = a.Title
		case "keywords":
			r["keywords"] = a.Keywords
		case "image":
			r["image"] = a.Image
		case "image_original":
			r["image_original"] = a.ImageOriginal
		case "summary":
			r["summary"] = a.Summary
		case "content":
			r["content"] = a.Content
		case "contype":
			r["contype"] = a.Contype
		case "created":
			r["created"] = a.Created
		case "updated":
			r["updated"] = a.Updated
		case "display":
			r["display"] = a.Display
		case "template":
			r["template"] = a.Template
		case "comments":
			r["comments"] = a.Comments
		case "close_comment":
			r["close_comment"] = a.CloseComment
		case "comment_auto_display":
			r["comment_auto_display"] = a.CommentAutoDisplay
		case "comment_allow_user":
			r["comment_allow_user"] = a.CommentAllowUser
		case "likes":
			r["likes"] = a.Likes
		case "hates":
			r["hates"] = a.Hates
		case "views":
			r["views"] = a.Views
		case "tags":
			r["tags"] = a.Tags
		case "price":
			r["price"] = a.Price
		case "slugify":
			r["slugify"] = a.Slugify
		}
	}
	return r
}

func (a *OfficialCommonArticle) ListPage(cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCommonArticle) ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error {
	_, err := pagination.NewLister(a, recv, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(a.Context())
	return err
}

func (a *OfficialCommonArticle) BatchValidate(kvset map[string]interface{}) error {
	if kvset == nil {
		kvset = a.AsRow()
	}
	return DBI.Fields.BatchValidate(a.Short_(), kvset)
}

func (a *OfficialCommonArticle) Validate(field string, value interface{}) error {
	return DBI.Fields.Validate(a.Short_(), field, value)
}
