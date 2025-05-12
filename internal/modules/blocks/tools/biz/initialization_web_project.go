package biz

import (
	"fmt"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCommon/foundation/log"
)

type InitializationWebProjectTool struct{}

func (s *InitializationWebProjectTool) Name() string {
	return "initialization_web_project"
}

func (s *InitializationWebProjectTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *InitializationWebProjectTool) Execute(workspaceDir string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
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

	command := "magicCodeCli"
	var args = []string{
		"-actionType", "project/web",
		"-executeAction", "create",
		"-projectName", name,
	}
	result, err = ExecuteCommand(workspaceDir, command, args, "")
	if err != nil {
		log.Errorf("failed to execute command: %v", err)
		return
	}

	return
}
