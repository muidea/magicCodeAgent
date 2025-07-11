package biz

import (
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

type CreateServiceInitiatorTool struct {
	commander Commander
}

func (s *CreateServiceInitiatorTool) Name() string {
	return "create_service_initiator"
}

func (s *CreateServiceInitiatorTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *CreateServiceInitiatorTool) Tags() fu.StringSet { return fu.StringSet{iipc.Service} }

func (s *CreateServiceInitiatorTool) Execute(workspaceName, workspaceDir, workspaceType, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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
	fullPath := workspaceDir
	if !fp.Exist(fullPath) {
		err = fmt.Errorf("workspace not exist: %s", fullPath)
		return
	}

	pathVal, pathOK := param["path"].(string)
	if !pathOK || pathVal == "" {
		err = fmt.Errorf("missing or invalid path parameter")
		return
	}
	name, ok := param["name"].(string)
	if !ok || name == "" {
		err = fmt.Errorf("missing or invalid name parameter")
		return
	}
	fullPath = filepath.Join(fullPath, pathVal)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}

	command := "magicCodeCli"
	var args = []string{
		"-actionType", "initiator",
		"-executeAction", "create",
		"-projectType", common.Service,
		"-moduleName", name,
	}

	runResult, runErr := s.commander.RunCommand(workspaceName, fullPath, command, true, 0, args)
	if runErr != nil {
		log.Errorf("run command failed, error:%s", runErr.Error())
		result = fmt.Sprintf("create service initiator failed, error:%s", runErr.Error())
		return
	}

	result = "create service initiator finished."
	if runResult != "" {
		result = fmt.Sprintf("%s\n%s", result, runResult)
	}
	return
}
