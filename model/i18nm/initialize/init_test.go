package initialize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialized(t *testing.T) {
	v := isInitialized()
	assert.False(t, v)
	v = isInitialized()
	assert.True(t, v)
}
