package article

import (
	"fmt"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/download/downloadByContent"
	"github.com/coscms/webfront/library/top"
	"github.com/coscms/webfront/library/xcommon"
	"github.com/coscms/webfront/library/xrole"
	"github.com/coscms/webfront/library/xrole/xroleutils"
	"github.com/coscms/webfront/model/official"
)

func NewArticle(ctx echo.Context) *Article {
	m := &Article{
		OfficialCommonArticle: dbschema.NewOfficialCommonArticle(ctx),
	}
	return m
}

type Article struct {
	*dbschema.OfficialCommonArticle
	DisallowCreateTags bool
}

func (f *Article) ListByTag(recv interface{}, mw func(db.Result) db.Result, page, size int, tag string) (func() int64, error) {
	return f.OfficialCommonArticle.List(recv, mw, page, size, f.TagCond(tag))
}

func (f *Article) GetCategories(rows ...*dbschema.OfficialCommonArticle) ([]dbschema.OfficialCommonCategory, error) {
	row := f.OfficialCommonArticle
	if len(rows) > 0 {
		row = rows[0]
	}
	cateM := official.NewCategory(f.Context())
	return cateM.Positions(row.CategoryId)
}

func (f *Article) TagCond(tag string) db.Compound {
	return official.TagCond(tag)
}

func (f *Article) check(old *dbschema.OfficialCommonArticle) error {
	if len(f.Title) < 1 {
		return f.Context().NewError(code.InvalidParameter, `标题不能为空`).SetZone(`title`)
	}
	if f.Price < 0 {
		return f.Context().NewError(code.InvalidParameter, `价格不能为负数`).SetZone(`price`)
	}
	if f.CategoryId > 0 {
		cateM := official.NewCategory(f.Context())
		parentIDs, err := cateM.PositionIds(f.CategoryId)
		if err != nil {
			return err
		}
		f.Category3 = 0
		f.Category2 = 0
		f.Category1 = 0
		switch len(parentIDs) {
		case 4:
			//f.CategoryId = parentIDs[3]
			fallthrough
		case 3:
			f.Category3 = parentIDs[2]
			fallthrough
		case 2:
			f.Category2 = parentIDs[1]
			fallthrough
		case 1:
			f.Category1 = parentIDs[0]
		}
	}
	var err error
	var oldTags []string
	if old != nil && len(old.Tags) > 0 {
		oldTags = strings.Split(old.Tags, `,`)
	}
	tagsM := official.NewTags(f.Context())
	tags, err := tagsM.UpdateTags(f.Id == 0, GroupName, oldTags, strings.Split(f.Tags, `,`), f.DisallowCreateTags)
	if err != nil {
		return err
	}
	f.Tags = strings.Join(tags, `,`)

	if len(f.Contype) == 0 || !Contype.Has(f.Contype) {
		f.Contype = `text`
	}
	f.Content = common.ContentEncode(f.Content, f.Contype)
	f.Slugify = top.Slugify(f.Title)
	return err
}

func (f *Article) SyncRemoteImage() (string, error) {
	newContent, err := downloadByContent.SyncRemoteImage(f.Context(), `default`, fmt.Sprint(f.Id), f.Content, f.Contype)
	return newContent, err
}

func (f *Article) CustomerTodayCount(customerID interface{}) (int64, error) {
	startTs, endTs := top.TodayTimestamp()
	return f.Count(nil, db.And(
		db.Cond{`owner_type`: `customer`},
		db.Cond{`owner_id`: customerID},
		db.Cond{`created`: db.Between(startTs, endTs)},
	))
}

func (f *Article) CustomerPendingCount(customerID interface{}) (int64, error) {
	return f.Count(nil, db.And(
		db.Cond{`owner_type`: `customer`},
		db.Cond{`owner_id`: customerID},
		db.Cond{`display`: `N`},
	))
}

func (f *Article) CustomerPendingTodayCount(customerID interface{}) (int64, error) {
	startTs, endTs := top.TodayTimestamp()
	return f.Count(nil, db.And(
		db.Cond{`owner_type`: `customer`},
		db.Cond{`owner_id`: customerID},
		db.Cond{`display`: `N`},
		db.Cond{`created`: db.Between(startTs, endTs)},
	))
}

func (f *Article) checkCustomerAdd(permission *xrole.RolePermission) error {
	err := xcommon.CheckRoleCustomerAdd(f.Context(), permission, BehaviorName, f.OwnerId, f)
	if err == nil {
		return err
	}
	switch err {
	case xcommon.ErrCustomerRoleDisabled:
		return f.Context().E(`当前角色不支持文章投稿`)
	case xcommon.ErrCustomerAddClosed:
		return f.Context().E(`文章投稿功能已关闭`)
	case xcommon.ErrCustomerAddMaxPerDay:
		return f.Context().E(`投稿失败。您的账号已达到今日最大投稿数量`)
	case xcommon.ErrCustomerAddMaxPending:
		return f.Context().E(`投稿失败。您的待审核文章数量已达上限，请等待审核通过后再投稿`)
	default:
		return err
	}
}

func (f *Article) Add() (pk interface{}, err error) {
	if f.OwnerType == `customer` {
		permission := xroleutils.CustomerPermission(f.Context())
		if err = f.checkCustomerAdd(permission); err != nil {
			return nil, err
		}
	}
	f.Context().Begin()
	if err = f.check(nil); err != nil {
		f.Context().Rollback()
		return
	}
	syncRemoteImage := f.Context().Formx(`syncRemoteImage`).Bool()
	if syncRemoteImage {
		f.Content, err = f.SyncRemoteImage()
		if err != nil {
			f.Context().Rollback()
			return
		}
	}
	pk, err = f.OfficialCommonArticle.Insert()
	if err != nil {
		f.Context().Rollback()
		return
	}
	err = f.Context().Commit()
	return
}

