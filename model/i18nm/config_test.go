package i18nm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopySlices(t *testing.T) {
	sl1 := []string{`A`, `B`, `C`}
	sl2 := []string{`D`, `E`}
	sl := make([]string, len(sl1)+len(sl2))

	n := copy(sl, sl1)
	copy(sl[n:], sl2)
	assert.Equal(t, []string{`A`, `B`, `C`, `D`, `E`}, sl)
}
