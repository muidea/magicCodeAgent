package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBulkURL(t *testing.T) {
	nlpAPIURL := GetNlpEmbedAPIUrl()
	assert.Equal(t, defaultNlpAPIUrl, nlpAPIURL)
}
