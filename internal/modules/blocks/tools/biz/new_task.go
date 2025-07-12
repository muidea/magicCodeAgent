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

type NewTaskTool struct{}

func (s *NewTaskTool) Name() string { return "new_task" }
func (s *NewTaskTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *NewTaskTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *NewTaskTool) Execute(workspaceName, workspaceDir, workspaceType, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
			Name: s.Name(),
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
			ret.Result.Content = append(ret.Result.Content, result)
		}
	}()

	mode, ok := param["mode"].(string)
	if !ok || mode == "" {
		err = fmt.Errorf("missing or invalid mode parameter")
		return
	}

	message, ok := param["message"].(string)
	if !ok || message == "" {
		err = fmt.Errorf("missing or invalid message parameter")
		return
	}

	result, err = NewTask(mode, message)
	if err != nil {
		return
	}

	return
}

/*
/*
## new_task
Description: Create a new task with a specified starting mode and initial message. This tool instructs the system to create a new Cline instance in the given mode with the provided message.
描述：创建一个具有指定启动模式和初始消息的新任务。该工具指示系统在给定模式下使用提供的消息创建新的Cline实例。

Parameters 参数:
- mode: (required) The slug of the mode to start the new task in (e.g., "code", "ask", "architect").
- mode: (必填) 启动新任务的模式标识符（例如："code"、"ask"、"architect"）。
- message: (required) The initial user message or instructions for this new task.
- message: (必填) 此新任务的初始用户消息或指令。

Usage 用法:
<new_task>
<mode>your-mode-slug-here</mode>
<message>Your initial instructions here</message>
</new_task>

Example 示例:
<new_task>
<mode>code</mode>
<message>Implement a new feature for the application.</message>
<message>实现应用程序的新功能。</message>
</new_task>
*/

// NewTask 实现new_task工具功能
// NewTask function implements the new_task tool functionality
func NewTask(mode, message string) (string, error) {
	// 返回新任务创建信息
	return fmt.Sprintf("Created new task in mode '%s' with message: %s", mode, message), nil
}
