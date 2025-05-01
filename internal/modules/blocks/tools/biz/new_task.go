package biz

import (
	"encoding/json"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type NewTaskTool struct{}

func (s *NewTaskTool) Name() string { return "new_task" }
func (s *NewTaskTool) Prompts() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *NewTaskTool) Execute(param imbtc.ToolParam) (imbtc.ToolResult, error) {
	mode, ok := param["mode"].(string)
	if !ok || mode == "" {
		return nil, fmt.Errorf("缺少或无效的模式参数")
	}

	message, ok := param["message"].(string)
	if !ok || message == "" {
		return nil, fmt.Errorf("缺少或无效的消息参数")
	}

	ret, err := NewTask(mode, message)
	if err != nil {
		return nil, fmt.Errorf("创建新任务失败: %v", err)
	}

	result, err := json.Marshal(ret)
	if err != nil {
		return nil, fmt.Errorf("序列化结果失败: %v", err)
	}

	return result, nil
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
