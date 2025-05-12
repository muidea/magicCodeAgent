package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuleName(t *testing.T) {
	moduleName := "abc123"

	assert.True(t, projectPattern.MatchString(moduleName))

	moduleName = "service/abc123"
	assert.True(t, projectPattern.MatchString(moduleName))

	moduleName = "web/abc123"
	assert.True(t, projectPattern.MatchString(moduleName))

	moduleName = "service/web/abc123"
	assert.False(t, projectPattern.MatchString(moduleName))

	moduleName = "afdf/abc123"
	assert.False(t, projectPattern.MatchString(moduleName))

	moduleName = "123abc/abc123"
	assert.False(t, projectPattern.MatchString(moduleName))
}
