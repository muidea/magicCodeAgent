package util

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathJoin(t *testing.T) {
	subItem := []string{"a", "b", "c"}
	reallyVal := path.Join(subItem...)
	assert.Equal(t, reallyVal, "a/b/c")

	subItem = []string{"a", "", "c"}
	reallyVal = path.Join(subItem...)
	assert.Equal(t, reallyVal, "a/c")
}
