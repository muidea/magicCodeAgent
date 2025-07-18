package biz

import (
	"encoding/json"
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

type ListStageFilesTool struct {
	workspacePtr *models.Workspace
}

func (s *ListStageFilesTool) Name() string { return "list_stage_files" }
func (s *ListStageFilesTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadInnerToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *ListStageFilesTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *ListStageFilesTool) Execute(param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	parentPath := s.workspacePtr.WorkspaceDir
	fullPath := filepath.Join(parentPath, stageVal)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}
	recursive, _ := param["recursive"].(bool)
	listResult, listErr := ListFiles(fullPath, path, recursive)
	if listErr != nil {
		log.Errorf("ListFiles failed, fullPath:%s, path:%s, error:%v", fullPath, path, listErr)
		result = fmt.Sprintf("list stage files failed, path:%s, error:%s", path, listErr.Error())
		return
	}
	byteVal, byteErr := json.MarshalIndent(listResult, "", "  ")
	if byteErr != nil {
		log.Errorf("failed to marshal list result, error: %s", byteErr.Error())
		result = fmt.Sprintf("list stage files failed, path:%s, error:%s", path, byteErr.Error())
		return
	}

	result = fmt.Sprintf("list stage files successed, path:%s\n---\n%s", path, string(byteVal))
	return
}
