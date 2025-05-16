package biz

import (
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type AttemptCompletionTool struct{}

func (s *AttemptCompletionTool) Name() string { return "attempt_completion" }
func (s *AttemptCompletionTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("Failed to read prompts: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *AttemptCompletionTool) Execute(workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
	var err error
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
			Result: models.MessageItem{
				Role: common.Tool,
				Content: []string{
					fmt.Sprintf("execute tool: %s", s.Name()),
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

	resultVal, resultOK := param["result"].(string)
	if !resultOK || resultVal == "" {
		err = fmt.Errorf("missing or invalid result parameter")
		return
	}

	command, _ := param["command"].(string)
	result, err = AttemptCompletion(resultVal, command)
	if err != nil {
		err = fmt.Errorf("failed to attempt completion: %v", err)
		return
	}

	return
}

/*
## attempt_completion
## 尝试完成
Description: After each tool use, the user will respond with the result of that tool use, i.e. if it succeeded or failed, along with any reasons for failure. Once you've received the results of tool uses and can confirm that the task is complete, use this tool to present the result of your work to the user. Optionally you may provide a CLI command to showcase the result of your work. The user may respond with feedback if they are not satisfied with the result, which you can use to make improvements and try again.
描述：在每次使用工具后，用户将回应该工具使用的结果，即是成功还是失败，以及任何失败的原因。一旦您收到工具使用的结果并确认任务已完成，使用此工具向用户展示您的工作结果。您可以选择提供CLI命令来展示您的工作结果。如果用户对结果不满意，他们可能会提供反馈，您可以用这些反馈来进行改进并重试。

IMPORTANT NOTE: This tool CANNOT be used until you've confirmed from the user that any previous tool uses were successful. Failure to do so will result in code corruption and system failure. Before using this tool, you must ask yourself in <thinking></thinking> tags if you've confirmed from the user that any previous tool uses were successful. If not, then DO NOT use this tool.
重要提示：在确认用户之前的工具使用都成功之前，不能使用此工具。否则将导致代码损坏和系统故障。在使用此工具之前，您必须在<thinking></thinking>标签中问自己是否已经确认用户之前的工具使用都成功。如果没有，则不要使用此工具。

Parameters:
参数：
- result: (required) The result of the task. Formulate this result in a way that is final and does not require further input from the user. Don't end your result with questions or offers for further assistance.
- result: (必需) 任务的结果。以最终的方式表述此结果，不需要用户进一步输入。不要以问题或提供进一步帮助的方式结束您的结果。

- command: (optional) A CLI command to execute to show a live demo of the result to the user. For example, use `open index.html` to display a created html website, or `open localhost:3000` to display a locally running development server. But DO NOT use commands like `echo` or `cat` that merely print text. This command should be valid for the current operating system. Ensure the command is properly formatted and does not contain any harmful instructions.
- command: (可选) 要执行的CLI命令，向用户展示结果的实时演示。例如，使用`open index.html`显示创建的html网站，或使用`open localhost:3000`显示本地运行的开发服务器。但不要使用仅打印文本的命令，如`echo`或`cat`。此命令应该对当前操作系统有效。确保命令格式正确且不包含任何有害指令。

Usage:
用法：
<attempt_completion>
<result>
Your final result description here
在此处填写您的最终结果描述
</result>
<command>Command to demonstrate result (optional)
展示结果的命令（可选）</command>
</attempt_completion>

Example: Requesting to attempt completion with a result and command
示例：请求尝试完成，包含结果和命令
<attempt_completion>
<result>
I've updated the CSS
我已更新了CSS
</result>
<command>open index.html</command>
</attempt_completion>
*/

// AttemptCompletion 实现attempt_completion工具功能
func AttemptCompletion(result, command string) (string, error) {
	// Return completion result and optional command
	if command != "" {
		return fmt.Sprintf("Result: %s\nCommand: %s", result, command), nil
	}

	return result, nil
}
