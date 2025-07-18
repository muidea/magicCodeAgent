package biz

import (
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCodeAgent/pkg/util"
)

type WorkspaceTool struct {
	commander    Commander
	workspacePtr *models.Workspace
	toolName     string
	toolPrompt   string
	toolCommand  string
}

func (s *WorkspaceTool) Name() string {
	return s.toolName
}

func (s *WorkspaceTool) Prompt() string {
	return s.toolPrompt
}

func (s *WorkspaceTool) Tags() fu.StringSet {
	return fu.StringSet{s.workspacePtr.Category}
}

func (s *WorkspaceTool) Execute(param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
	var result string
	defer func() {
		ret = &imbtc.ToolResult{
			Name: s.Name(),
			Result: models.MessageItem{
				Role: common.Tool,
				Content: []string{
					fmt.Sprintf("execute tool: %s", s.Name()),
				},
			},
		}
		if err != nil {
			ret.Result.Content = append(ret.Result.Content, err.Error())
		} else {
			ret.Result.Content = append(ret.Result.Content, result)
		}
	}()

	tmplParam := map[string]interface{}{}
	for k, v := range param {
		tmplParam[util.ToTitle(k)] = v
	}

	tplPtr, tplErr := template.New(s.Name()).Parse(s.toolCommand)
	if tplErr != nil {
		err = tplErr
		log.Errorf("parse tool command template failed:%s", tplErr.Error())
		return
	}

	var builder strings.Builder
	err = tplPtr.Execute(&builder, tmplParam)
	if err != nil {
		log.Errorf("execute tool command template failed:%s", err.Error())
		return
	}

	finalCommand := builder.String()
	runResult, runErr := s.commander.RunCommand(s.workspacePtr.GetName(), s.workspacePtr.GetWorkingPath(), finalCommand)
	if runErr != nil {
		log.Errorf("run tool:%s, command:%s failed: %s", s.Name(), finalCommand, runErr.Error())
		result = fmt.Sprintf("run %s failed, error:%s", s.Name(), runErr.Error())
		return
	}
	result = fmt.Sprintf("run %s finished.", s.Name())
	if runResult != "" {
		result = fmt.Sprintf("%s\n%s", result, runResult)
	}
	return
}

func extractCommandName(text string) string {
	re := regexp.MustCompile(`(?m)^###\s+(\S+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func NewWorkspaceTool(commander Commander, workspacePtr *models.Workspace, toolContent string) (ret imbtc.Tool, err error) {
	items := strings.Split(toolContent, "---")
	if len(items) != 2 {
		err = errors.New("invalid tool content format, split by --- should be 2")
		return
	}

	toolPrompt := items[0]
	toolName := extractCommandName(toolPrompt)
	if toolName == "" {
		err = errors.New("invalid tool content format, missing tool name")
		return
	}
	toolCommand := strings.TrimSpace(items[1])

	ret = &WorkspaceTool{
		commander:    commander,
		workspacePtr: workspacePtr,
		toolName:     toolName,
		toolPrompt:   toolPrompt,
		toolCommand:  toolCommand,
	}
	return
}
