package perm

import (
	"github.com/coscms/webfront/library/xrole"
)

func New() *RolePermission {
	return xrole.NewRolePermission()
}

type RolePermission = xrole.RolePermission
