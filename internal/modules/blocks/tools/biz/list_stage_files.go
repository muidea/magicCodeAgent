package biz

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type ListStageFilesTool struct{}

func (s *ListStageFilesTool) Name() string { return "list_stage_files" }
func (s *ListStageFilesTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *ListStageFilesTool) Execute(workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
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
	fullPath := filepath.Join(workspaceDir, stageVal)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}
	recursive, _ := param["recursive"].(bool)
	listResult, listErr := ListFiles(fullPath, path, recursive)
	if listErr != nil {
		err = listErr
		log.Errorf("ListFiles failed, fullPath:%s, path:%s err: %v", fullPath, path, err)
		log.Errorf("failed to list files: %v", err)
		return
	}
	byteVal, byteErr := json.MarshalIndent(listResult, "", "  ")
	if byteErr != nil {
		err = byteErr
		return
	}

	result = string(byteVal)
	return
}
