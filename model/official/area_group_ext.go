package official

import "github.com/coscms/webfront/dbschema"

type AreaGroupExt struct {
	*dbschema.OfficialCommonAreaGroup
	Areas []*dbschema.OfficialCommonArea `db:"-,relation=id:area_ids|notEmpty|split"`
}
