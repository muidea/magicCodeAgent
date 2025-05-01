package biz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchInstructions(t *testing.T) {
	// 测试用例1：请求create_mcp_server指令
	task1 := "create_mcp_server"
	expectedResult1 := "Instructions for creating an MCP server..."

	result1, err1 := FetchInstructions(task1)

	// 验证结果
	assert.NoError(t, err1, "获取create_mcp_server指令应该成功")
	assert.Equal(t, expectedResult1, result1, "返回的指令应该符合预期")

	// 测试用例2：请求create_mode指令
	task2 := "create_mode"
	expectedResult2 := "Instructions for creating a mode..."

	result2, err2 := FetchInstructions(task2)

	// 验证结果
	assert.NoError(t, err2, "获取create_mode指令应该成功")
	assert.Equal(t, expectedResult2, result2, "返回的指令应该符合预期")

	// 测试用例3：请求未知指令
	task3 := "unknown_task"
	expectedResult3 := ""

	result3, err3 := FetchInstructions(task3)

	// 验证结果
	assert.NoError(t, err3, "获取未知指令应该不会返回错误")
	assert.Equal(t, expectedResult3, result3, "未知指令应该返回空字符串")

	// 测试用例4：空任务名
	task4 := ""
	expectedResult4 := ""

	result4, err4 := FetchInstructions(task4)

	// 验证结果
	assert.NoError(t, err4, "空任务名应该不会返回错误")
	assert.Equal(t, expectedResult4, result4, "空任务名应该返回空字符串")
}
