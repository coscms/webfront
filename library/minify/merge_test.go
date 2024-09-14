package minify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	v := `
  <link combine href="{{AssetsURL}}/css/form/input-webx/webx.min.css?t={{BuildTime}}" rel="stylesheet" charset="utf-8">
  <link combine href="{{AssetsXURL}}/css/custom.min.css?t={{BuildTime}}" rel="stylesheet" charset="utf-8">
  <link combine="a" href="{{AssetsURL}}/css/form/input-webx/webx-a.min.css?t={{BuildTime}}" rel="stylesheet" charset="utf-8">
  <link combine="a" href="{{AssetsXURL}}/css/custom-a.min.css?t={{BuildTime}}" rel="stylesheet" charset="utf-8">
<script combine src="{{AssetsURL}}/js/compatible.min.js?t={{BuildTime}}" type="text/javascript"></script>
<script combine src="{{AssetsURL}}/js/jquery3.6.min.js?t={{BuildTime}}"></script>
`
	r := Merge([]byte(v), true)
	assert.Equal(t, `<link href="{{AssetsURL}}/combined/0/0-d35057682c082b6e80edbf5e95aa41ff.min.css" rel="stylesheet" /><link href="{{AssetsURL}}/combined/0/a-79947a278f6e37e603ba92b6deafdf19.min.css" rel="stylesheet" /><script src="{{AssetsURL}}/combined/0/0-2c506cfe58fb7ec6e13b57de1a6492c6.min.js"></script>`, string(r))
}

func TestResolveURLPath(t *testing.T) {
	r := resolveURLPath(`/a/b/c/d`, `/a/b/c/d/e/f`)
	assert.Equal(t, `../../`, r)
	r = resolveURLPath(`/a/b/c/d`, `/a/b`)
	assert.Equal(t, `c/d`, r)

	r = resolveURLPath(`/a/b/c`, `/d/e/f`)
	assert.Equal(t, `../../../a/b/c`, r)

	r = replaceCSSImportURL(`url('../img/a.jpg')`, `/public/assets/backend/css/css.css`, `/public/assets/backend/combined`)
	assert.Equal(t, `url(../css/img/a.jpg)`, r)
}
