package biz

import (
	"encoding/json"
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type AccessMcpResourceTool struct{}

func (s *AccessMcpResourceTool) Name() string { return "access_mcp_resource" }
func (s *AccessMcpResourceTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *AccessMcpResourceTool) Execute(workspaceDir string, param imbtc.ToolParam) (imbtc.ToolResult, error) {
	// 参数校验
	serverName, ok := param["server_name"].(string)
	if !ok || serverName == "" {
		return nil, fmt.Errorf("缺少或无效的服务器名称参数")
	}

	uri, ok := param["uri"].(string)
	if !ok || uri == "" {
		return nil, fmt.Errorf("缺少或无效的URI参数")
	}

	// 调用工具函数
	ret, err := AccessMcpResource(serverName, uri)
	if err != nil {
		return nil, fmt.Errorf("访问MCP资源失败: %v", err)
	}

	// 处理返回值
	result, err := json.Marshal(ret)
	if err != nil {
		return nil, fmt.Errorf("序列化结果失败: %v", err)
	}

	return result, nil
}

/*
## access_mcp_resource
Description: Request to access a resource provided by a connected MCP server. Resources represent data sources that can be used as context, such as files, API responses, or system information.
描述：请求访问由连接的MCP服务器提供的资源。资源代表可用作上下文的数据源，例如文件、API响应或系统信息。
Parameters:
- server_name: (required) The name of the MCP server providing the resource
  服务器名称：(必填) 提供资源的MCP服务器名称
- uri: (required) The URI identifying the specific resource to access
  资源URI：(必填) 标识要访问的特定资源的URI
Usage:
<access_mcp_resource>
<server_name>server name here</server_name>
<uri>resource URI here</uri>
</access_mcp_resource>

示例：请求访问MCP资源

<access_mcp_resource>
<server_name>weather-server</server_name>
<uri>weather://san-francisco/current</uri>
</access_mcp_resource>
*/

// AccessMcpResource 实现access_mcp_resource工具功能
// AccessMcpResource function implements the access_mcp_resource tool functionality
func AccessMcpResource(serverName, uri string) (string, error) {
	// 这里应该实现与MCP服务器的资源访问逻辑
	// 由于这是一个模拟实现，我们只返回一个示例响应
	return fmt.Sprintf("Accessed resource '%s' on server '%s'", uri, serverName), nil
}
