package biz

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/google/shlex"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/task"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbcc "github.com/muidea/magicCodeAgent/internal/modules/blocks/commander/pkg/common"
	imbtmc "github.com/muidea/magicCodeAgent/internal/modules/blocks/templatemanager/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Commander interface {
	RunCommand(workspaceName, workspaceDir, command string) (result string, err *cd.Error)
}

type TemplateQuerier interface {
	QueryTemplate(name string) (result *models.Template, err *cd.Error)
}

type Tools struct {
	biz.Base

	toolRegistry *ToolRegistry
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
	promptsRegistry iipc.PromptsRegistry,
) *Tools {
	ptr := &Tools{
		Base: biz.New(imbtc.ToolsModule, eventHub, backgroundRoutine),
	}
	ptr.toolRegistry = NewToolRegistry(ptr, ptr, promptsRegistry)

	ptr.SubscribeFunc(imbtc.ReloadTools, ptr.Reload)
	ptr.SubscribeFunc(imbtc.ExecuteTool, ptr.Execute)
	return ptr
}

func (s *Tools) Run() (err *cd.Error) {
	return
}

func (s *Tools) execDataContent(dataContent string) (imbtc.ToolParam, *imbtc.ToolResult, error) {
	toolName, toolParam, toolErr := parseParamContent(dataContent)
	if toolErr != nil {
		log.Errorf("failed to parse param content: %s", toolErr.Error())
		return nil, &imbtc.ToolResult{
			Name: toolName,
			Result: models.MessageItem{
				Role:    common.Tool,
				Content: []string{fmt.Sprintf("parse tool param failed, param:\n%s\nerror:%s", dataContent, toolErr.Error())},
			},
		}, nil
	}
	if toolName == "" {
		log.Errorf("tool name cannot be empty")
		return nil, &imbtc.ToolResult{
			Name: toolName,
			Result: models.MessageItem{
				Role:    common.Tool,
				Content: []string{"illegal tool param, tool name cannot be empty. param:\n%s", dataContent},
			},
		}, nil
	}

	if toolName == "thinking" {
		return nil, &imbtc.ToolResult{
			Name: toolName,
			Result: models.MessageItem{
				Role:    common.Tool,
				Content: []string{"thinking ok!"},
			},
		}, nil
	}

	executeResult, executeErr := s.toolRegistry.Execute(toolName, toolParam)
	return toolParam, executeResult, executeErr
}

// 解析执行参数
// content里解析出执行函数名和参数
// content的格式是xml
// 执行函数是xml的root节点名
// 执行函数的参数是root的子节点及子节点的值
// 参照示例参见各个工具函数对应的注释信息
func parseParamContent(content string) (string, imbtc.ToolParam, error) {
	// 检查内容是否为空
	if strings.TrimSpace(content) == "" {
		return "", nil, fmt.Errorf("empty content")
	}

	// 提取根元素名称
	start := strings.Index(content, "<")
	if start == -1 {
		return "", nil, fmt.Errorf("invalid XML: missing opening tag")
	}

	// 检查是否是XML声明而不是根元素
	if strings.HasPrefix(strings.TrimSpace(content[start:]), "<?xml") {
		return "", nil, fmt.Errorf("invalid XML: missing root element")
	}

	end := strings.Index(content[start+1:], ">")
	if end == -1 {
		return "", nil, fmt.Errorf("invalid XML: missing closing tag")
	}

	// 获取根元素名称
	tagName := strings.TrimSpace(content[start+1 : start+1+end])
	if strings.Contains(tagName, " ") {
		tagName = strings.Split(tagName, " ")[0]
	}

	// 检查是否有匹配的结束标签
	closingTag := "</" + tagName + ">"
	if !strings.Contains(content, closingTag) {
		return "", nil, fmt.Errorf("invalid XML: missing closing tag for root element '%s'", tagName)
	}

	// 创建参数映射
	param := imbtc.ToolParam{}

	// 特殊处理包含特殊XML字符的工具
	if tagName == "apply_diff" || tagName == "write_to_file" {
		// 提取path参数
		pathStart := strings.Index(content, "<path>")
		pathEnd := strings.Index(content, "</path>")
		if pathStart != -1 && pathEnd != -1 {
			pathValue := content[pathStart+6 : pathEnd]
			param["path"] = strings.TrimSpace(pathValue)
		}

		switch tagName {
		case "apply_diff":
			diffStart := strings.Index(content, "<diff>")
			diffEnd := strings.Index(content, "</diff>")
			if diffStart != -1 && diffEnd != -1 {
				diffValue := content[diffStart+6 : diffEnd]
				param["diff"] = diffValue
			}
		case "write_to_file":
			// 提取content参数
			contentStart := strings.Index(content, "<content>")
			contentEnd := strings.Index(content, "</content>")
			if contentStart != -1 && contentEnd != -1 {
				contentValue := content[contentStart+9 : contentEnd]
				// 移除内容开头的换行符
				if len(contentValue) > 0 && contentValue[0] == '\n' {
					contentValue = contentValue[1:]
				}
				param["content"] = contentValue
			}

			// 提取line_count参数
			lineCountStart := strings.Index(content, "<line_count>")
			lineCountEnd := strings.Index(content, "</line_count>")
			if lineCountStart != -1 && lineCountEnd != -1 {
				lineCountValue := content[lineCountStart+12 : lineCountEnd]
				param["line_count"] = strings.TrimSpace(lineCountValue)
			}
		}
		return tagName, param, nil
	}

	// 对于其他工具，使用标准XML解析
	// 预处理XML内容，转义特殊字符
	processedContent := content
	processedContent = strings.ReplaceAll(processedContent, "&", "&amp;")

	decoder := xml.NewDecoder(strings.NewReader(processedContent))
	originalTagName := tagName
	tagName = ""

	for {
		tokenVal, tokenErr := decoder.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			}
			// 如果XML解析失败，但我们已经有了根元素名称，尝试使用手动解析
			if originalTagName != "" {
				return originalTagName, param, nil
			}
			return "", nil, tokenErr
		}

		switch se := tokenVal.(type) {
		case xml.StartElement:
			if tagName == "" {
				tagName = se.Name.Local
			} else {
				// 处理特殊情况：如果是follow_up元素，可能包含多个suggest子元素
				if se.Name.Local == "follow_up" {
					// 跳过follow_up标签的处理，因为我们已经在特殊处理部分处理了
					continue
				}

				var value string
				if err := decoder.DecodeElement(&value, &se); err != nil {
					return "", nil, err
				}

				// 特殊处理content字段，移除开头的换行符
				if se.Name.Local == "content" && len(value) > 0 && value[0] == '\n' {
					value = value[1:]
				}

				param[se.Name.Local] = value
			}
		}
	}

	// 如果XML解析没有找到标签名，但我们已经从字符串解析中获取了标签名
	if tagName == "" && originalTagName != "" {
		tagName = originalTagName
	}

	if tagName == "" {
		return "", nil, fmt.Errorf("invalid tool content")
	}

	return tagName, param, nil
}

