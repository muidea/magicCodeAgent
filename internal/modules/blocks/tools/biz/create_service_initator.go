package biz

import (
	"fmt"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCommon/foundation/log"
)

type CreateServiceInitatorTool struct{}

func (s *CreateServiceInitatorTool) Name() string {
	return "create_service_initator"
}

func (s *CreateServiceInitatorTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *CreateServiceInitatorTool) Execute(workspaceDir string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
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

	command := "magicCodeCli"
	var args = []string{
		"-actionType", "initator",
		"-executeAction", "create",
		"-moduleName", name,
	}
	result, err = ExecuteCommand(workspaceDir, command, args, path)
	if err != nil {
		log.Errorf("failed to execute command: %v", err)
		return
	}

	return
}
