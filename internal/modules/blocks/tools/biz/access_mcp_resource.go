package biz

import (
	"fmt"

	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type AccessMcpResourceTool struct{}

func (s *AccessMcpResourceTool) Name() string { return "access_mcp_resource" }
func (s *AccessMcpResourceTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *AccessMcpResourceTool) Execute(workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
	var err error
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
			Result: models.MessageItem{
				Content: []string{
					fmt.Sprintf("execute tool: %s, param:%v", s.Name(), param),
				},
			},
		}
		if err != nil {
			ret.Final = true
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

	uri, ok := param["uri"].(string)
	if !ok || uri == "" {
		err = fmt.Errorf("uri is empty")
		return
	}

	resourceVal, resourceErr := AccessMcpResource(serverName, uri)
	if resourceErr != nil {
		err = resourceErr
		return
	}

	result = resourceVal
	return
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

// AccessMcpResource implements the access_mcp_resource tool functionality
// AccessMcpResource function implements the access_mcp_resource tool functionality
func AccessMcpResource(serverName, uri string) (string, error) {
	// Should implement the resource access logic with MCP server
	// As this is a mock implementation, we just return a sample response
	return fmt.Sprintf("Accessed resource '%s' on server '%s'", uri, serverName), nil
}
