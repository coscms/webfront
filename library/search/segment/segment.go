/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package segment

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/once"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/echo"
	"go.uber.org/atomic"
)

type Filter func(string) bool

var (
	DefaultEngine  = atomic.NewString(`sego`) // gojieba / sego / jiebago
	stopWords      []string
	stopWordsMap   map[string]bool
	Filters        []Filter
	defaultSegment atomic.Value
	onceSegment    once.Once
	onceStopword   once.Once
)

// initDefaultSegment initializes the default segmentation engine by loading it from the DefaultEngine configuration
// and storing it in the defaultSegment variable. This function is called during package initialization.
func initDefaultSegment() {
	log.Debug("[segment]Default engine:", DefaultEngine)
	defaultSegment.Store(Get(DefaultEngine.Load()))
}

// IsInitialized reports whether the default segment has been initialized.
func IsInitialized() bool {
	return defaultSegment.Load() != nil
}

// Default returns the default Segment instance, initializing it if necessary.
// The initialization is thread-safe and will only occur once.
func Default() Segment {
	onceSegment.Do(initDefaultSegment)
	return defaultSegment.Load().(Segment)
}

// ResetSegment closes the current segment if it exists and resets the segment initialization state.
// This allows for a new segment to be created on next use.
func ResetSegment() {
	seg := defaultSegment.Load()
	if seg != nil {
		seg.(Segment).Close()
	}
	onceSegment.Reset()
}

// ResetStopwords resets the stopwords initialization state, allowing stopwords to be reloaded on next use.
func ResetStopwords() {
	onceStopword.Reset()
}

// StopWords returns the list of stop words that are loaded once and cached.
// The stop words are initialized on first call and returned from cache on subsequent calls.
func StopWords() []string {
	onceStopword.Do(initLoadStopWordsDict)
	return stopWords
}

// LoadStopWordsDict loads stop words from the specified file into memory.
// If rebuild is true, it clears existing stop words before loading new ones.
// Non-empty lines in the file are treated as stop words after trimming whitespace.
func LoadStopWordsDict(stopWordsFile string, args ...bool) {
	var rebuild bool
	if len(args) > 0 {
		rebuild = args[0]
	}
	b, err := os.ReadFile(stopWordsFile)
	if err != nil {
		log.Debug(stopWordsFile+`:`, err)
		return
	}
	words := strings.Split(strings.TrimSpace(string(b)), "\n")
	if rebuild {
		stopWords = []string{}
		stopWordsMap = nil
	}
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 0 {
			stopWords = append(stopWords, word)
		}
	}
}

// initLoadStopWordsDict initializes the stop words dictionary by searching for stopwords.txt
// in predefined locations (data/sego/stopwords.txt and GOPATH/src/github.com/coscms/webfront/library/search/segment/stopwords.txt).
// It loads the first found stop words file and stops searching. Debug logs are printed for files that cannot be accessed.
func initLoadStopWordsDict() {
	var stopWordsFiles []string
	stopWordsFiles = append(stopWordsFiles, filepath.Join(echo.Wd(), `data`, `sego/stopwords.txt`))
	goPath := os.Getenv(`GOPATH`)
	if len(goPath) > 0 {
		stopWordsFiles = append(stopWordsFiles, filepath.Join(goPath, `src`, `github.com/coscms/webfront/library/search/segment/stopwords.txt`))
	}
	for _, stopWordsFile := range stopWordsFiles {
		_, err := os.Stat(stopWordsFile)
		if err == nil {
			LoadStopWordsDict(stopWordsFile)
			return
		}
		log.Debug(stopWordsFile+`:`, err)
	}
}

// CleanStopWords removes all stop words from the input string and replaces them with spaces.
// It returns the cleaned string with stop words removed.
func CleanStopWords(v string) string {
	for _, word := range StopWords() {
		v = strings.Replace(v, word, ` `, -1)
	}
	return v
}

// CleanStopWordsFromSlice removes stop words from the input string slice and returns a new slice
// containing only non-stop words. Stop words are loaded from StopWords() if not already initialized.
func CleanStopWordsFromSlice(v []string) (r []string) {
	if stopWordsMap == nil {
		stopWordsMap = make(map[string]bool)
		for _, word := range StopWords() {
			stopWordsMap[word] = true
		}
	}
	r = make([]string, 0)
	for _, w := range v {
		if _, ok := stopWordsMap[w]; !ok {
			r = append(r, w)
		}
	}
	return r
}

// DoFilter checks if the given string passes all registered filters.
// Returns true if all filters accept the string, false otherwise.
func DoFilter(v string) bool {
	for _, f := range Filters {
		if !f(v) {
			return false
		}
	}
	return true
}

// AddFilter appends a new filter to the global Filters collection.
func AddFilter(filter Filter) {
	Filters = append(Filters, filter)
}

// ApplySegmentConfig applies segment configuration from the given config object.
// It initializes or updates the segment engine based on the configuration.
// If the engine is 'api', it will register a new API segment with the provided URL and key.
// The function handles engine switching and cleanup of previous segment instances.
func ApplySegmentConfig(c *config.Config) {
	segmentCfg := c.Extend.GetStore(`segment`)
	segmentEngine := segmentCfg.String(`engine`)
	if len(segmentEngine) == 0 {
		return
	}
	if DefaultEngine.Load() != segmentEngine {
		seg := defaultSegment.Load()
		if seg != nil {
			seg.(Segment).Close()
		}
		DefaultEngine.Store(segmentEngine)
	}
	switch segmentEngine {
	case `api`:
		segmentApiURL := segmentCfg.String(`apiURL`)
		segmentApiKey := segmentCfg.String(`apiKey`)
		seg := Get(segmentEngine)
		reg := true
		if !IsNop(seg) {
			if apiSeg, ok := seg.(*apiSegment); ok {
				reg = apiSeg.apiKey != segmentApiKey || apiSeg.apiURL != segmentApiURL
			}
		}
		if reg {
			a := NewAPI(segmentApiURL, segmentApiKey)
			Register(segmentEngine, func() Segment {
				return a
			})
			ResetSegment()
		}
	default:
	}
}

// Segment interface
type Segment interface {
	//载入词典（词典路径，词典类型）
	LoadDict(string, ...string) error

	//分词（文本，词性）
	Segment(string, ...string) []string

	//分词（文本，分词模式，词性）
	SegmentBy(string, string, ...string) []string

	//关闭或释放资源
	Close() error
}

type nopSegment struct {
}

// LoadDict implements the Segment interface but does nothing, always returning nil.
// dictFile is the dictionary file path to load (ignored).
// dictType specifies optional dictionary type (ignored).
func (s *nopSegment) LoadDict(dictFile string, dictType ...string) error {
	return nil
}

// Segment implements the Segmenter interface but returns an empty slice for any input.
// This is a no-op implementation that performs no actual segmentation.
func (s *nopSegment) Segment(text string, args ...string) []string {
	return []string{}
}

// SegmentBy segments the input text according to the specified mode and returns an empty slice.
// This is a no-op implementation that always returns an empty result.
func (s *nopSegment) SegmentBy(text string, mode string, args ...string) []string {
	return []string{}
}

// Close implements the Segment interface with a no-op operation.
func (s *nopSegment) Close() error {
	return nil
}
