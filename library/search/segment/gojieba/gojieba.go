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

package gojieba

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/yanyiwu/gojieba"

	. "github.com/coscms/webfront/library/search/segment"
)

// init registers the jieba segmenter with the "gojieba" name using the New function
func init() {
	Register(`gojieba`, New)
}

// New creates and returns a new Jieba segmenter instance with default dictionary directory.
func New() Segment {
	return &Jieba{
		defaultDictDir: filepath.Join(echo.Wd(), `data`, `sego`),
	}
}

type Jieba struct {
	segmenter      *gojieba.Jieba
	defaultDictDir string
	dictLoaded     atomic.Bool
	once           sync.Once
}

// SetDictDir sets the dictionary directory path and initializes all dictionary file paths
// including main dictionary, HMM model, user dictionary, IDF and stop words files.
func SetDictDir(dictDir string) {
	gojieba.DICT_DIR = dictDir
	gojieba.DICT_PATH = filepath.Join(gojieba.DICT_DIR, "dict.txt")
	gojieba.HMM_PATH = filepath.Join(gojieba.DICT_DIR, "hmm.model.txt")
	gojieba.USER_DICT_PATH = filepath.Join(gojieba.DICT_DIR, "user.dict.txt")
	gojieba.IDF_PATH = filepath.Join(gojieba.DICT_DIR, "idf.txt")
	gojieba.STOP_WORDS_PATH = filepath.Join(gojieba.DICT_DIR, "stopwords.txt")
}

// LoadDict loads Jieba dictionary from specified file or directory.
// If dictFile is a file path, it will use the containing directory as dictionary location.
// Returns error if dictionary file/directory cannot be accessed.
// Automatically closes previous segmenter and initializes new one on success.
func (s *Jieba) LoadDict(dictFile string, dictType ...string) error {
	dictFile = strings.SplitN(dictFile, `,`, 2)[0]
	log.Debug(`[gojieba]LoadDict:`, dictFile)
	fi, err := os.Stat(dictFile)
	if err != nil {
		log.Error(`[gojieba]LoadDict:`, err)
		return err
	}
	var dictDir string
	if !fi.IsDir() {
		dictDir = filepath.Dir(dictFile)
	} else {
		dictDir = dictFile
	}
	SetDictDir(dictDir)
	s.Close()
	s.segmenter = gojieba.NewJieba()
	s.dictLoaded.Store(true)
	return nil
}

func (s *Jieba) initDict() {
	s.LoadDict(s.defaultDictDir, `default`)
}

func (s *Jieba) Segment(text string, args ...string) []string {
	if !s.dictLoaded.Load() {
		s.once.Do(s.initDict)
	}
	var (
		words     []string
		rets      []string //精确模式
		wordTypes []string //获取指定类型的词语,如仅仅获取名词则为n
	)
	if len(args) > 0 {
		wordTypes = args
	}

	if len(wordTypes) > 0 {

		words = s.segmenter.Tag(text)

		for _, word := range words {
			p := strings.LastIndex(word, `/`)
			if p < 0 {
				continue
			}
			_word := word
			word = _word[0:p]
			typ := _word[p+1:]
			for _, wt := range wordTypes {
				if !strings.Contains(typ, wt) {
					continue
				}
				word = strings.TrimSpace(word)
				if len(word) > 0 && DoFilter(word) {
					rets = append(rets, word)
				}
			}
		}
	} else {

		words = s.segmenter.Cut(text, false)

		for _, word := range words {
			word = strings.TrimSpace(word)
			if len(word) > 0 && DoFilter(word) {
				rets = append(rets, word)
			}
		}
	}
	rets = CleanStopWordsFromSlice(rets)
	return rets
}

// SegmentBy performs Chinese text segmentation using Jieba with specified mode.
// Available modes:
//   - "all": full mode that segments all possible words
//   - "new": new word recognition mode
//   - "search": search engine mode optimized for indexing
//   - "tag": POS tagging mode that can filter by word types (n, v, etc.)
//   - "keywords": extracts top N keywords (default 50)
//   - default: precise mode for accurate segmentation
//
// Args can specify word types for filtering (in "tag" mode) or top N count (in "keywords" mode).
// Returns filtered and cleaned word segments based on mode.
func (s *Jieba) SegmentBy(text string, mode string, args ...string) []string {
	if !s.dictLoaded.Load() {
		s.once.Do(s.initDict)
	}
	var (
		words     []string
		rets      []string //精确模式
		wordTypes []string //获取指定类型的词语,如仅仅获取名词则为n
	)
	if len(args) > 0 {
		wordTypes = args
	}

	if len(wordTypes) > 0 {
		mode = `tag`
	}
	switch mode {
	case `all`:
		//log.Println(`all mode:`, text)
		words = s.segmenter.CutAll(text)
	case `new`: //新词识别
		//log.Println(`new mode:`, text)
		words = s.segmenter.Cut(text, true)
	case `search`: //搜索引擎模式
		//log.Println(`search mode:`, text)
		words = s.segmenter.CutForSearch(text, true)
	case `tag`: //词性标注
		words = s.segmenter.Tag(text)
		if len(wordTypes) > 0 {
			for _, word := range words {
				p := strings.LastIndex(word, `/`)
				if p < 0 {
					continue
				}
				_word := word
				word = _word[0:p]
				typ := _word[p+1:]
				for _, wt := range wordTypes {
					if !strings.Contains(typ, wt) {
						continue
					}
					word = strings.TrimSpace(word)
					if len(word) > 0 && DoFilter(word) {
						rets = append(rets, word)
					}
				}
			}
			rets = CleanStopWordsFromSlice(rets)
			return rets
		}
	case `keywords`: //关键词提取
		topN := 50
		if len(args) > 0 {
			if n, err := strconv.Atoi(args[0]); err == nil {
				topN = n
			}
		}
		words = s.segmenter.Extract(text, topN)
	default: //精确模式
		words = s.segmenter.Cut(text, false)
	}
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 0 && DoFilter(word) {
			rets = append(rets, word)
		}
	}
	rets = CleanStopWordsFromSlice(rets)
	return rets
}

// Close releases resources associated with the Jieba segmenter.
// It safely frees the underlying segmenter if it exists.
// Returns nil on success (always succeeds).
func (s *Jieba) Close() error {
	if s.segmenter != nil {
		s.segmenter.Free()
	}
	return nil
}
