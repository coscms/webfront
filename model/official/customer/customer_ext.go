package customer

import (
	"github.com/webx-top/db/factory"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webfront/dbschema"
	modelLevel "github.com/coscms/webfront/model/official/level"
)

type CustomerBase struct {
	Id     uint64 `db:"id" bson:"id,omitempty" comment:"ID" json:"id" xml:"id"`
	Name   string `db:"name" bson:"name" comment:"名称" json:"name" xml:"name"`
	Gender string `db:"gender" bson:"gender" comment:"性别(male-男;female-女;secret-保密)" json:"gender" xml:"gender"`
	Avatar string `db:"avatar" bson:"avatar" comment:"头像" json:"avatar" xml:"avatar"`
}

func (c *CustomerBase) Short_() string {
	return "official_customer"
}

func (c *CustomerBase) Name_() string {
	b := c
	if b == nil {
		b = &CustomerBase{}
	}
	return dbschema.WithPrefix(factory.TableNamerGet(b.Short_())(b))
}

func (c *CustomerBase) SelectColumns() []interface{} {
	return []interface{}{`id`, `name`, `gender`, `avatar`}
}

type CustomerAndGroup struct {
	*dbschema.OfficialCustomer
	Group  *dbschema.OfficialCommonGroup    `db:"-,relation=id:group_id|gtZero"`
	Levels []*modelLevel.RelationExt        `db:"-,relation=customer_id:id|gtZero" json:",omitempty"`
	Agent  map[string]interface{}           `db:"-" json:",omitempty" xml:",omitempty"`
	Roles  []*dbschema.OfficialCustomerRole `db:"-,relation=id:role_ids|notEmpty|split"`
}

func (d *CustomerAndGroup) AsMap() param.Store {
	m := d.OfficialCustomer.AsMap()
	if d.Group != nil {
		m[`Group`] = d.Group.AsMap()
	}
	if len(d.Roles) > 0 {
		roles := make([]param.Store, len(d.Roles))
		for k, v := range d.Roles {
			roles[k] = v.AsMap()
		}
		m[`Roles`] = roles
	}
	return m
}
