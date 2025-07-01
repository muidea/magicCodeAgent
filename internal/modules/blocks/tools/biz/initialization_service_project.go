package biz

import (
	"errors"
	"fmt"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
)

type InitializationServiceProjectTool struct {
	commander Commander
}

func (s *InitializationServiceProjectTool) Name() string {
	return "initialization_service_project"
}

func (s *InitializationServiceProjectTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *InitializationServiceProjectTool) Tags() fu.StringSet { return fu.StringSet{iipc.Service} }

func (s *InitializationServiceProjectTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	name, ok := param["name"].(string)
	if !ok || name == "" {
		err = fmt.Errorf("missing or invalid name parameter")
		return
	}

	description := ""
	descriptionVal, descriptionOK := param["description"]
	if descriptionOK {
		description = descriptionVal.(string)
	}
	command := "magicCodeCli"
	var args = []string{
		"-actionType", "project",
		"-executeAction", "create",
		"-projectName", name,
		"-projectType", common.Service,
		"-projectDescription", description,
	}
	runResult, runErr := s.commander.RunCommand(workspaceName, fullPath, command, true, 0, args)
	if runErr != nil {
		err = errors.New(runErr.Error())
		return
	}
	result = runResult
	if result == "" {
		result = "初始化项目成功"
	}
	return
}
