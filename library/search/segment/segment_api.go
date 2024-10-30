package segment

import (
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/restyclient"
)

func NewAPI(apiURL string, apiKey string) *apiSegment {
	return &apiSegment{apiURL: apiURL, apiKey: apiKey}
}

type apiSegment struct {
	apiURL string
	apiKey string
}

func (s *apiSegment) LoadDict(dictFile string, dictType ...string) error {
	return nil
}

func (s *apiSegment) Segment(text string, args ...string) []string {
	return s.request(text, ``, args...)
}

func (s *apiSegment) request(text string, mode string, args ...string) []string {
	data := map[string]string{
		`text`: text,
	}
	if len(mode) > 0 {
		data[`mode`] = mode
	}
	if len(args) > 0 {
		data[`args`] = strings.Join(args, `,`)
	}
	req := restyclient.Retryable()
	if len(s.apiKey) > 0 {
		b, _ := com.JSONEncode(data)
		token := com.Token(s.apiKey, b)
		req.SetAuthToken(token)
	}
	response, err := req.SetFormData(data).Post(s.apiURL)
	if err != nil {
		log.Errorf(`%s: %v`, s.apiURL, err)
		return nil
	}
	if response.IsError() {
		log.Errorf(`%s: [%d] %v`, s.apiURL, response.StatusCode(), com.StripTags(com.Bytes2str(response.Body())))
		return nil
	}
	result := response.String()
	if len(result) == 0 {
		return []string{}
	}
	return strings.Split(result, ` `)
}

func (s *apiSegment) SegmentBy(text string, mode string, args ...string) []string {
	return s.request(text, mode, args...)
}

func (s *apiSegment) Close() error {
	return nil
}
