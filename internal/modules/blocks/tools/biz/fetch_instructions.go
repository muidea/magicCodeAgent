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

type FetchInstructionsTool struct{}

func (s *FetchInstructionsTool) Name() string { return "fetch_instructions" }
func (s *FetchInstructionsTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("Failed to read prompts: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *FetchInstructionsTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *FetchInstructionsTool) Execute(workspaceName, workspaceDir, workspaceType, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	task, ok := param["task"].(string)
	if !ok || task == "" {
		err = fmt.Errorf("task is empty")
		return
	}

	result, err = FetchInstructions(task)
	if err != nil {
		log.Errorf("Failed to fetch instructions: %v", err)
		return
	}

	return
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
	// Return corresponding instructions based on task type
	switch task {
	case "create_mcp_server":
		return "Instructions for creating an MCP server...", nil
	case "create_mode":
		return "Instructions for creating a mode...", nil
	default:
		return "", nil
	}
}
