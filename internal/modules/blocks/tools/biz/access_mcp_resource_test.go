package biz

import (
	"fmt"
	"testing"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
	"github.com/stretchr/testify/assert"
)

func TestAccessMcpResource(t *testing.T) {
	// 创建一个Tool实例用于测试
	eventHub := event.NewHub(100)
	backgroundRoutine := task.NewBackgroundRoutine(100)
	tool := New(eventHub, backgroundRoutine, nil)

	// 测试用例1：正常情况
	serverName := "weather-server"
	uri := "weather://san-francisco/current"
	expectedResult := fmt.Sprintf("Accessed resource '%s' on server '%s'", uri, serverName)

	result, err := tool.AccessMcpResource(serverName, uri)

	// 验证结果
	assert.NoError(t, err, "应该没有错误")
	assert.Equal(t, expectedResult, result, "返回结果应该符合预期")

	// 测试用例2：空服务器名
	emptyServerName := ""
	uri2 := "resource://path"
	expectedResult2 := fmt.Sprintf("Accessed resource '%s' on server '%s'", uri2, emptyServerName)

	result2, err2 := tool.AccessMcpResource(emptyServerName, uri2)

	// 验证结果
	assert.NoError(t, err2, "空服务器名应该不会导致错误")
	assert.Equal(t, expectedResult2, result2, "返回结果应该符合预期")

	// 测试用例3：空URI
	serverName3 := "test-server"
	emptyUri := ""
	expectedResult3 := fmt.Sprintf("Accessed resource '%s' on server '%s'", emptyUri, serverName3)

	result3, err3 := tool.AccessMcpResource(serverName3, emptyUri)

	// 验证结果
	assert.NoError(t, err3, "空URI应该不会导致错误")
	assert.Equal(t, expectedResult3, result3, "返回结果应该符合预期")

	// 测试用例4：特殊字符
	serverName4 := "special-server!@#"
	uri4 := "resource://path?query=value&special=!@#"
	expectedResult4 := fmt.Sprintf("Accessed resource '%s' on server '%s'", uri4, serverName4)

	result4, err4 := tool.AccessMcpResource(serverName4, uri4)

	// 验证结果
	assert.NoError(t, err4, "特殊字符应该不会导致错误")
	assert.Equal(t, expectedResult4, result4, "返回结果应该符合预期")
}
