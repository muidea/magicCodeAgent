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

type LintWebProjectTool struct {
	commander Commander
}

func (s *LintWebProjectTool) Name() string {
	return "lint_web_project"
}

func (s *LintWebProjectTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *LintWebProjectTool) Tags() fu.StringSet { return fu.StringSet{iipc.Web} }

func (s *LintWebProjectTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	fullPath := filepath.Join(workspaceDir, stage)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}

	command := "yarn"
	var args = []string{
		"run", "lint",
	}
	runResult, runErr := s.commander.RunCommand(workspaceName, fullPath, command, true, args)
	if runErr != nil {
		err = errors.New(runErr.Error())
		return
	}
	result = runResult
	if result == "" {
		result = "代码检查结束"
	}
	return
}
