package biz

import (
	"encoding/json"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type UseMcpToolTool struct{}

func (s *UseMcpToolTool) Name() string { return "use_mcp_tool" }
func (s *UseMcpToolTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadInnerToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *UseMcpToolTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *UseMcpToolTool) Execute(workspaceName, workspaceDir, workspaceType, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	serverName, ok := param["server_name"].(string)
	if !ok || serverName == "" {
		err = fmt.Errorf("server_name is empty")
		return
	}

	toolName, ok := param["tool_name"].(string)
	if !ok || toolName == "" {
		err = fmt.Errorf("tool_name is empty")
		return
	}

	arguments, ok := param["arguments"].(string)
	if !ok || arguments == "" {
		err = fmt.Errorf("arguments is empty")
		return
	}

	resourceVal, resourceErr := UseMcpTool(serverName, toolName, json.RawMessage(arguments))
	if resourceErr != nil {
		log.Errorf("UseMcpTool failed, serverName:%s, toolName:%s, error: %v", serverName, toolName, err)
		result = fmt.Sprintf("use mcp tool failed, error: %s", resourceErr.Error())
		return
	}

	result = resourceVal
	return
}

/*
## use_mcp_tool
Description: Request to use a tool provided by a connected MCP server. Each MCP server can provide multiple tools with different capabilities. Tools have defined input schemas that specify required and optional parameters.
描述：请求使用连接的MCP服务器提供的工具。每个MCP服务器可以提供具有不同功能的多个工具。工具具有定义的输入模式，用于指定必需和可选参数。

Parameters 参数:
- server_name: (required) The name of the MCP server providing the tool
  服务器名称：（必需）提供工具的MCP服务器的名称
- tool_name: (required) The name of the tool to execute
  工具名称：（必需）要执行的工具的名称
- arguments: (required) A JSON object containing the tool's input parameters, following the tool's input schema
  参数：（必需）包含工具输入参数的JSON对象，遵循工具的输入模式

Usage 使用方法:
<use_mcp_tool>
<server_name>server name here</server_name>
<tool_name>tool name here</tool_name>
<arguments>
{
  "param1":"value1",
  "param2":"value2"
}
</arguments>
</use_mcp_tool>

Example 示例: Requesting to use an MCP tool 请求使用MCP工具

<use_mcp_tool>
<server_name>weather-server</server_name>
<tool_name>get_forecast</tool_name>
<arguments>
{
  "city":"San Francisco",
  "days":5
}
</arguments>
</use_mcp_tool>
*/

// UseMcpTool 实现use_mcp_tool工具功能
func UseMcpTool(serverName, toolName string, arguments json.RawMessage) (string, error) {
	// 这里应该实现与MCP服务器的通信逻辑
	// 由于这是一个模拟实现，我们只返回一个示例响应
	return fmt.Sprintf("Executed tool '%s' on server '%s' with arguments: %s", toolName, serverName, string(arguments)), nil
}
