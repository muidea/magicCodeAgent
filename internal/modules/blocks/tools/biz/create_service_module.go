package biz

import (
	"fmt"
	"os"
	"path/filepath"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
)

type CreateServiceModuleTool struct{}

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

func (s *CreateServiceModuleTool) Execute(workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
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

	path, ok := param["path"].(string)
	if !ok || path == "" {
		err = fmt.Errorf("missing or invalid path parameter")
		return
	}
	name, ok := param["name"].(string)
	if !ok || name == "" {
		err = fmt.Errorf("missing or invalid name parameter")
		return
	}
	fullPath := filepath.Join(workspaceDir, stage)
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
		"-moduleName", name,
	}
	result, err = ExecuteCommand(fullPath, command, args, path)
	if err != nil {
		log.Errorf("failed to execute command: %v", err)
		return
	}

	return
}
