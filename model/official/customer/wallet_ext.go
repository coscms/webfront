package customer

import "github.com/coscms/webfront/dbschema"

type WalletExt struct {
	*dbschema.OfficialCustomerWallet
	Customer      *CustomerBase `db:"-,relation=id:customer_id"`
	AssetTypeName string        `db:"-"`
}

type WalletFlowExt struct {
	*dbschema.OfficialCustomerWalletFlow
	Customer      *CustomerBase `db:"-,relation=id:customer_id"`
	SrcCustomer   *CustomerBase `db:"-,relation=id:source_customer|gtZero"`
	AssetTypeName string        `db:"-"`
}
