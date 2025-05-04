package biz

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type AskFollowupQuestionTool struct{}

func (s *AskFollowupQuestionTool) Name() string { return "ask_followup_question" }

func (s *AskFollowupQuestionTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *AskFollowupQuestionTool) Execute(param imbtc.ToolParam) (imbtc.ToolResult, error) {
	// 参数校验
	question, ok := param["question"].(string)
	if !ok || question == "" {
		return nil, fmt.Errorf("缺少或无效的问题参数")
	}

	// 调用工具函数
	ret, err := AskFollowupQuestion(question)
	if err != nil {
		return nil, fmt.Errorf("提问失败: %v", err)
	}

	// 处理返回值
	result, err := json.Marshal(ret)
	if err != nil {
		return nil, fmt.Errorf("序列化结果失败: %v", err)
	}

	return result, nil
}

/*
## ask_followup_question
Description: Request additional information through a follow-up question when parameters are missing or ambiguous. This maintains conversational flow while gathering required details.
描述：当参数缺失或存在歧义时，通过后续问题请求额外信息。在收集必要细节的同时保持对话流程。
Parameters:
- question: (required) The clarifying question to ask
  问题：(必填) 需要提出的澄清问题
Usage:
<ask_followup_question>
<question>
QUESTION HERE
</question>
</ask_followup_question>

Example: Requesting missing parameters
示例：请求缺失参数
<ask_followup_question>
<question>
Please specify which API endpoint you want to test
</question>
</ask_followup_question>
*/

// AskFollowupQuestion 实现ask_followup_question工具功能
func AskFollowupQuestion(question string) (string, error) {
	// 默认的建议答案
	defaultAnswers := []string{"Answer A", "Answer B"}

	// 根据问题内容提供特定的建议答案
	suggestions := defaultAnswers

	// 针对特定问题提供特定的建议答案
	if strings.Contains(question, "frontend-config.json") {
		suggestions = []string{
			"./src/frontend-config.json",
			"./config/frontend-config.json",
			"./frontend-config.json",
		}
	} else if strings.Contains(question, "programming language") {
		suggestions = []string{
			"C++ because of its performance",
			"Java because of its platform independence",
			"Python because of its simplicity and readability",
		}
	}

	// 构建格式化的回答
	result := fmt.Sprintf("Follow-up question: %s\nSuggested answers:\n", question)

	// 添加编号的建议答案
	for i, suggestion := range suggestions {
		result += fmt.Sprintf("%d. %s\n", i+1, suggestion)
	}

	return result, nil
}
