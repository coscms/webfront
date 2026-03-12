package advert

import (
	"html/template"
	"math/rand"
	"time"

	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/com"
)

var _ Adverter = (*ItemResponse)(nil)

type ItemResponse struct {
	Content     string `json:"content" xml:"content"`
	Contype     string `json:"contype" xml:"contype"`
	URL         string `json:"url" xml:"url"`
	Start       uint   `json:"start,omitempty" xml:"start,omitempty"`
	End         uint   `json:"end,omitempty" xml:"end,omitempty"`
	Width       uint   `json:"width,omitempty" xml:"width,omitempty"`
	Height      uint   `json:"height,omitempty" xml:"height,omitempty"`
	Rendered    string `json:"rendered,omitempty" xml:"rendered,omitempty"`
	Title       string `json:"title" xml:"title"`
	Description string `json:"description" xml:"description"`
}

func (i *ItemResponse) GetWidth() uint {
	if i == nil {
		return 0
	}
	return i.Width
}

func (i *ItemResponse) GetHeight() uint {
	if i == nil {
		return 0
	}
	return i.Height
}

func (i *ItemResponse) GetURL() string {
	if i == nil {
		return ``
	}
	return i.URL
}

func (i *ItemResponse) GetContent() string {
	if i == nil {
		return ``
	}
	return i.Content
}

func (i *ItemResponse) GetContype() string {
	if i == nil {
		return ``
	}
	return i.Contype
}

func (i *ItemResponse) GetTitle() string {
	if i == nil {
		return ``
	}
	return i.Title
}

func (i *ItemResponse) GetDescription() string {
	if i == nil {
		return ``
	}
	return i.Description
}

func (i *ItemResponse) GenHTML() *ItemResponse {
	if i == nil {
		return i
	}
	i.Rendered = Render(i)
	return i
}

func (i *ItemResponse) HTML() template.HTML {
	if i == nil {
		return template.HTML(``)
	}
	if len(i.Rendered) == 0 {
		i.Rendered = Render(i)
	}
	return template.HTML(i.Rendered)
}

func NewItemResponse(item *dbschema.OfficialAdItem, position *dbschema.OfficialAdPosition) *ItemResponse {
	resp := ItemResponse{
		Content:     item.Content,
		Contype:     item.Contype,
		URL:         item.Url,
		Start:       item.Start,
		End:         item.End,
		Width:       position.Width,
		Height:      position.Height,
		Title:       item.Title,
		Description: item.Description,
	}
	if len(resp.Title) == 0 {
		resp.Title = position.Title
	}
	if len(resp.Description) == 0 {
		resp.Description = position.Description
	}
	return &resp
}

type ItemsResponse []*ItemResponse

func (i ItemsResponse) Rand() *ItemResponse {
	switch len(i) {
	case 1:
		return i[0]
	case 0:
		return nil
	default:
		return com.RandSlicex(i)
	}
}

func (i ItemsResponse) Shuffle() ItemsResponse {
	r := make(ItemsResponse, len(i))
	copy(r, i)
	var random = rand.New(rand.NewSource(time.Now().Unix()))
	for index := len(r) - 1; index > 0; index-- {
		chooseIndex := random.Intn(index + 1)
		r[chooseIndex], r[index] = r[index], r[chooseIndex]
	}
	return r
}

func (i ItemsResponse) Valid() bool {
	return len(i) > 0
}

func (i ItemsResponse) Size() int {
	return len(i)
}

func (i ItemsResponse) First() *ItemResponse {
	if len(i) < 1 {
		return nil
	}
	return i[0]
}

func (i ItemsResponse) Last() *ItemResponse {
	if len(i) < 1 {
		return nil
	}
	return i[len(i)-1]
}

func (i *ItemsResponse) GenHTML() *ItemsResponse {
	for _, ad := range *i {
		ad.GenHTML()
	}
	return i
}

func (c *ItemsResponse) Place(_ string) *ItemsResponse {
	return c
}

type PositionAdverts map[string]ItemsResponse

func (c *PositionAdverts) GenHTML() *PositionAdverts {
	for _, adList := range *c {
		adList.GenHTML()
	}
	return c
}

func (c PositionAdverts) Place(ident string) ItemsResponse {
	return c[ident]
}

func (c PositionAdverts) Valid() bool {
	return len(c) > 0
}