func (s *Tools) Reload(ev event.Event, re event.Result) {
	workspaceVal := ev.Data()
	if workspaceVal == nil {
		log.Errorf("reload tools failed, workspace is nil ")
		if re != nil {
			re.Set(nil, cd.NewError(cd.IllegalParam, "invalid workspace"))
		}
		return
	}
	workspacePtr, ok := workspaceVal.(*models.Workspace)
	if !ok {
		log.Errorf("reload tools failed, workspace is invalid")
		if re != nil {
			re.Set(nil, cd.NewError(cd.IllegalParam, "invalid workspace"))
		}
		return
	}
	err := s.toolRegistry.Reload(workspacePtr)
	if re != nil {
		re.Set(nil, err)
	}
}

func (s *Tools) Execute(ev event.Event, re event.Result) {
	var err *cd.Error
	var executeParam imbtc.ToolParam
	var executeResult *imbtc.ToolResult
	var executeErr error
	defer func() {
		if re != nil {
			re.Set(executeResult, err)
			re.SetVal("param", executeParam)
		}
	}()

	dataVal := ev.Data()
	contentVal, contentOK := dataVal.(string)
	if !contentOK {
		err = cd.NewError(cd.IllegalParam, "execute tool failed, invalid tool content")
		log.Warnf("content is not string, dataVal:%v", dataVal)
		return
	}
	executeParam, executeResult, executeErr = s.execDataContent(contentVal)
	if executeErr != nil {
		err = cd.NewError(cd.Unexpected, fmt.Sprintf("execute tool failed, executeErr:%v", executeErr.Error()))
		log.Errorf("execute tool failed, executeErr:%v", executeErr)
		return
	}
}

func (s *Tools) RunCommand(workspaceName, workspaceDir, command string) (result string, err *cd.Error) {
	defer func() {
		if err != nil {
			log.Errorf("run command failed, err:%v", err)
		}
	}()

	itemsVal, itemErr := shlex.Split(strings.TrimSpace(command))
	if itemErr != nil {
		err = cd.NewError(cd.IllegalParam, itemErr.Error())
		return
	}
	command = itemsVal[0]
	var args []string
	if len(itemsVal) > 1 {
		args = itemsVal[1:]
	}
	commandEvent := event.NewEvent(imbcc.CommandRun, s.ID(), imbcc.CommanderModule, nil, workspaceName)
	paramPtr := &imbcc.CommanderParam{
		Command:      command,
		Args:         args,
		WorkspaceDir: workspaceDir,
	}
	commandEvent.SetData("command", paramPtr)

	eventResult := s.Base.SendEvent(commandEvent)
	if eventResult == nil {
		err = cd.NewError(cd.Unexpected, "run command failed, result is nil")
		return
	}
	commandResult, commandErr := eventResult.Get()
	if commandErr != nil {
		err = commandErr
		return
	}
	commandResultPtr := commandResult.(*imbcc.CommanderResult)
	if commandResultPtr.Workspace != workspaceName {
		err = cd.NewError(cd.Unexpected, "run command failed, workspace not match")
		return
	}
	stringBuff := &strings.Builder{}
	if commandResultPtr.StatusMessage != "" {
		stringBuff.WriteString(commandResultPtr.StatusMessage + "\n\n")
	}
	if commandResultPtr.ErrorMessage != "" {
		stringBuff.WriteString(commandResultPtr.ErrorMessage + "\n\n")
	}
	if commandResultPtr.ConsoleMessage != "" {
		stringBuff.WriteString(commandResultPtr.ConsoleMessage + "\n\n")
	}
	result = stringBuff.String()
	return
}

func (s *Tools) QueryTemplate(name string) (result *models.Template, err *cd.Error) {
	queryEvent := event.NewEvent(imbtmc.QueryTemplate, s.ID(), imbtmc.TemplatemanagerModule, nil, name)
	queryResult := s.SendEvent(queryEvent)
	templateVal, templateErr := queryResult.Get()
	if templateErr != nil {
		log.Errorf("queryTemplate failed, name:%s, error:%s", name, templateErr.Error())
		return
	}
	result = templateVal.(*models.Template)
	return
}
