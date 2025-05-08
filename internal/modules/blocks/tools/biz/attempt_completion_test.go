package biz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttemptCompletion(t *testing.T) {
	// 测试用例1：只有结果，没有命令
	result1 := "Task completed successfully"
	output1, err := AttemptCompletion(result1, "")
	assert.NoError(t, err, "尝试完成任务应该成功")
	assert.Equal(t, "Task completed successfully", output1, "输出应该只包含结果信息")

	// 测试用例2：有结果和命令
	result2 := "Created a new web application"
	command2 := "open index.html"
	output2, err := AttemptCompletion(result2, command2)
	assert.NoError(t, err, "尝试完成任务应该成功")
	expected2 := "Result: Created a new web application\nCommand: open index.html"
	assert.Equal(t, expected2, output2, "输出应该包含结果和命令信息")

	// 测试用例3：空结果
	result3 := ""
	output3, err := AttemptCompletion(result3, "")
	assert.NoError(t, err, "尝试完成任务应该成功，即使结果为空")
	assert.Equal(t, "", output3, "输出应该包含空结果信息")

	// 测试用例4：空结果但有命令
	result4 := ""
	command4 := "npm start"
	output4, err := AttemptCompletion(result4, command4)
	assert.NoError(t, err, "尝试完成任务应该成功，即使结果为空")
	expected4 := "Result: \nCommand: npm start"
	assert.Equal(t, expected4, output4, "输出应该包含空结果和命令信息")
}
