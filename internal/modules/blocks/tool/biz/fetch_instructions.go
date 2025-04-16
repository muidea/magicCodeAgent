package biz

/*
## fetch_instructions
Description: Request to fetch instructions to perform a task
Parameters:
- task: (required) The task to get instructions for.  This can take the following values:
  create_mcp_server
  create_mode

Example: Requesting instructions to create an MCP Server

<fetch_instructions>
<task>create_mcp_server</task>
</fetch_instructions>
*/

// FetchInstructions 实现fetch_instructions工具功能
func (s *Tool) FetchInstructions(task string) (string, error) {
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
