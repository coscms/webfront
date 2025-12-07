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

package sego

import (
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/admpub/log"
	"github.com/huichen/sego"
	"github.com/webx-top/echo"

	. "github.com/coscms/webfront/library/search/segment"
)

// init registers the 'sego' segmenter with the New constructor function
func init() {
	Register(`sego`, New)
}

// New creates and returns a new Sego segmenter instance with default dictionary path.
// The default dictionary is located at "data/sego/dict.txt" relative to the working directory.
func New() Segment {
	return &Sego{
		segmenter:   &sego.Segmenter{},
		defaultDict: filepath.Join(echo.Wd(), `data`, `sego/dict.txt`),
	}
}

type Sego struct {
	segmenter   *sego.Segmenter
	defaultDict string
	dictLoaded  atomic.Bool
	once        sync.Once
}

// LoadDict loads a dictionary file for segmentation. Multiple dictionary files can be specified by comma-separated paths.
// dictFile specifies the path to the dictionary file(s)
// dictType optionally specifies the type of dictionary (not currently used)
// Returns nil on success or an error if loading fails
func (s *Sego) LoadDict(dictFile string, dictType ...string) error {
	log.Debug(`[sego]LoadDict:`, dictFile)
	s.segmenter.LoadDictionary(dictFile) //多个字典文件用半角“,”逗号分隔
	s.dictLoaded.Store(true)
	return nil
}

func (s *Sego) initDict() {
	s.LoadDict(s.defaultDict, `default`)
}

// Segment segments the input text into words with optional filtering by word types.
//
// text: the input string to be segmented
// args: optional word types to filter (e.g. "n" for nouns only)
// returns: a slice of segmented words after filtering and cleaning
func (s *Sego) Segment(text string, args ...string) []string {
	if !s.dictLoaded.Load() {
		s.once.Do(s.initDict)
	}
	segments := s.segmenter.Segment([]byte(text))
	var (
		words     []string
		wordTypes []string //获取指定类型的词语,如仅仅获取名词则为n
	)
	if len(args) > 0 {
		wordTypes = args
	}
	typeLength := len(wordTypes)
	for _, seg := range segments {
		//排除指定词性的词
		if typeLength > 0 {
			var ok bool
			for _, wt := range wordTypes {
				if seg.Token().Pos() == wt {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		content := seg.Token().Text()
		content = strings.Replace(content, `　`, ``, -1)
		content = strings.TrimSpace(content)
		if len(content) > 0 && DoFilter(content) {
			words = append(words, content)
		}
	}
	words = CleanStopWordsFromSlice(words)
	return words
}

// SegmentBy segments the input text using the specified mode and additional arguments.
// It delegates to the Segment method with the provided text and arguments.
func (s *Sego) SegmentBy(text string, mode string, args ...string) []string {
	return s.Segment(text, args...)
}

// Close implements io.Closer interface to clean up resources when the Sego instance is no longer needed.
func (s *Sego) Close() error {
	return nil
}
