package biz

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/task"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/config"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCodeAgent/pkg/util"
)

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
		Base:         biz.New(imbtc.ToolsModule, eventHub, backgroundRoutine),
		toolRegistry: NewToolRegistry(promptsRegistry),
	}

	ptr.toolRegistry.Register(&ReadFileTool{})
	ptr.toolRegistry.Register(&FetchInstructionsTool{})
	ptr.toolRegistry.Register(&SearchFilesTool{})
	ptr.toolRegistry.Register(&ListFilesTool{})
	ptr.toolRegistry.Register(&ListCodeDefinitionNamesTool{})
	ptr.toolRegistry.Register(&ApplyDiffTool{})
	ptr.toolRegistry.Register(&WriteToFileTool{})
	ptr.toolRegistry.Register(&InsertContentTool{})
	ptr.toolRegistry.Register(&SearchAndReplaceTool{})
	ptr.toolRegistry.Register(&ExecuteCommandTool{})
	ptr.toolRegistry.Register(&UseMcpToolTool{})
	ptr.toolRegistry.Register(&AccessMcpResourceTool{})
	ptr.toolRegistry.Register(&AskFollowupQuestionTool{})
	ptr.toolRegistry.Register(&AttemptCompletionTool{})
	ptr.toolRegistry.Register(&SwitchModeTool{})
	ptr.toolRegistry.Register(&NewTaskTool{})

	ptr.SubscribeFunc(imbtc.ExecuteTool, ptr.Execute)
	return ptr
}

func (s *Tools) Run() (err *cd.Error) {
	currentWorkDir := fu.GetCurrentWorkDir()

	dataFile := config.GetDataFile()
	if dataFile != "" {
		resultVal, resultErr := s.execDataFile(currentWorkDir, dataFile)
		if resultErr != nil {
			log.Errorf("exec data file failed!\n%v", resultErr.Error())
			return
		}

		log.Infof("exec data file result.\n%s", resultVal)
		return
	}

	dataContent := config.GetDataContent()
	if dataContent != "" {
		_, resultVal := s.execDataContent(currentWorkDir, dataContent)
		if resultVal != nil {
			log.Infof("exec data content result.\n%v", resultVal.Result.String())
		}
		return
	}

	return
}

func (s *Tools) execDataFile(workspaceDir, dataFile string) (string, error) {
	if !fp.Exist(dataFile) {
		log.Errorf("data file does not exist: %s", dataFile)
		return "", fmt.Errorf("data file does not exist: %s", dataFile)
	}

	byteContent, byteErr := os.ReadFile(dataFile)
	if byteErr != nil {
		log.Errorf("failed to read data file: %s", byteErr.Error())
		return "", fmt.Errorf("failed to read data file: %s", byteErr.Error())
	}

	blocksList, blockErr := util.SpliteBlocks(string(byteContent))
	if blockErr != nil {
		log.Errorf("failed to splite blocks: %s", blockErr.Error())
		return "", fmt.Errorf("failed to splite blocks: %s", blockErr.Error())
	}

	byteBuilder := strings.Builder{}
	for idx, block := range blocksList {
		_, toolResult := s.execDataContent(workspaceDir, block)
		if toolResult != nil {
			if idx > 0 {
				byteBuilder.WriteString("=======\n")
			}
			byteBuilder.WriteString(toolResult.Result.String() + "\n")
		}
	}
	return byteBuilder.String(), nil
}

func (s *Tools) execDataContent(workspaceDir, dataContent string) (string, *imbtc.ToolResult) {
	toolName, toolParam, toolErr := parseParamContent(dataContent)
	if toolErr != nil {
		log.Errorf("failed to parse param content: %s", toolErr.Error())
		return "", &imbtc.ToolResult{
			Final: true,
			Result: models.MessageItem{
				Content: []string{fmt.Sprintf("failed to parse param content: %s", toolErr.Error())},
			},
		}
	}
	if toolName == "" {
		log.Errorf("tool name cannot be empty")
		return "", &imbtc.ToolResult{
			Final: true,
			Result: models.MessageItem{
				Content: []string{"tool name cannot be empty"},
			},
		}
	}

	if toolName == "thinking" {
		return "thinking", &imbtc.ToolResult{
			Final: true,
			Result: models.MessageItem{
				Content: []string{"thinking ok!"},
			},
		}
	}

	return toolName, s.toolRegistry.Execute(workspaceDir, toolName, toolParam)
}

// 解析执行参数
// content里解析出执行函数名和参数
// content的格式是xml
// 执行函数是xml的root节点名
// 执行函数的参数是root的子节点及子节点的值
// 参照示例参见各个工具函数对应的注释信息
func parseParamContent(content string) (string, map[string]any, error) {
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
	param := make(map[string]any)

	// 特殊处理包含特殊XML字符的工具
	if tagName == "apply_diff" || tagName == "write_to_file" {
		// 提取path参数
		pathStart := strings.Index(content, "<path>")
		pathEnd := strings.Index(content, "</path>")
		if pathStart != -1 && pathEnd != -1 {
			pathValue := content[pathStart+6 : pathEnd]
			param["path"] = strings.TrimSpace(pathValue)
		}

		// 根据工具类型提取不同参数
		if tagName == "apply_diff" {
			// 提取diff参数
			diffStart := strings.Index(content, "<diff>")
			diffEnd := strings.Index(content, "</diff>")
			if diffStart != -1 && diffEnd != -1 {
				diffValue := content[diffStart+6 : diffEnd]
				param["diff"] = diffValue
			}
		} else if tagName == "write_to_file" {
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

func (s *Tools) Execute(ev event.Event, re event.Result) {
	var err *cd.Error
	var executeToolName string
	var executeResult *imbtc.ToolResult
	defer func() {
		if re != nil {
			re.Set(executeResult, err)
			re.SetVal("toolName", executeToolName)
		}
	}()

	dataVal := ev.Data()
	contentVal, contentOK := dataVal.(string)
	if !contentOK {
		err = cd.NewError(cd.IllegalParam, "execute tool failed, invalid tool content")
		log.Warnf("content is not string, dataVal:%v", dataVal)
		return
	}
	workspaceDirVal := ev.GetData("workspaceDir")
	workspaceDir, workspaceDirOK := workspaceDirVal.(string)
	if !workspaceDirOK {
		err = cd.NewError(cd.IllegalParam, "execute tool failed, workspaceDir is not string")
		log.Warnf("workspaceDir is not string, dataVal:%v", workspaceDirVal)
		return
	}

	executeToolName, executeResult = s.execDataContent(workspaceDir, contentVal)
}
