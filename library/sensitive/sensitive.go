package sensitive

import (
	"errors"
	"io"
	"strings"

	"github.com/admpub/log"
	syncOnce "github.com/admpub/once"
	"github.com/admpub/sensitive"
	"github.com/webx-top/db"
	"github.com/webx-top/echo/defaults"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/extend"
	"github.com/coscms/webfront/dbschema"
)

func init() {
	extend.Register(`sensitive`, func() interface{} {
		return &Config{}
	})
}

type Config struct {
	DictPath string `json:"dictPath"`
	DictURL  string `json:"dictURL"`
}

var (
	defaultFilter *sensitive.Filter
	once          syncOnce.Once
)

func IsInitialized() bool {
	return defaultFilter != nil
}

func initDefaultFilter() {
	defaultFilter = sensitive.New()
	ctx := defaults.NewMockContext()
	m := dbschema.NewOfficialCommonSensitive(ctx)
	m.ListByOffset(nil, nil, 0, -1, db.Cond{`disabled`: `N`})
	var noises []string
	for _, f := range m.Objects() {
		if len(f.Words) == 0 {
			continue
		}
		switch f.Type {
		case `bad`:
			defaultFilter.AddWord(f.Words)
		case `noise`:
			noises = append(noises, f.Words)
		}
	}
	if len(noises) > 0 {
		defaultFilter.UpdateNoisePattern(strings.Join(noises, `|`))
	}
	cfg, ok := config.FromFile().Extend.Get(`sensitive`).(*Config)
	if !ok {
		return
	}
	if len(cfg.DictPath) > 0 {
		log.Info(`[sensitive] load word dict: ` + cfg.DictPath)
		if err := defaultFilter.LoadWordDict(cfg.DictPath); err != nil {
			log.Error(err)
		}
	}
	if len(cfg.DictURL) > 0 {
		log.Info(`[sensitive] load word dict: ` + cfg.DictURL)
		if err := defaultFilter.LoadNetWordDict(cfg.DictURL); err != nil {
			log.Error(err)
		}
	}
}

func Default() *sensitive.Filter {
	once.Do(initDefaultFilter)
	return defaultFilter
}

// Reset 重置
func Reset() {
	once.Reset()
}

// UpdateNoisePattern 更新去噪模式
func UpdateNoisePattern(pattern string) {
	Default().UpdateNoisePattern(pattern)
}

// LoadWordDict 加载敏感词字典
func LoadWordDict(path string) error {
	return Default().LoadWordDict(path)
}

// LoadNetWordDict 加载网络敏感词字典
func LoadNetWordDict(url string) error {
	return Default().LoadNetWordDict(url)
}

// Load common method to add words
func Load(rd io.Reader) error {
	return Default().Load(rd)
}

// AddWord 添加敏感词
func AddWord(words ...string) {
	Default().AddWord(words...)
}

// DelWord 删除敏感词
func DelWord(words ...string) {
	Default().DelWord(words...)
}

// Filter 过滤敏感词
func Filter(text string) string {
	return Default().Filter(text)
}

// Replace 和谐敏感词
func Replace(text string, repl rune) string {
	return Default().Replace(text, repl)
}

// FindIn 检测敏感词
func FindIn(text string) (bool, string) {
	return Default().FindIn(text)
}

// FindAll 找到所有匹配词
func FindAll(text string) []string {
	return Default().FindAll(text)
}

// Validate 检测字符串是否合法
func Validate(text string) (bool, string) {
	return Default().Validate(text)
}

// RemoveNoise 去除空格等噪音
func RemoveNoise(text string) string {
	return Default().RemoveNoise(text)
}

var ErrSensitive = errors.New("发送失败: 不能包含违禁词。(failed to send: cannot contain prohibited words)")
