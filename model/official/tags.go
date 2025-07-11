package official

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/top"
)

func NewTags(ctx echo.Context) *Tags {
	return &Tags{
		OfficialCommonTags: dbschema.NewOfficialCommonTags(ctx),
	}
}

type Tags struct {
	*dbschema.OfficialCommonTags
}

func (f *Tags) ListByGroup(group string, limit int) []*dbschema.OfficialCommonTags {
	f.ListByOffset(nil, nil, 0, limit, db.Cond{`group`: group})
	return f.Objects()
}

func (f *Tags) check() error {
	f.Name = strings.TrimSpace(f.Name)
	if len(f.Name) == 0 {
		return f.Context().NewError(code.InvalidParameter, `标签名无效`).SetZone(`name`)
	}
	f.Group = strings.TrimSpace(f.Group)
	return nil
}

func (f *Tags) Add() (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return nil, err
	}
	cond := db.NewCompounds()
	cond.AddKV(`name`, f.Name)
	cond.AddKV(`group`, f.Group)
	m := dbschema.NewOfficialCommonTags(f.Context())
	err = m.Get(nil, cond.And())
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		return f.OfficialCommonTags.Insert()
	}

	if f.Num != 0 {
		err = f.IncrNum(f.Group, f.Name, int(f.Num))
	}
	return
}

func (f *Tags) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialCommonTags.Update(mw, args...)
}

func (f *Tags) Clean() (err error) {
	return f.Delete(nil, db.Cond{`num`: 0})
}

func (f *Tags) IncrNum(group, name string, n ...int) error {
	var _n int
	if len(n) > 0 {
		_n = n[0]
	} else {
		_n = 1
	}
	if _n == 0 {
		return nil
	}
	err := f.UpdateField(nil, `num`, db.Raw(`num+`+param.AsString(_n)), db.And(
		db.Cond{`name`: name},
		db.Cond{`group`: group},
	))
	return err
}

func (f *Tags) IncrNumByNames(group string, names []string, n ...int) error {
	var _n int
	if len(n) > 0 {
		_n = n[0]
	} else {
		_n = 1
	}
	if _n == 0 {
		return nil
	}
	err := f.UpdateField(nil, `num`, db.Raw(`num+`+param.AsString(_n)), db.And(
		db.Cond{`name`: db.In(names)},
		db.Cond{`group`: group},
	))
	return err
}

func (f *Tags) DecrNum(group string, name []string, n ...int) error {
	var _n int
	if len(n) > 0 {
		_n = n[0]
	} else {
		_n = 1
	}
	if _n == 0 {
		return nil
	}
	err := f.UpdateField(nil, `num`, db.Raw(`num-`+param.AsString(_n)), db.And(
		db.Cond{`name`: db.In(name)},
		db.Cond{`group`: group},
		db.Cond{`num`: db.Gt(0)},
	))
	return err
}

func (f *Tags) UpdateTags(group string, oldTags []string, postTags []string, disallowCreateTags ...bool) ([]string, error) {
	var (
		delTags        []string
		tags           = make([]string, 0, len(postTags))
		err            error
		disallowCreate bool
	)
	if len(disallowCreateTags) > 0 {
		disallowCreate = disallowCreateTags[0]
	}
	if oldTags == nil {
		oldTags = []string{}
	}
	uniqueTags := map[string]int{}
	if len(postTags) > 0 && len(group) > 0 {
		//获取提交tag的唯一值
		for idx, tag := range postTags {
			tag = strings.TrimSpace(tag)
			if len(tag) == 0 {
				continue
			}
			if _, y := uniqueTags[tag]; !y {
				uniqueTags[tag] = idx
				tags = append(tags, tag)
			}
		}
	}
	if len(tags) > 0 { //如果有提交tags
		delTags = com.StringSliceDiff(oldTags, tags) // 比较出被删除的tags
		_, err = f.ListByOffset(nil, nil, 0, -1, db.And(
			db.Cond{`group`: group},
			db.Cond{`name`: db.In(tags)},
		))
		if err != nil {
			return nil, err
		}
		// 找出新增tags
		for _, tagRow := range f.Objects() {
			delete(uniqueTags, tagRow.Name)         //从提交的tags中清除掉已经存在的tags
			if !com.InSlice(tagRow.Name, oldTags) { //旧数据tags中没有时，代表本次新增了此标签
				err = f.IncrNum(group, tagRow.Name)
				if err != nil {
					return nil, err
				}
			}
		}
		if disallowCreate { // 不允许创建新标签，则需要清理掉新标签
			var filtered []string
			for _, tag := range tags {
				if _, ok := uniqueTags[tag]; !ok {
					filtered = append(filtered, tag)
				}
			}
			tags = filtered
		} else {
			// 添加新tags
			for name := range uniqueTags {
				f.Reset()
				f.Name = name
				f.Group = group
				f.Num = 1
				f.Display = `Y`
				_, err = f.OfficialCommonTags.Insert()
				if err != nil {
					return nil, err
				}
			}
		}
	} else { //如果没有提交tags，则删除旧tags
		delTags = oldTags
	}
	if len(delTags) > 0 { // 删除标签时才减去使用次数
		err = f.DecrNum(group, delTags)
	}
	return tags, err
}

func TagCond(tag string) db.Compound {
	return top.CondFindInSet(`tags`, tag)
}
