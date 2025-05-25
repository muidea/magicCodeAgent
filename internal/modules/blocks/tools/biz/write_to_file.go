package biz

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
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

func (s *WriteToFileTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *WriteToFileTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
			Name: s.Name(),
			Result: models.MessageItem{
				Role: common.Tool,
				Content: []string{
					fmt.Sprintf("execute tool: %s, param:%v", s.Name(), param),
				},
			},
		}
		if err != nil {
			ret.Result.Content = append(ret.Result.Content, err.Error())
		} else {
			ret.Result.Content = append(ret.Result.Content, result)
		}
	}()

	path, ok := param["path"].(string)
	if !ok || path == "" {
		err = fmt.Errorf("missing or invalid path parameter")
		return
	}

	content, ok := param["content"].(string)
	if !ok || content == "" {
		err = fmt.Errorf("missing or invalid content parameter")
		return
	}

	lineCount := 0
	if lineCountParam, ok := param["line_count"]; ok {
		switch v := lineCountParam.(type) {
		case float64:
			lineCount = int(v)
		case int:
			lineCount = v
		case string:
			_, err = fmt.Sscanf(v, "%d", &lineCount)
		default:
			err = fmt.Errorf("line_count parameter must be a number")
		}
		if lineCount < 1 {
			err = fmt.Errorf("line_count parameter must be greater than 0")
		}
	}
	if err != nil {
		return
	}
	fullPath := filepath.Join(workspaceDir, stage)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}
	err = WriteToFile(fullPath, path, content, lineCount)
	if err != nil {
		log.Errorf("WriteToFile failed, fullPath:%s, path:%s err: %v", fullPath, path, err)
		return
	}

	result = "write to file success!"
	return
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
