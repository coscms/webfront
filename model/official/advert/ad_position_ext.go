package advert

import "github.com/coscms/webfront/dbschema"

var _ Adverter = (*PositionWithRendered)(nil)

type PositionWithRendered struct {
	*dbschema.OfficialAdPosition
	Rendered string
}

func (i *PositionWithRendered) GetWidth() uint {
	return i.Width
}

func (i *PositionWithRendered) GetHeight() uint {
	return i.Height
}

func (i *PositionWithRendered) GetURL() string {
	return i.Url
}

func (i *PositionWithRendered) GetContent() string {
	return i.Content
}

func (i *PositionWithRendered) GetContype() string {
	return i.Contype
}

func (i *PositionWithRendered) GetTitle() string {
	return i.Title
}

func (i *PositionWithRendered) GetDescription() string {
	return i.Description
}
