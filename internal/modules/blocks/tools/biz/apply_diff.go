package biz

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type ApplyDiffTool struct{}

func (s *ApplyDiffTool) Name() string { return "apply_diff" }
func (s *ApplyDiffTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *ApplyDiffTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *ApplyDiffTool) Execute(workspaceName, workspaceDir, workspaceType, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	fullPath := workspaceDir
	if !fp.Exist(fullPath) {
		err = fmt.Errorf("workspace not exist: %s", fullPath)
		return
	}

	path, ok := param["path"].(string)
	if !ok || path == "" {
		err = fmt.Errorf("missing or invalid path parameter")
		return
	}

	diff, ok := param["diff"].(string)
	if !ok || diff == "" {
		err = fmt.Errorf("missing or invalid diff parameter")
		return
	}

	applyErr := ApplyDiff(fullPath, path, diff)
	if applyErr != nil {
		log.Errorf("ApplyDiff failed, fullPath:%s, path:%s err: %v", fullPath, path, applyErr)
		result = "apply diff failed!\n" + applyErr.Error()
		return
	}

	result = "apply diff success!"
	return
}

/*
1. 单个DiffBlock定义
```
<<<<<<< SEARCH
xxx
>>>>>>> REPLACE
```
之间的内容是一个DiffBlock

2. 单个DiffBlock包含的结构如下
```
替换位置信息
-------
原内容
=======
新内容
```
3. 替换位置信息的结构如下
```
:start_line:4
:end_line:5
```

```
:start_line:4
```

```
<line_count:33</line_count>
```
替换位置信息是可选，
*/
type DiffBlock struct {
	StartLine      int    `json:"start_line"`
	EndLine        int    `json:"end_line"`
	OriginContent  string `json:"origin_content"`
	ReplaceContent string `json:"replace_content"`
}

func extractNumbers(input string) []string {
	// 定义匹配数字的正则表达式
	re := regexp.MustCompile(`\d+`)
	// FindAllString 返回所有匹配的数字字符串
	return re.FindAllString(input, -1)
}

// 从diff中分解出各个待处理的diff块
// 每个diff块的格式如下
// <<<<<<< SEARCH
// :start_line:1
// :end_line:5
// -------
// [exact content to find including whitespace]
// =======
// [new content to replace with]
// >>>>>>> REPLACE
// 其中，SEARCH和REPLACE部分是必须的，SEARCH部分的内容必须完全匹配原始内容（包括空格和缩进）。
// diff 可能会包含多个DiffBlock，每个DiffBlock包含一个原始内容和一个替换内容。
// splitDiffBlock 函数的作用是将一个完整的diff字符串分解成多个DiffBlock，每个DiffBlock包含一个原始内容和一个替换内容。
// 并且校验分解出来的DiffBlock是否有效，包括行号是否正确、原始内容是否匹配等。
// 将分解出来的DiffBlock按照起始行号进行排序，确保后续处理的正确性。
func splitDiffBlock(diff string) ([]*DiffBlock, error) {
	var blocks []*DiffBlock
	remaining := diff

	if strings.TrimSpace(diff) == "" {
		return nil, nil
	}

	for {
		// 查找下一个SEARCH块的起始位置
		searchIdx := strings.Index(remaining, "<<<<<<< SEARCH")
		if searchIdx == -1 {
			break
		}

		// 提取块内容
		blockPart := remaining[searchIdx:]
		sepIdx := strings.Index(blockPart, "=======")
		replaceEndIdx := strings.Index(blockPart, ">>>>>>> REPLACE")

		if sepIdx == -1 || replaceEndIdx == -1 {
			return nil, fmt.Errorf("malformed diff block: missing separator or end marker")
		}

		// 提取元数据行
		metaPart := blockPart[len("<<<<<<< SEARCH"):]
		contentSep := strings.Index(metaPart, "-------")
		if contentSep == -1 {
			// 到这里说明是不包含替换位置信息
			contentSep = 0
		}
		metaLines := strings.Split(strings.TrimSpace(metaPart[:contentSep]), "\n")

		var startLine, endLine int
		startLineFound, endLineFound := false, false
		if len(metaLines) > 0 {
			// 到这里说明替换位置信息区有内容，这里尝试解析出对应的信息
			for _, line := range metaLines {
				if strings.HasPrefix(line, ":start_line:") {
					if _, err := fmt.Sscanf(line, ":start_line:%d", &startLine); err != nil {
						return nil, fmt.Errorf("invalid start line: %v", err)
					}
					startLineFound = true
				} else if strings.HasPrefix(line, ":end_line:") {
					if _, err := fmt.Sscanf(line, ":end_line:%d", &endLine); err != nil {
						return nil, fmt.Errorf("invalid end line: %v", err)
					}
					endLineFound = true
				}
			}

			// 由于DeepSeek返回过一下这种格式的数据
			// <line_count:33</line_count>
			// 这里尝试获取行数当成startLine
			numberList := []string{}
			for _, line := range metaLines {
				items := extractNumbers(line)
				numberList = append(numberList, items...)
			}
			if len(numberList) > 0 {
				startLine, _ = strconv.Atoi(numberList[0])
			}
			if len(numberList) > 1 {
				endLine, _ = strconv.Atoi(numberList[1])
			}
		}

		// 验证行号
		if !startLineFound {
			startLine = 0
		}

		// 提取原始内容和替换内容
		contentSep = strings.Index(blockPart, "-------")
		if contentSep == -1 {
			contentSep = len("<<<<<<< SEARCH")
		} else {
			contentSep += len("-------")
		}
		rawContent := strings.Trim(blockPart[contentSep:sepIdx], "\n")

		// 如果没有指定结束行号，根据原始内容的行数计算
		if !endLineFound && startLine > 0 {
			lineCount := len(strings.Split(rawContent, "\n"))
			endLine = startLine + lineCount - 1
		}

		if startLine > endLine {
			return nil, fmt.Errorf("start line cannot be greater than end line: startLine=%d, endLine=%d", startLine, endLine)
		}

		// 提取替换内容
		replaceContent := strings.Trim(blockPart[sepIdx+len("======="):replaceEndIdx], "\n")
		// 创建diff块
		blocks = append(blocks, &DiffBlock{
			StartLine:      startLine,
			EndLine:        endLine,
			OriginContent:  rawContent,
			ReplaceContent: replaceContent,
		})

		// 移动到下一个块
		remaining = blockPart[replaceEndIdx+len(">>>>>>> REPLACE"):]
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("no valid diff blocks found")
	}

	// 对块按起始行号排序
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].StartLine < blocks[j].StartLine
	})

	return blocks, nil
}

