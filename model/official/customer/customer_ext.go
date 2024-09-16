package customer

import (
	"github.com/webx-top/echo/param"

	"github.com/coscms/webfront/dbschema"
	modelLevel "github.com/coscms/webfront/model/official/level"
)

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
