package biz

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type ReadStageFileTool struct{}

func (s *ReadStageFileTool) Name() string {
	return "read_stage_file"
}

func (s *ReadStageFileTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *ReadStageFileTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *ReadStageFileTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
	var err error
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
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

	stageVal, ok := param["stage"].(string)
	if !ok || stageVal == "" {
		err = fmt.Errorf("missing or invalid stageVal parameter")
		return
	}

	path, ok := param["path"].(string)
	if !ok || path == "" {
		err = fmt.Errorf("missing or invalid path parameter")
		return
	}

	if strings.Contains(path, "..") || strings.HasPrefix(path, "/") {
		err = fmt.Errorf("invalid path parameter")
		return
	}

	startLine := 1
	if startLineParam, exists := param["start_line"]; exists {
		switch v := startLineParam.(type) {
		case float64:
			startLine = int(v)
		case int:
			startLine = v
		default:
			err = fmt.Errorf("start_line parameter must be a number")
			return
		}
		if startLine < 1 {
			err = fmt.Errorf("start_line parameter must be greater than 0")
			return
		}
	}

	endLine := 0 // 0表示读取到文件末尾
	if endLineParam, exists := param["end_line"]; exists {
		switch v := endLineParam.(type) {
		case float64:
			endLine = int(v)
		case int:
			endLine = v
		default:
			err = fmt.Errorf("end_line parameter must be a number")
			return
		}
		if endLine > 0 && endLine < startLine {
			err = fmt.Errorf("end_line parameter must be greater than start_line")
			return
		}
	}
	fullPath := filepath.Join(workspaceDir, stageVal)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}
	result, err = ReadFile(fullPath, path, startLine, endLine)
	if err != nil {
		log.Errorf("ReadFile failed, fullPath:%s, path:%s err: %v", fullPath, path, err)
	}
	return
}
