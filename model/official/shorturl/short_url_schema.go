package shorturl

import (
	"github.com/coscms/webfront/dbschema"
)

type ShortURLVisitWithURL struct {
	*dbschema.OfficialShortUrlVisit
	Num uint64 `db:"num" bson:"num" comment:"数量" json:"num" xml:"num"`
	URL *dbschema.OfficialShortUrl
}
