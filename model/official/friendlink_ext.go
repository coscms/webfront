package official

import "github.com/coscms/webfront/dbschema"

type FriendlinkExt struct {
	*dbschema.OfficialCommonFriendlink
	Category *dbschema.OfficialCommonCategory `db:"-,relation=id:category_id|gtZero"`
}
