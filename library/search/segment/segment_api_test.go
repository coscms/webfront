package segment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApi(t *testing.T) {
	tm, err := time.Parse(time.RFC1123, `Fri, 31 Dec 1999 23:59:59 GMT`)
	require.NoError(t, err)
	assert.Equal(t, time.Date(1999, time.December, 31, 23, 59, 59, 0, time.UTC), tm.UTC())

	a := NewAPI(`http://dev.coscms.com:8181/segment`, `1.0.1`)
	r := a.Segment(`世界上的海洋`)
	assert.Equal(t, []string{"世界", "海洋"}, r)
}
