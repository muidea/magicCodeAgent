package biz

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
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

func (s *InitializationWebProjectTool) Tags() fu.StringSet { return fu.StringSet{iipc.Web} }

func (s *InitializationWebProjectTool) Execute(workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
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
	fullPath := filepath.Join(workspaceDir, stage)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}

	command := "magicCodeCli"
	var args = []string{
		"-actionType", "project",
		"-executeAction", "create",
		"-projectName", path.Join("web", name),
		"-projectDescription", description,
	}
	result, err = ExecuteCommand(fullPath, command, args, "")
	if err != nil {
		log.Errorf("failed to execute command: %v", err)
		return
	}

	return
}
