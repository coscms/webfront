package customer

import "github.com/coscms/webfront/dbschema"

type InvitationCustomerExt struct {
	*dbschema.OfficialCustomerInvitationUsed
	Customer   *dbschema.OfficialCustomer       `db:"-,relation=id:customer_id|gtZero"`
	Level      *dbschema.OfficialCustomerLevel  `db:"-,relation=id:level_id|gtZero"`
	AgentLevel map[string]interface{}           `db:"-" json:",omitempty" xml:",omitempty"`
	RoleList   []*dbschema.OfficialCustomerRole `db:"-,relation=id:role_ids|split"`
}

type InvitationCustomerWithCode struct {
	*dbschema.OfficialCustomerInvitationUsed
	Customer   *dbschema.OfficialCustomer           `db:"-,relation=id:customer_id|gtZero"`
	Level      *dbschema.OfficialCustomerLevel      `db:"-,relation=id:level_id|gtZero"`
	AgentLevel map[string]interface{}               `db:"-" json:",omitempty" xml:",omitempty"`
	RoleList   []*dbschema.OfficialCustomerRole     `db:"-,relation=id:role_ids|split"`
	Invitation *dbschema.OfficialCustomerInvitation `db:"-,relation=id:invitation_id|gtZero"`
}
