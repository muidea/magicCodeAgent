package biz

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type WriteToFileTool struct{}

func (s *WriteToFileTool) Name() string { return "write_to_file" }
func (s *WriteToFileTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *WriteToFileTool) Execute(workspaceDir string, param imbtc.ToolParam) *imbtc.ToolResult {
	path, ok := param["path"].(string)
	if !ok || path == "" {
		return &imbtc.ToolResult{
			Final: true,
			Result: models.MessageItem{
				Content: []string{
					fmt.Sprintf("execute tool: %s", s.Name()),
					"missing or invalid path parameter",
				},
			},
		}
	}

	content, ok := param["content"].(string)
	if !ok || content == "" {
		return &imbtc.ToolResult{
			Final: true,
			Result: models.MessageItem{
				Content: []string{
					fmt.Sprintf("execute tool: %s", s.Name()),
					"missing or invalid content parameter",
				},
			},
		}
	}

	// 获取行数参数
	lineCount := 0
	if lineCountParam, ok := param["line_count"].(float64); ok {
		lineCount = int(lineCountParam)
	}

	err := WriteToFile(workspaceDir, path, content, lineCount)
	if err != nil {
		return &imbtc.ToolResult{
			Final: true,
			Result: models.MessageItem{
				Content: []string{
					fmt.Sprintf("execute tool: %s", s.Name()),
					fmt.Sprintf("failed to write to file: %v", err),
				},
			},
		}
	}

	return &imbtc.ToolResult{
		Result: models.MessageItem{
			Content: []string{
				fmt.Sprintf("execute tool: %s", s.Name()),
				"write to file success!",
			},
		},
	}
}

/*
## write_to_file
Description: Request to write full content to a file at the specified path. If the file exists, it will be overwritten with the provided content. If the file doesn't exist, it will be created. This tool will automatically create any directories needed to write the file.
Parameters:
- path: (required) The path of the file to write to (relative to the current workspace directory /home/rangh/webspace)
- content: (required) The content to write to the file. ALWAYS provide the COMPLETE intended content of the file, without any truncation or omissions. You MUST include ALL parts of the file, even if they haven't been modified. Do NOT include the line numbers in the content though, just the actual content of the file.
- line_count: (required) The number of lines in the file. Make sure to compute this based on the actual content of the file, not the number of lines in the content you're providing.
Usage:
<write_to_file>
<path>File path here</path>
<content>
Your file content here
</content>
<line_count>total number of lines in the file, including empty lines</line_count>
</write_to_file>

Example: Requesting to write to frontend-config.json
<write_to_file>
<path>frontend-config.json</path>
<content>
{
  \"apiEndpoint\":\"https://api.example.com\",
  \"theme\":{
    \"primaryColor\":\"#007bff\",
    \"secondaryColor\":\"#6c757d\",
    \"fontFamily\":\"Arial, sans-serif\"
  },
  \"features\":{
    \"darkMode\":true,
    \"notifications\":true,
    \"analytics\":false
  },
  \"version\":\"1.0.0\"
}
</content>
<line_count>14</line_count>
</write_to_file>

*/

// WriteToFile 实现write_to_file工具功能
func WriteToFile(workspaceDir, path, content string, lineCount int) error {
	fullPath := filepath.Join(workspaceDir, path)

	// 创建必要的目录
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(fullPath, []byte(content), 0644)
}
