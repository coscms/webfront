package segment

import (
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/restyclient"
)

// NewAPI creates a new apiSegment instance with the provided API URL and API key.
func NewAPI(apiURL string, apiKey string) *apiSegment {
	return &apiSegment{apiURL: apiURL, apiKey: apiKey}
}

type apiSegment struct {
	apiURL string
	apiKey string
}

// LoadDict loads a dictionary file of specified type for the segmenter.
// dictFile specifies the path to the dictionary file.
// dictType optionally specifies the dictionary type (e.g., "main", "stopwords").
// Returns an error if the dictionary fails to load.
func (s *apiSegment) LoadDict(dictFile string, dictType ...string) error {
	return nil
}

// Segment splits the input text into segments using the configured segmentation method.
// The optional args parameter allows for additional segmentation parameters.
// Returns a slice of segmented strings.
func (s *apiSegment) Segment(text string, args ...string) []string {
	return s.request(text, ``, args...)
}

// request sends a text segmentation request to the API endpoint.
// It takes the input text, segmentation mode, and optional arguments,
// then sends them to the configured API URL with authentication if an API key is set.
// Returns the segmented result as a string slice or nil on error.
// Errors are logged internally and not returned.
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

// SegmentBy splits the input text into segments using the specified mode and optional arguments.
// Returns a slice of segmented strings.
func (s *apiSegment) SegmentBy(text string, mode string, args ...string) []string {
	return s.request(text, mode, args...)
}

// Close implements the io.Closer interface for apiSegment. It's a no-op method that always returns nil.
func (s *apiSegment) Close() error {
	return nil
}
