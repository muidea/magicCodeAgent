package biz

import (
	"errors"
	"fmt"
	"path/filepath"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type DebugWebProjectTool struct {
	commander Commander
}

func (s *DebugWebProjectTool) Name() string {
	return "debug_web_project"
}

func (s *DebugWebProjectTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *DebugWebProjectTool) Tags() []string { return []string{iipc.Web} }

func (s *DebugWebProjectTool) Execute(workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
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
	workspaceName, workspaceOK := param["name"]
	if !workspaceOK {
		err = fmt.Errorf("error: workspace name not found")
		return
	}
	var restart bool
	restartVal, restartOK := param["restart"]
	if restartOK {
		restart = restartVal.(bool)
	}

	var runStatusMsg, runErrorMsg string
	var runErr *cd.Error
	fullPath := filepath.Join(workspaceDir, stage)
	// yarn install
	_, runErrorMsg, runErr = s.commander.RunCommand(workspaceName.(string), fullPath, "yarn", restart, []string{"install"})
	if runErr != nil {
		err = errors.New(runErr.Error())
		return
	}
	if runErrorMsg != "" {
		result = runErrorMsg
		return
	}

	// yarn lint
	runStatusMsg, runErrorMsg, runErr = s.commander.RunCommand(workspaceName.(string), fullPath, "yarn", restart, []string{"lint"})
	if runErr != nil {
		err = errors.New(runErr.Error())
		return
	}
	if runErrorMsg != "" {
		result = runErrorMsg
		return
	}

	result = runStatusMsg

	// TODO 执行调试命令
	// yarn install 安装依赖
	// lint
	// test
	// build
	// 这其中任意一个命令执行出错，就返回失败

	return
}
