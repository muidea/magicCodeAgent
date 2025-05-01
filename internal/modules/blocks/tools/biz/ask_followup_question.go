package biz

import (
	"encoding/json"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type AskFollowupQuestionTool struct{}

func (s *AskFollowupQuestionTool) Name() string { return "ask_followup_question" }

func (s *AskFollowupQuestionTool) Prompts() string {
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
	// 返回问题
	return fmt.Sprintf("Follow-up question: %s", question), nil
}
