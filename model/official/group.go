package official

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webfront/dbschema"
)

func NewGroup(ctx echo.Context) *Group {
	return &Group{
		OfficialCommonGroup: dbschema.NewOfficialCommonGroup(ctx),
	}
}

type Group struct {
	*dbschema.OfficialCommonGroup
}

func (f *Group) Exists(name string) (bool, error) {
	return f.OfficialCommonGroup.Exists(nil, db.Cond{`name`: name})
}

func (f *Group) ExistsOther(name string, id uint) (bool, error) {
	return f.OfficialCommonGroup.Exists(nil, db.Cond{`name`: name, `id <>`: id})
}

func (f *Group) check() error {
	f.Name = strings.TrimSpace(f.Name)
	var err error
	if len(f.Name) == 0 {
		err = f.Context().NewError(code.InvalidParameter, `用户组名称不能为空`).SetZone(`name`)
		return err
	}
	var exists bool
	if f.Id > 0 {
		exists, err = f.ExistsOther(f.Name, f.Id)
	} else {
		exists, err = f.Exists(f.Name)
	}
	if err != nil {
		return err
	}
	if exists {
		err = f.Context().NewError(code.DataAlreadyExists, `用户组名称“%s”已经存在`, f.Name).SetZone(`name`)
		return err
	}
	return err
}

func (f *Group) Add() (pk interface{}, err error) {
	err = f.check()
	if err != nil {
		return
	}
	return f.OfficialCommonGroup.Insert()
}

func (f *Group) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	err := f.check()
	if err != nil {
		return err
	}
	return f.OfficialCommonGroup.Update(mw, args...)
}
