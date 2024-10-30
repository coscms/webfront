package segment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	a := NewAPI(`http://dev.coscms.com:8181/segment`, `1.0.1`)
	r := a.Segment(`世界上的海洋`)
	assert.Equal(t, []string{"世界", "海洋"}, r)
}
