package biz

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
)

type CreateServiceModuleTool struct {
	commander Commander
}

func (s *CreateServiceModuleTool) Name() string {
	return "create_service_module"
}

func (s *CreateServiceModuleTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *CreateServiceModuleTool) Tags() fu.StringSet { return fu.StringSet{iipc.Service} }

func (s *CreateServiceModuleTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	pathVal, pathOK := param["path"].(string)
	if !pathOK || pathVal == "" {
		err = fmt.Errorf("missing or invalid path parameter")
		return
	}
	name, pathOK := param["name"].(string)
	if !pathOK || name == "" {
		err = fmt.Errorf("missing or invalid name parameter")
		return
	}
	fullPath := filepath.Join(workspaceDir, stage, pathVal)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}

	command := "magicCodeCli"
	var args = []string{
		"-actionType", "module",
		"-executeAction", "create",
		"-projectType", common.Service,
		"-moduleName", name,
	}
	runResult, runErr := s.commander.RunCommand(workspaceName, fullPath, command, true, args)
	if runErr != nil {
		err = errors.New(runErr.Error())
		return
	}
	result = runResult
	if result == "" {
		result = "创建Module成功"
	}
	return
}
