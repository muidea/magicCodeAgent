package biz

import (
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type AskFollowupQuestionTool struct{}

func (s *AskFollowupQuestionTool) Name() string { return "ask_followup_question" }

func (s *AskFollowupQuestionTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *AskFollowupQuestionTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *AskFollowupQuestionTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
	var err error
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
			Result: models.MessageItem{
				Role: common.Tool,
				Content: []string{
					fmt.Sprintf("execute tool: %s, param:%v", s.Name(), param),
				},
			},
		}
		if err != nil {
			ret.Result.Content = append(ret.Result.Content, err.Error())
		} else {
			ret.Final = true
			ret.Result.Content = append(ret.Result.Content, result)
		}
	}()

	question, ok := param["question"].(string)
	if !ok || question == "" {
		err = fmt.Errorf("missing or invalid question parameter")
		return
	}
	flowUp, ok := param["follow_up"].([]string)
	if !ok || len(flowUp) == 0 {
		err = fmt.Errorf("missing or invalid follow_up parameter")
		return
	}

	result, err = AskFollowupQuestion(question, flowUp)
	if err != nil {
		err = fmt.Errorf("failed to ask followup question: %v", err)
		return
	}

	return
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
func AskFollowupQuestion(question string, flowUp []string) (string, error) {
	// 构建格式化的回答
	result := fmt.Sprintf("Follow-up question: %s\nSuggested answers:\n", question)

	// 添加编号的建议答案
	for i, suggestion := range flowUp {
		result += fmt.Sprintf("%d. %s\n", i+1, suggestion)
	}

	return result, nil
}
