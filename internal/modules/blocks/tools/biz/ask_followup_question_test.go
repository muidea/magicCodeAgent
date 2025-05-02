package biz

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAskFollowupQuestion(t *testing.T) {
	// 测试用例1：正常情况 - 多个建议答案
	question := "What is the path to the frontend-config.json file?"
	followUp := []string{
		"./src/frontend-config.json",
		"./config/frontend-config.json",
		"./frontend-config.json",
	}

	result, err := AskFollowupQuestion(question)
	assert.NoError(t, err, "应该没有错误")
	// 验证结果
	assert.True(t, strings.HasPrefix(result, "Follow-up question: "), "返回结果应该以'Follow-up question: '开头")
	assert.Contains(t, result, "Suggested answers:", "返回结果应该包含'Suggested answers:'")
	for i, suggestion := range followUp {
		assert.Contains(t, result, fmt.Sprintf("%d. %s", i+1, suggestion), "返回结果应该包含编号的答案")
	}

	// 测试用例2：空问题
	emptyQuestion := ""

	result2, err2 := AskFollowupQuestion(emptyQuestion)
	assert.NoError(t, err2, "空问题应该不会导致错误")
	// 验证结果
	assert.True(t, strings.HasPrefix(result2, "Follow-up question: "), "返回结果应该以'Follow-up question: '开头")
	assert.Contains(t, result2, "Suggested answers:", "返回结果应该包含'Suggested answers:'")
	assert.Contains(t, result2, "1. Answer A", "返回结果应该包含编号的答案")
	assert.Contains(t, result2, "2. Answer B", "返回结果应该包含编号的答案")

	// 测试用例3：空建议答案列表
	question3 := "What is your name?"
	result3, err3 := AskFollowupQuestion(question3)
	assert.NoError(t, err3, "空建议答案列表应该不会导致错误")
	// 验证结果
	assert.True(t, strings.HasPrefix(result3, "Follow-up question: "), "返回结果应该以'Follow-up question: '开头")
	assert.Contains(t, result3, "Suggested answers:", "返回结果应该包含'Suggested answers:'")

	// 测试用例4：特殊字符
	question4 := "What is your favorite programming language? (C++, Java, Python)"

	result4, err4 := AskFollowupQuestion(question4)
	assert.NoError(t, err4, "特殊字符应该不会导致错误")
	// 验证结果
	assert.True(t, strings.HasPrefix(result4, "Follow-up question: "), "返回结果应该以'Follow-up question: '开头")
	assert.Contains(t, result4, "Suggested answers:", "返回结果应该包含'Suggested answers:'")
	assert.Contains(t, result4, "1. C++ because of its performance", "返回结果应该包含编号的答案")
	assert.Contains(t, result4, "2. Java because of its platform independence", "返回结果应该包含编号的答案")
	assert.Contains(t, result4, "3. Python because of its simplicity and readability", "返回结果应该包含编号的答案")

	// 测试用例5：验证格式
	question5 := "Simple question?"
	result5, err5 := AskFollowupQuestion(question5)

	// 验证结果
	assert.NoError(t, err5, "应该没有错误")
	assert.True(t, strings.HasPrefix(result5, "Follow-up question: "), "返回结果应该以'Follow-up question: '开头")
	assert.Contains(t, result5, "Suggested answers:", "返回结果应该包含'Suggested answers:'")
	assert.Contains(t, result5, "1. Answer A", "返回结果应该包含编号的答案")
	assert.Contains(t, result5, "2. Answer B", "返回结果应该包含编号的答案")

}
