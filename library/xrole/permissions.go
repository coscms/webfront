package xrole

import (
	"github.com/coscms/webcore/library/role"
	"github.com/coscms/webfront/dbschema"
)

type CustomerRoleWithPermissions struct {
	*dbschema.OfficialCustomerRole
	Permissions []*dbschema.OfficialCustomerRolePermission `db:"-,relation=role_id:id"`
}

func (u *CustomerRoleWithPermissions) GetPermissions() []role.PermissionConfiger {
	r := make([]role.PermissionConfiger, len(u.Permissions))
	for k, v := range u.Permissions {
		r[k] = v
	}
	return r
}