// ApplyDiff 实现apply_diff工具功能
func ApplyDiff(workspaceDir, path string, diff string) error {
	fullPath := filepath.Join(workspaceDir, path)
	fullPath = filepath.Clean(fullPath)
	if !fp.Exist(fullPath) {
		return fmt.Errorf("file not found: %s", path)
	}

	diffBlocks, err := splitDiffBlock(diff)
	if err != nil {
		return fmt.Errorf("diff块格式非法: %v", err)
	}
	if len(diffBlocks) == 0 {
		return nil
	}

	// 读取文件内容，用于在未指定行号时查找匹配内容
	fileContent, err := os.ReadFile(fullPath)
	if err != nil {
		log.Errorf("read raw file failed, error:%v", err)
		return fmt.Errorf("failed to read file:%s, err:%v", path, err)
	}
	fileLines := strings.Split(string(fileContent), "\n")

	// 处理未指定行号的情况，通过内容匹配确定行号
	for _, block := range diffBlocks {
		// 如果起始行号未指定，通过内容匹配确定
		if block.StartLine == 0 {
			// 在文件中查找匹配的内容
			lineNum := findContentInFile(fileLines, block.OriginContent)
			if lineNum == -1 {
				return fmt.Errorf("未在文件中找到diff块替换内容: \n%s", block.OriginContent+"\n\n 需要重新全量刷新该文件！\n\n")
			}
			block.StartLine = lineNum
			// 如果结束行号也未指定，根据原始内容的行数计算
			if block.EndLine == 0 {
				block.EndLine = lineNum + len(strings.Split(block.OriginContent, "\n")) - 1
			}
		} else if block.EndLine == 0 {
			// 如果只有结束行号未指定，根据原始内容的行数计算
			block.EndLine = block.StartLine + len(strings.Split(block.OriginContent, "\n")) - 1
		}

		// 验证行号的有效性
		if block.StartLine <= 0 || block.EndLine <= 0 || block.StartLine > block.EndLine {
			return fmt.Errorf("diff块行号无效: startLine=%d, endLine=%d，diff块%s",
				block.StartLine, block.EndLine, block.OriginContent)
		}
	}

	// 对块按起始行号重新排序
	sort.Slice(diffBlocks, func(i, j int) bool {
		return diffBlocks[i].StartLine < diffBlocks[j].StartLine
	})

	// 创建临时文件
	tmpPath := fullPath + ".tmp"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		log.Errorf("create temp file failed, error:%v", err)
		return fmt.Errorf("写入文件失败, error:%v", err)
	}
	defer tmpFile.Close()

	// 打开原文件
	origFile, err := os.Open(fullPath)
	if err != nil {
		log.Errorf("open file failed, error:%v", err)
		return fmt.Errorf("写入文件失败, error:%v", err)
	}
	defer origFile.Close()

	// 逐行读取并处理
	scanner := bufio.NewScanner(origFile)
	writer := bufio.NewWriter(tmpFile)
	lineNum := 1
	currentBlock := 0

	firstLine := true
	for scanner.Scan() {
		line := scanner.Text()
		// 检查是否需要替换
		if currentBlock < len(diffBlocks) {
			block := diffBlocks[currentBlock]

			// 如果当前行在替换范围内
			if lineNum >= block.StartLine && lineNum <= block.EndLine {
				// 第一行时写入替换内容
				if lineNum == block.StartLine {
					// 写入换行符
					if !firstLine {
						if _, err := writer.WriteString("\n"); err != nil {
							log.Errorf("write file failed, error:%v", err)
							return fmt.Errorf("写入文件失败, error:%v", err)
						}
					}
					firstLine = false
					if _, err := writer.WriteString(block.ReplaceContent); err != nil {
						log.Errorf("write file failed, error:%v", err)
						return fmt.Errorf("写入文件失败, error:%v", err)
					}
				}

				// 如果到达块的结束行，移动到下一个块
				if lineNum == block.EndLine {
					currentBlock++
				}

				lineNum++
				continue
			}
		}

		// 写入换行符
		if !firstLine {
			if _, err := writer.WriteString("\n"); err != nil {
				log.Errorf("write string failed, error:%v", err)
				return fmt.Errorf("写入文件失败, error:%v", err)
			}
		}
		firstLine = false
		// 不在替换范围内，写入原始行
		if _, err := writer.WriteString(line); err != nil {
			log.Errorf("write file failed, error:%v")
			return fmt.Errorf("写入文件失败, error:%v", err)
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		log.Errorf("scanner.Err() error:%v", err)
		return fmt.Errorf("写入文件失败, error:%v", err)
	}

	// 刷新缓冲区
	if err := writer.Flush(); err != nil {
		log.Errorf("flush buffer failed, error:%v", err)
		return fmt.Errorf("写入文件失败, error:%v", err)
	}

	// 关闭文件
	tmpFile.Close()
	origFile.Close()

	// 重命名临时文件
	if err := os.Rename(tmpPath, fullPath); err != nil {
		log.Errorf("os.rename failed, error:%v", err)
		return fmt.Errorf("写入文件失败, error:%v", err)
	}

	return nil
}

