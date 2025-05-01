package biz

import (
	"fmt"
)

type AskFollowupQuestion struct{}

func (s *AskFollowupQuestion) Name() string { return "ask_followup_question" }
func (s *AskFollowupQuestion) Prompts() string {
	return ``
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
func (s *Tools) AskFollowupQuestion(question string, followUp []string) (string, error) {
	// 返回问题和建议的回答
	result := fmt.Sprintf("Question: %s\nSuggested answers:\n", question)
	for i, suggestion := range followUp {
		result += fmt.Sprintf("%d. %s\n", i+1, suggestion)
	}
	return result, nil
}