func (f *Article) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	old := dbschema.NewOfficialCommonArticle(f.Context())
	err := old.Get(nil, args...)
	if err != nil {
		return err
	}
	if err := f.check(old); err != nil {
		return err
	}
	syncRemoteImage := f.Context().Formx(`syncRemoteImage`).Bool()
	if syncRemoteImage {
		f.Content, err = f.SyncRemoteImage()
		if err != nil {
			return err
		}
	}
	err = f.OfficialCommonArticle.Update(mw, args...)
	return err
}

func (f *Article) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.OfficialCommonArticle.Delete(mw, args...)
	return err
}

var OrderQuerier = func(c echo.Context, customer *dbschema.OfficialCustomer, sourceId, sourceTable string) error {
	return c.E(`很抱歉，您不是当前商品买家，本文只有当前商品买家才能评论`)
}

// IsAllowedComment 是否可以评论
func (f *Article) IsAllowedComment(customer *dbschema.OfficialCustomer) error {
	articleM := f
	c := f.Context()
	if articleM.CloseComment == `Y` {
		return c.E(`本文已经关闭评论`)
	}
	if articleM.CommentAllowUser == `all` {
		return nil
	}
	if execute, y := AllowComment[articleM.CommentAllowUser]; y {
		return execute(c, customer, f.OfficialCommonArticle)
	}
	return c.E(`很抱歉，您没有权限评论此文`)
}

func (f *Article) ListPageSimple(cond *db.Compounds, orderby ...interface{}) ([]*ArticleWithOwner, error) {
	if len(orderby) == 0 {
		orderby = []interface{}{`-id`}
	}
	rows := []*ArticleWithOwner{}
	_, err := common.NewLister(f, &rows, func(r db.Result) db.Result {
		return r.Select(factory.DBIGet().OmitSelect(f.OfficialCommonArticle, `content`)...).OrderBy(orderby...)
	}, cond.And()).Paging(f.Context())
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (f *Article) ListByOffsetSimple(cond *db.Compounds, limit int, offset int, orderby ...interface{}) ([]*ArticleWithOwner, error) {
	if len(orderby) == 0 {
		orderby = []interface{}{`-id`}
	}
	rows := []*ArticleWithOwner{}
	_, err := f.OfficialCommonArticle.ListByOffset(&rows, func(r db.Result) db.Result {
		return r.Select(factory.DBIGet().OmitSelect(f.OfficialCommonArticle, `content`)...).OrderBy(orderby...)
	}, offset, limit, cond.And())
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (f *Article) ListPage(cond *db.Compounds, orderby ...interface{}) ([]*ArticleAndSourceInfo, error) {
	list := []*ArticleAndSourceInfo{}
	_, err := common.NewLister(f, &list, func(r db.Result) db.Result {
		return r.Select(factory.DBIGet().OmitSelect(f, `content`)...).OrderBy(orderby...)
	}, cond.And()).Paging(f.Context())
	if err != nil {
		return list, err
	}
	err = WithSourceInfo(f.Context(), list)
	if err != nil {
		return nil, err
	}
	tg := make([]official.ICategory, len(list))
	for k, v := range list {
		tg[k] = v
	}
	categoryM := official.NewCategory(f.Context())
	err = categoryM.FillTo(tg)
	return list, err
}

func (f *Article) NextRow(currentID uint64, extraCond *db.Compounds) (*dbschema.OfficialCommonArticle, error) {
	row := dbschema.NewOfficialCommonArticle(nil)
	row.CPAFrom(f.OfficialCommonArticle)
	cond := db.NewCompounds()
	cond.AddKV(`display`, `Y`)
	cond.AddKV(`id`, db.Lt(currentID))
	if extraCond != nil {
		cond.From(extraCond)
	}
	err := row.Get(func(r db.Result) db.Result {
		return r.Select(`id`, `title`, `image`, `created`).OrderBy(`-id`)
	}, cond.And())
	return row, err
}

func (f *Article) PrevRow(currentID uint64, extraCond *db.Compounds) (*dbschema.OfficialCommonArticle, error) {
	row := dbschema.NewOfficialCommonArticle(nil)
	row.CPAFrom(f.OfficialCommonArticle)
	cond := db.NewCompounds()
	cond.AddKV(`display`, `Y`)
	cond.AddKV(`id`, db.Gt(currentID))
	if extraCond != nil {
		cond.From(extraCond)
	}
	err := row.Get(func(r db.Result) db.Result {
		return r.Select(`id`, `title`, `image`, `created`).OrderBy(`id`)
	}, cond.And())
	return row, err
}

func (f *Article) RelationList(limit int, orderby ...interface{}) []*ArticleWithOwner {
	cond := db.NewCompounds()
	if len(f.SourceTable) > 0 {
		cond.Add(db.Cond{`source_table`: f.SourceTable})
		if len(f.SourceId) > 0 {
			cond.Add(db.Cond{`source_id`: f.SourceId})
		}
	}
	return f.CommonQueryList(cond, limit, 0, orderby...)
}

func (f *Article) CommonQueryList(cond *db.Compounds, limit int, offset int, orderby ...interface{}) []*ArticleWithOwner {
	cond.Add(db.Cond{`display`: `Y`})
	rows, _ := f.ListByOffsetSimple(cond, limit, offset, orderby...)
	return rows
}

func (f *Article) QueryList(query string, limit int, offset int, orderby ...interface{}) []*ArticleWithOwner {
	cond := db.NewCompounds()
	r := common.NewSortedURLValues(query)
	r.ApplyCond(cond)
	return f.CommonQueryList(cond, limit, offset, orderby...)
}
