package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBulkURL(t *testing.T) {
	bgeAPIURL := GetBgeAPIUrl()
	assert.Equal(t, defaultBgeAPIUrl, bgeAPIURL)
}
