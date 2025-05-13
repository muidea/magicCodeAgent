package biz

import (
	"fmt"
	"path"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCommon/foundation/log"
)

type InitializationServiceProjectTool struct{}

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

func (s *InitializationServiceProjectTool) Execute(workspaceDir string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
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
		"-projectName", path.Join("service", name),
		"-projectDescription", description,
	}
	result, err = ExecuteCommand(workspaceDir, command, args, "")
	if err != nil {
		log.Errorf("failed to execute command: %v", err)
		return
	}

	return
}
