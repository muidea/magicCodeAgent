package biz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkspaceResourceURI(t *testing.T) {
	ptr := &MCP{}

	idVal := "test"
	rawVal := ptr.workspaceFileResourceContentURI(idVal)
	nID, nErr := ptr.extractID(rawVal)
	assert.Nil(t, nErr)
	assert.Equal(t, idVal, nID)
}
