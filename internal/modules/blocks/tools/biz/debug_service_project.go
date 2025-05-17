package biz

import (
	"fmt"
	"path/filepath"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCommon/foundation/log"
)

type DebugServiceProjectTool struct {
	commander Commander
}

func (s *DebugServiceProjectTool) Name() string {
	return "debug_service_project"
}

func (s *DebugServiceProjectTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *DebugServiceProjectTool) Execute(workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
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

	fullPath := filepath.Join(workspaceDir, stage)
	s.commander.RunCommand(workspaceName.(string), fullPath, "go", restart, []string{"lint", "tidy"})

	//
	// TODO 执行调试命令
	// lint
	// test
	// build
	// 这其中任意一个命令执行出错，就返回失败
	return
}
