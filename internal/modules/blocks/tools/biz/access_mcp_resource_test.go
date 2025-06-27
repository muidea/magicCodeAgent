package biz

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessMcpResource(t *testing.T) {
	// Test case 1: Normal scenario
	serverName := "weather-server"
	uri := "weather://san-francisco/current"
	expectedResult := fmt.Sprintf("Accessed resource '%s' on server '%s'", uri, serverName)

	result, err := AccessMcpResource(serverName, uri)

	// Verify result
	assert.NoError(t, err, "should not have error")
	assert.Equal(t, expectedResult, result, "result should match expectation")

	// Test case 2: Empty server name
	emptyServerName := ""
	uri2 := "resource://path"
	expectedResult2 := fmt.Sprintf("Accessed resource '%s' on server '%s'", uri2, emptyServerName)

	result2, err2 := AccessMcpResource(emptyServerName, uri2)

	// Verify result
	assert.NoError(t, err2, "empty server name should not cause error")
	assert.Equal(t, expectedResult2, result2, "result should match expectation")

	// Test case 3: Empty URI
	serverName3 := "test-server"
	emptyUri := ""
	expectedResult3 := fmt.Sprintf("Accessed resource '%s' on server '%s'", emptyUri, serverName3)

	result3, err3 := AccessMcpResource(serverName3, emptyUri)

	// Verify result
	assert.NoError(t, err3, "empty URI should not cause error")
	assert.Equal(t, expectedResult3, result3, "result should match expectation")

	// Test case 4: Special characters
	serverName4 := "special-server!@#"
	uri4 := "resource://path?query=value&special=!@#"
	expectedResult4 := fmt.Sprintf("Accessed resource '%s' on server '%s'", uri4, serverName4)

	result4, err4 := AccessMcpResource(serverName4, uri4)

	// Verify result
	assert.NoError(t, err4, "special characters should not cause error")
	assert.Equal(t, expectedResult4, result4, "result should match expectation")
}
