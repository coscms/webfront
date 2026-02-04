package official

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
)

func NewAreaCountry(ctx echo.Context) *AreaCountry {
	return &AreaCountry{
		OfficialCommonAreaCountry: dbschema.NewOfficialCommonAreaCountry(ctx),
	}
}

type AreaCountry struct {
	*dbschema.OfficialCommonAreaCountry
}

func (f *AreaCountry) Exists(abbr string) (bool, error) {
	return f.OfficialCommonAreaCountry.Exists(nil, db.Cond{`abbr`: abbr})
}

func (f *AreaCountry) ExistsOther(abbr string, id uint) (bool, error) {
	return f.OfficialCommonAreaCountry.Exists(nil, db.And(
		db.Cond{`abbr`: abbr},
		db.Cond{`id`: db.NotEq(id)},
	))
}

func (f *AreaCountry) check() error {
	f.Name = strings.TrimSpace(f.Name)
	if len(f.Name) == 0 {
		return f.Context().NewError(code.InvalidParameter, `请输入国家名称`).SetZone(`name`)
	}
	f.Short = strings.TrimSpace(f.Short)
	if len(f.Short) == 0 {
		return f.Context().NewError(code.InvalidParameter, `请输入国家简称`).SetZone(`short`)
	}
	f.Abbr = strings.TrimSpace(f.Abbr)
	if len(f.Abbr) == 0 {
		return f.Context().NewError(code.InvalidParameter, `请输入国家缩写`).SetZone(`abbr`)
	}
	if len(f.Abbr) != 2 || !com.StrIsAlpha(f.Abbr) {
		return f.Context().NewError(code.InvalidParameter, `请输入两个字母的国家缩写`).SetZone(`abbr`)
	}
	f.Lng = strings.TrimSpace(f.Lng)
	f.Lat = strings.TrimSpace(f.Lat)
	f.Abbr = strings.ToUpper(f.Abbr)
	f.Disabled = common.GetBoolFlag(f.Disabled)
	f.Code = strings.TrimSpace(f.Code)
	var (
		exists bool
		err    error
	)
	if f.Id < 1 {
		exists, err = f.Exists(f.Abbr)
	} else {
		exists, err = f.ExistsOther(f.Abbr, f.Id)
	}
	if err != nil {
		return err
	}
	if exists {
		return f.Context().E(`缩写“%s”已存在`, f.Abbr)
	}
	return err
}

func (f *AreaCountry) Add() (pk interface{}, err error) {
	if err = f.check(); err != nil {
		return nil, err
	}
	return f.OfficialCommonAreaCountry.Insert()
}

func (f *AreaCountry) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	return f.OfficialCommonAreaCountry.Update(mw, args...)
}
