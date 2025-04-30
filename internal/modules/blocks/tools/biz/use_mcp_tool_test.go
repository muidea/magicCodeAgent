package biz

import (
	"encoding/json"
	"testing"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
	"github.com/stretchr/testify/assert"
)

func TestUseMcpTool(t *testing.T) {
	// 创建一个Tool实例用于测试
	eventHub := event.NewHub(100)
	backgroundRoutine := task.NewBackgroundRoutine(100)
	tool := New(eventHub, backgroundRoutine)

	// 测试用例1：简单参数
	serverName1 := "weather-server"
	toolName1 := "get_forecast"
	arguments1 := json.RawMessage(`{"city":"San Francisco","days":5}`)
	expectedResult1 := "Executed tool 'get_forecast' on server 'weather-server' with arguments: {\"city\":\"San Francisco\",\"days\":5}"

	result1, err := tool.UseMcpTool(serverName1, toolName1, arguments1)

	// 验证结果
	assert.NoError(t, err, "使用MCP工具应该成功")
	assert.Equal(t, expectedResult1, result1, "返回结果应该符合预期")

	// 测试用例2：空参数
	serverName2 := "test-server"
	toolName2 := "empty-args-tool"
	arguments2 := json.RawMessage(`{}`)
	expectedResult2 := "Executed tool 'empty-args-tool' on server 'test-server' with arguments: {}"

	result2, err := tool.UseMcpTool(serverName2, toolName2, arguments2)

	// 验证结果
	assert.NoError(t, err, "空参数应该不会导致错误")
	assert.Equal(t, expectedResult2, result2, "返回结果应该符合预期")

	// 测试用例3：复杂嵌套参数
	serverName3 := "data-server"
	toolName3 := "process-data"
	arguments3 := json.RawMessage(`{"data":{"items":[1,2,3],"config":{"process":true,"format":"json"}}}`)
	expectedResult3 := "Executed tool 'process-data' on server 'data-server' with arguments: {\"data\":{\"items\":[1,2,3],\"config\":{\"process\":true,\"format\":\"json\"}}}"

	result3, err := tool.UseMcpTool(serverName3, toolName3, arguments3)

	// 验证结果
	assert.NoError(t, err, "复杂嵌套参数应该不会导致错误")
	assert.Equal(t, expectedResult3, result3, "返回结果应该符合预期")

	// 测试用例4：空服务器名
	serverName4 := ""
	toolName4 := "some-tool"
	arguments4 := json.RawMessage(`{"param":"value"}`)
	expectedResult4 := "Executed tool 'some-tool' on server '' with arguments: {\"param\":\"value\"}"

	result4, err := tool.UseMcpTool(serverName4, toolName4, arguments4)

	// 验证结果
	assert.NoError(t, err, "空服务器名应该不会导致错误")
	assert.Equal(t, expectedResult4, result4, "返回结果应该符合预期")

	// 测试用例5：空工具名
	serverName5 := "some-server"
	toolName5 := ""
	arguments5 := json.RawMessage(`{"param":"value"}`)
	expectedResult5 := "Executed tool '' on server 'some-server' with arguments: {\"param\":\"value\"}"

	result5, err := tool.UseMcpTool(serverName5, toolName5, arguments5)

	// 验证结果
	assert.NoError(t, err, "空工具名应该不会导致错误")
	assert.Equal(t, expectedResult5, result5, "返回结果应该符合预期")

	// 测试用例6：nil参数（应该被视为空JSON对象）
	serverName6 := "nil-server"
	toolName6 := "nil-tool"
	var arguments6 json.RawMessage
	expectedResult6 := "Executed tool 'nil-tool' on server 'nil-server' with arguments: "

	result6, err := tool.UseMcpTool(serverName6, toolName6, arguments6)

	// 验证结果
	assert.NoError(t, err, "nil参数应该不会导致错误")
	assert.Equal(t, expectedResult6, result6, "返回结果应该符合预期")
}
