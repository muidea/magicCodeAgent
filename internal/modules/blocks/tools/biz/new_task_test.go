package biz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	// 测试用例1：正常情况
	mode1 := "code"
	message1 := "Implement a new feature"
	expectedResult1 := "Created new task in mode 'code' with message: Implement a new feature"

	result1, err := NewTask(mode1, message1)

	// 验证结果
	assert.NoError(t, err, "创建新任务应该成功")
	assert.Equal(t, expectedResult1, result1, "返回结果应该符合预期")

	// 测试用例2：不同的模式
	mode2 := "ask"
	message2 := "How does this work?"
	expectedResult2 := "Created new task in mode 'ask' with message: How does this work?"

	result2, err := NewTask(mode2, message2)

	// 验证结果
	assert.NoError(t, err, "创建新任务应该成功")
	assert.Equal(t, expectedResult2, result2, "返回结果应该符合预期")

	// 测试用例3：空模式
	mode3 := ""
	message3 := "Test message"
	expectedResult3 := "Created new task in mode '' with message: Test message"

	result3, err := NewTask(mode3, message3)

	// 验证结果
	assert.NoError(t, err, "空模式应该不会导致错误")
	assert.Equal(t, expectedResult3, result3, "返回结果应该符合预期")

	// 测试用例4：空消息
	mode4 := "architect"
	message4 := ""
	expectedResult4 := "Created new task in mode 'architect' with message: "

	result4, err := NewTask(mode4, message4)

	// 验证结果
	assert.NoError(t, err, "空消息应该不会导致错误")
	assert.Equal(t, expectedResult4, result4, "返回结果应该符合预期")

	// 测试用例5：特殊字符
	mode5 := "code!@#"
	message5 := "Special chars: !@#$%^&*()"
	expectedResult5 := "Created new task in mode 'code!@#' with message: Special chars: !@#$%^&*()"

	result5, err := NewTask(mode5, message5)

	// 验证结果
	assert.NoError(t, err, "特殊字符应该不会导致错误")
	assert.Equal(t, expectedResult5, result5, "返回结果应该符合预期")
}