// findContentInFile 在文件内容中查找匹配的内容，返回匹配内容的起始行号（1-based）
// 如果找不到匹配的内容，返回-1
// 支持精确匹配和模糊匹配（忽略空白字符的差异）
func findContentInFile(fileLines []string, content string) int {
	if content == "" {
		return -1
	}

	contentLines := strings.Split(content, "\n")
	if len(contentLines) == 0 {
		return -1
	}

	// 单行内容直接查找
	if len(contentLines) == 1 {
		// 精确匹配
		for i, line := range fileLines {
			if line == contentLines[0] {
				return i + 1 // 返回1-based行号
			}
		}

		// 尝试忽略空白字符的匹配
		contentTrimmed := strings.TrimSpace(contentLines[0])
		if contentTrimmed != "" {
			for i, line := range fileLines {
				if strings.TrimSpace(line) == contentTrimmed {
					return i + 1 // 返回1-based行号
				}
			}
		}
		return -1
	}

	// 多行内容需要连续匹配
	// 1. 先尝试精确匹配
	for i := 0; i <= len(fileLines)-len(contentLines); i++ {
		matched := true
		for j := 0; j < len(contentLines); j++ {
			if i+j >= len(fileLines) || fileLines[i+j] != contentLines[j] {
				matched = false
				break
			}
		}
		if matched {
			return i + 1 // 返回1-based行号
		}
	}

	// 2. 如果精确匹配失败，尝试忽略空白字符的匹配
	// 预处理内容行，去除空白
	contentTrimmed := make([]string, len(contentLines))
	for i, line := range contentLines {
		contentTrimmed[i] = strings.TrimSpace(line)
	}

	for i := 0; i <= len(fileLines)-len(contentLines); i++ {
		matched := true
		for j := 0; j < len(contentLines); j++ {
			if i+j >= len(fileLines) || strings.TrimSpace(fileLines[i+j]) != contentTrimmed[j] {
				matched = false
				break
			}
		}
		if matched {
			return i + 1 // 返回1-based行号
		}
	}

	return -1
}
