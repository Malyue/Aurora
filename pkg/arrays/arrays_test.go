package arrays

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsContain(t *testing.T) {
	is := IsContain([]string{"a", "b", "c"}, "a")
	assert.Equal(t, is, true)
}
