package minify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	v := `
  <link combine href="{{AssetsURL}}/css/form/input-webx/webx.min.css?t={{BuildTime}}" rel="stylesheet" charset="utf-8">
  <link combine href="{{AssetsXURL}}/css/custom.min.css?t={{BuildTime}}" rel="stylesheet" charset="utf-8">
<script combine src="{{AssetsURL}}/js/compatible.min.js?t={{BuildTime}}" type="text/javascript"></script>
<script combine src="{{AssetsURL}}/js/jquery3.6.min.js?t={{BuildTime}}"></script>
`
	r := Merge([]byte(v), nil)
	assert.Equal(t, `<link href="{{AssetsURL}}/combined/0/d41d8cd98f00b204e9800998ecf8427e.min.css" rel="stylesheet" /><script src="{{AssetsURL}}/combined/0/d41d8cd98f00b204e9800998ecf8427e.min.js"></script>`, string(r))
}
