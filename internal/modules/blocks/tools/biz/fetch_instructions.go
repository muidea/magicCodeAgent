package biz

import (
	"encoding/json"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type FetchInstructionsTool struct{}

func (s *FetchInstructionsTool) Name() string { return "fetch_instructions" }
func (s *FetchInstructionsTool) Prompts() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *FetchInstructionsTool) Execute(param imbtc.ToolParam) (imbtc.ToolResult, error) {
	// 参数校验
	task, ok := param["task"].(string)
	if !ok || task == "" {
		return nil, fmt.Errorf("缺少或无效的任务类型参数")
	}

	// 调用工具函数
	result, err := FetchInstructions(task)
	if err != nil {
		return nil, fmt.Errorf("获取指令失败: %v", err)
	}

	// 处理返回值
	ret, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("序列化结果失败: %v", err)
	}

	return ret, nil
}

/*
## fetch_instructions
Description: Request to fetch instructions to perform a task
描述：请求获取执行任务的指导说明
Parameters:
- task: (required) The task to get instructions for.  This can take the following values:
  任务类型：(必填) 需要获取指导说明的任务类型，可选值：
  create_mcp_server
  create_mode

Example: Requesting instructions to create an MCP Server
示例：请求创建MCP服务器的指导说明
<fetch_instructions>
<task>create_mcp_server</task>
</fetch_instructions>
*/

// FetchInstructions 实现fetch_instructions工具功能
// FetchInstructions function implements the fetch_instructions tool functionality
func FetchInstructions(task string) (string, error) {
	// 根据任务类型返回相应的指令
	switch task {
	case "create_mcp_server":
		return "Instructions for creating an MCP server...", nil
	case "create_mode":
		return "Instructions for creating a mode...", nil
	default:
		return "", nil
	}
}
