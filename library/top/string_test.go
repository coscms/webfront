package top

import (
	"html/template"
	"net/url"
	"testing"
	"time"

	"github.com/coscms/webcore/library/config"
	"github.com/stretchr/testify/assert"
)

func TestOutputContent(t *testing.T) {
	v := OutputContent(`1 1 2 3`, `text`)
	assert.Equal(t, template.HTML(`1 1 2 3`), v)

	v = OutputContent(`1 &nbsp; 1 2 3`, `text`)
	assert.Equal(t, template.HTML(`1 &amp;nbsp; 1 2 3`), v)

	v = OutputContent("1 \r\n 1 2 3", `text`)
	assert.Equal(t, template.HTML(`1 <br /> 1 2 3`), v)
}

func TestMakeEncodedURL(t *testing.T) {
	config.NewConfig().SetDefaults().AsDefault()
	expiry := time.Now().Add(time.Hour * 24).Unix()
	rawURL := `https://coscms.com/download?a=测试&b=1#c`
	v, err := MakeEncodedURL(rawURL, expiry)
	assert.NoError(t, err)
	assert.True(t, len(v) > 0)
	t.Logf(`encodedURL: %s`, v)
	u, err := url.Parse(v)
	assert.NoError(t, err)
	_rawURL, _expiry, err := ParseEncodedURL(u.Query().Get(`url`))
	assert.NoError(t, err)
	assert.Equal(t, _rawURL, rawURL)
	assert.Equal(t, _expiry, expiry)
}

func TestTrimOverflowTextSlice(t *testing.T) {
	v := TrimOverflowTextSlice([]string{`100000`, `20000`, `300000`}, 10, `,`)
	assert.Equal(t, []string{`100000`}, v)
	v = TrimOverflowTextSlice([]string{`100000`, `20000`, `300000`}, 15, `,`)
	assert.Equal(t, []string{`100000`, `20000`}, v)
	v = TrimOverflowTextSlice([]string{`100000`, `20000`, `300000`}, 5, `,`)
	assert.Equal(t, []string{}, v)
}

func TestTrimOverflowText(t *testing.T) {
	v := TrimOverflowText(`100000,20000,300000`, 10, `,`)
	assert.Equal(t, `100000`, v)
	v = TrimOverflowText(`100000,20000,300000`, 15, `,`)
	assert.Equal(t, `100000,20000`, v)
	v = TrimOverflowText(`100000,20000,300000`, 7, `,`)
	assert.Equal(t, `100000`, v)
	v = TrimOverflowText(`100000,20000,300000`, 6, `,`)
	assert.Equal(t, `100000`, v)
	v = TrimOverflowText(`100000,20000,300000`, 5, `,`)
	assert.Equal(t, ``, v)
	v = TrimOverflowText(`哈哈哈`, 5)
	assert.Equal(t, `哈`, v)
	v = TrimOverflowText(`哈12345`, 5)
	assert.Equal(t, `哈12`, v)
	v = TrimOverflowText(`12哈345`, 5)
	assert.Equal(t, `12哈`, v)
	v = TrimOverflowText(`1😊哈345`, 5)
	assert.Equal(t, `1😊`, v)
	v = TrimOverflowText(`1哈哈345`, 5)
	assert.Equal(t, `1哈`, v)
}
