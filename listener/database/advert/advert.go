package advert

import (
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/model/official/advert"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
)

func init() {
	// - official_ad_item
	dbschema.DBI.On(factory.EventDeleting, func(m factory.Model, _ ...string) error {
		fm := m.(*dbschema.OfficialAdItem)
		if fm.PositionId > 0 && fm.Disabled == common.BoolN {
			posM := dbschema.NewOfficialAdPosition(fm.Context())
			posM.Get(func(r db.Result) db.Result {
				return r.Select(`ident`)
			}, `id`, fm.PositionId)
			if len(posM.Ident) > 0 {
				advert.DeleteCachedAdvert(fm.Context(), posM.Ident)
			}
		}
		return nil
	}, `official_ad_item`)
	dbschema.DBI.On(factory.EventCreated, func(m factory.Model, _ ...string) error {
		fm := m.(*dbschema.OfficialAdItem)
		if fm.PositionId > 0 && fm.Disabled == common.BoolN {
			posM := dbschema.NewOfficialAdPosition(fm.Context())
			posM.Get(func(r db.Result) db.Result {
				return r.Select(`ident`)
			}, `id`, fm.PositionId)
			if len(posM.Ident) > 0 {
				advert.DeleteCachedAdvert(fm.Context(), posM.Ident)
			}
		}
		return nil
	}, `official_ad_item`)
	dbschema.DBI.On(factory.EventUpdated, func(m factory.Model, _ ...string) error {
		fm := m.(*dbschema.OfficialAdItem)
		if fm.PositionId > 0 {
			posM := dbschema.NewOfficialAdPosition(fm.Context())
			posM.Get(func(r db.Result) db.Result {
				return r.Select(`ident`)
			}, `id`, fm.PositionId)
			if len(posM.Ident) > 0 {
				advert.DeleteCachedAdvert(fm.Context(), posM.Ident)
			}
		}
		return nil
	}, `official_ad_item`)

	// - official_ad_position
	dbschema.DBI.On(factory.EventDeleting, func(m factory.Model, _ ...string) error {
		fm := m.(*dbschema.OfficialAdPosition)
		if len(fm.Ident) > 0 {
			advert.DeleteCachedAdvert(fm.Context(), fm.Ident)
		}
		return nil
	}, `official_ad_position`)
	dbschema.DBI.On(factory.EventUpdated, func(m factory.Model, _ ...string) error {
		fm := m.(*dbschema.OfficialAdPosition)
		if len(fm.Ident) > 0 {
			advert.DeleteCachedAdvert(fm.Context(), fm.Ident)
		}
		return nil
	}, `official_ad_position`)
}
