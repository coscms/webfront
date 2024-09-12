package customer

import (
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/echo"
)

type FollowingAndCustomer struct {
	*dbschema.OfficialCustomerFollowing
	Customer echo.H
}
