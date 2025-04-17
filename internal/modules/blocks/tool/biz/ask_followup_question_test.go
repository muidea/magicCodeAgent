package biz

import (
	"fmt"
	"strings"
	"testing"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
	"github.com/stretchr/testify/assert"
)

func TestAskFollowupQuestion(t *testing.T) {
	// 创建一个Tool实例用于测试
	eventHub := event.NewHub(100)
	backgroundRoutine := task.NewBackgroundRoutine(100)
	tool := New(eventHub, backgroundRoutine)

	// 测试用例1：正常情况 - 多个建议答案
	question := "What is the path to the frontend-config.json file?"
	followUp := []string{
		"./src/frontend-config.json",
		"./config/frontend-config.json",
		"./frontend-config.json",
	}

	result, err := tool.AskFollowupQuestion(question, followUp)

	// 验证结果
	assert.NoError(t, err, "应该没有错误")
	assert.Contains(t, result, question, "返回结果应该包含问题")
	for i, suggestion := range followUp {
		assert.Contains(t, result, suggestion, "返回结果应该包含建议答案")
		assert.Contains(t, result, fmt.Sprintf("%d.", i+1), "返回结果应该包含序号")
	}

	// 测试用例2：空问题
	emptyQuestion := ""
	followUp2 := []string{"Answer 1", "Answer 2"}

	result2, err2 := tool.AskFollowupQuestion(emptyQuestion, followUp2)

	// 验证结果
	assert.NoError(t, err2, "空问题应该不会导致错误")
	assert.Contains(t, result2, "Question: ", "返回结果应该包含问题标记")
	for _, suggestion := range followUp2 {
		assert.Contains(t, result2, suggestion, "返回结果应该包含建议答案")
	}

	// 测试用例3：空建议答案列表
	question3 := "What is your name?"
	emptyFollowUp := []string{}

	result3, err3 := tool.AskFollowupQuestion(question3, emptyFollowUp)

	// 验证结果
	assert.NoError(t, err3, "空建议答案列表应该不会导致错误")
	assert.Contains(t, result3, question3, "返回结果应该包含问题")
	assert.Contains(t, result3, "Suggested answers:", "返回结果应该包含建议答案标记")

	// 测试用例4：特殊字符
	question4 := "What is your favorite programming language? (C++, Java, Python)"
	followUp4 := []string{
		"C++ because of its performance",
		"Java because of its platform independence",
		"Python because of its simplicity and readability",
	}

	result4, err4 := tool.AskFollowupQuestion(question4, followUp4)

	// 验证结果
	assert.NoError(t, err4, "特殊字符应该不会导致错误")
	assert.Contains(t, result4, question4, "返回结果应该包含问题")
	for _, suggestion := range followUp4 {
		assert.Contains(t, result4, suggestion, "返回结果应该包含建议答案")
	}

	// 测试用例5：验证格式
	question5 := "Simple question?"
	followUp5 := []string{"Answer A", "Answer B"}

	result5, err5 := tool.AskFollowupQuestion(question5, followUp5)

	// 验证结果
	assert.NoError(t, err5, "应该没有错误")
	assert.True(t, strings.HasPrefix(result5, "Question: "), "返回结果应该以'Question: '开头")
	assert.Contains(t, result5, "Suggested answers:", "返回结果应该包含'Suggested answers:'")
	assert.Contains(t, result5, "1. Answer A", "返回结果应该包含编号的答案")
	assert.Contains(t, result5, "2. Answer B", "返回结果应该包含编号的答案")
}
