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

package jiebago

import (
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/admpub/log"
	"github.com/wangbin/jiebago"
	"github.com/webx-top/echo"

	. "github.com/coscms/webfront/library/search/segment"
)

// init registers the jiebago segmenter with the default New constructor
func init() {
	Register(`jiebago`, New)
}

// New creates and returns a new Jieba segmenter instance with default dictionary path.
// The default dictionary is located at 'data/sego/dict.txt' relative to the working directory.
func New() Segment {
	return &Jieba{
		segmenter:   &jiebago.Segmenter{},
		defaultDict: filepath.Join(echo.Wd(), `data`, `sego/dict.txt`),
	}
}

type Jieba struct {
	segmenter   *jiebago.Segmenter
	defaultDict string
	dictLoaded  atomic.Bool
	once        sync.Once
}

// LoadDict loads dictionary files for Jieba segmenter. Multiple dictionary files can be specified
// separated by commas. The first file is loaded as main dictionary, subsequent files are loaded
// as user dictionaries. Returns nil after loading all dictionaries.
func (s *Jieba) LoadDict(dictFile string, dictType ...string) error {
	for index, file := range strings.Split(dictFile, `,`) { //多个字典文件用半角“,”逗号分隔
		var err error
		if index == 0 {
			log.Debug(`[jiebago]Load dictionary:`, file)
			err = s.segmenter.LoadDictionary(file)
		} else {
			log.Debug(`[jiebago]Load user dictionary:`, file)
			err = s.segmenter.LoadUserDictionary(file)
		}
		if err != nil {
			log.Error(`[jiebago]LoadDict:`, err)
		}
	}
	s.dictLoaded.Store(true)
	return nil
}

// initDict initializes the Jieba segmenter by loading the default dictionary.
func (s *Jieba) initDict() {
	s.LoadDict(s.defaultDict, `default`)
}

// Segment splits the input text into words using Jieba segmentation.
// It performs the following steps:
// - Initializes the dictionary if not loaded
// - Uses precise mode segmentation
// - Filters out empty words and applies word filtering
// - Removes stop words from the result
// Returns a slice of segmented and filtered words.
func (s *Jieba) Segment(text string, args ...string) []string {
	if !s.dictLoaded.Load() {
		s.once.Do(s.initDict)
	}
	var (
		words = []string{}
		ch    <-chan string //精确模式
	)
	ch = s.segmenter.Cut(text, false)

	for word := range ch {
		word = strings.TrimSpace(word)
		if len(word) > 0 && DoFilter(word) {
			words = append(words, word)
		}
	}
	words = CleanStopWordsFromSlice(words)
	return words
}

// SegmentBy segments the input text using the specified mode.
// Supported modes: "all" (full mode), "new" (new word recognition),
// "search" (search engine mode), or default (accurate mode).
// Returns a slice of filtered and cleaned words after segmentation.
// The words are filtered by DoFilter and cleaned from stop words.
func (s *Jieba) SegmentBy(text string, mode string, args ...string) []string {
	if !s.dictLoaded.Load() {
		s.once.Do(s.initDict)
	}
	var (
		words = []string{}
		ch    <-chan string //精确模式
	)
	switch mode {
	case `all`:
		//log.Println(`all mode:`, text)
		ch = s.segmenter.CutAll(text)
	case `new`: //新词识别
		//log.Println(`new mode:`, text)
		ch = s.segmenter.Cut(text, true)
	case `search`: //搜索引擎模式
		//log.Println(`search mode:`, text)
		ch = s.segmenter.CutForSearch(text, true)

	//TODO
	//case `tag`: //词性标注
	//case `keywords`: //关键词提取

	default: //精确模式
		ch = s.segmenter.Cut(text, false)
	}
	for word := range ch {
		word = strings.TrimSpace(word)
		if len(word) > 0 && DoFilter(word) {
			words = append(words, word)
		}
	}
	words = CleanStopWordsFromSlice(words)
	return words
}

// Close implements io.Closer interface for Jieba segmenter, currently a no-op.
func (s *Jieba) Close() error {
	return nil
}
