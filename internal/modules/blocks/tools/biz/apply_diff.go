package biz

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
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

func (s *ApplyDiffTool) Execute(workspaceDir string, param imbtc.ToolParam) (imbtc.ToolResult, error) {
	// 参数校验
	path, ok := param["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("missing or invalid file path parameter")
	}

	diff, ok := param["diff"].(string)
	if !ok || diff == "" {
		return nil, fmt.Errorf("missing or invalid diff parameter")
	}

	// 调用工具函数
	err := ApplyDiff(workspaceDir, path, diff)
	if err != nil {
		return nil, fmt.Errorf("failed to apply diff: %v", err)
	}

	// 处理返回值
	result, err := json.Marshal("差异应用成功")
	if err != nil {
		return nil, fmt.Errorf("failed to serialize result: %v", err)
	}
	return result, nil
}

/*
## apply_diff
Description: Request to replace existing code using a search and replace block.
This tool allows for precise, surgical replaces to files by specifying exactly what content to search for and what to replace it with.
The tool will maintain proper indentation and formatting while making changes.
Only a single operation is allowed per tool use.
The SEARCH section must exactly match existing content including whitespace and indentation.
If you're not confident in the exact content to search for, use the read_file tool first to get the exact content.
When applying the diffs, be extra careful to remember to change any closing brackets or other syntax that may be affected by the diff farther down in the file.
ALWAYS make as many changes in a single 'apply_diff' request as possible using multiple SEARCH/REPLACE blocks

描述：使用搜索替换块请求替换现有代码。该工具通过精确指定要搜索的内容和替换内容，实现文件的精准修改。工具在更改时会保持正确的缩进和格式。每次工具使用仅允许单个操作。SEARCH部分必须完全匹配现有内容（包括空格和缩进）。如果不确定要搜索的确切内容，请先使用read_file工具获取准确内容。应用差异时要特别注意修改可能受影响的闭合括号或其他语法。尽量在单个apply_diff请求中使用多个SEARCH/REPLACE块完成所有更改。

Parameters:
- path: (required) The path of the file to modify (relative to the current workspace directory /home/rangh/webspace)
  路径：(必填) 要修改的文件路径（相对于当前工作空间目录/home/rangh/webspace）
- diff: (required) The search/replace block defining the changes.
  差异内容：(必填) 定义更改的搜索/替换块

Diff format:
```
<<<<<<< SEARCH
:start_line: (required) The line number of original content where the search block starts.
:end_line: (required) The line number of original content  where the search block ends.
-------
[exact content to find including whitespace]
=======
[new content to replace with]
>>>>>>> REPLACE

```


Example:
示例：

Original file:
```
1 | def calculate_total(items):
2 |     total = 0
3 |     for item in items:
4 |         total += item
5 |     return total
```

Search/Replace content:
```
<<<<<<< SEARCH
:start_line:1
:end_line:5
-------
def calculate_total(items):
    total = 0
    for item in items:
        total += item
    return total
=======
def calculate_total(items):
    \"\"\"Calculate total with 10% markup\"\"\"
    return sum(item * 1.1 for item in items)
>>>>>>> REPLACE

```

Search/Replace content with multi edits:
```
<<<<<<< SEARCH
:start_line:1
:end_line:2
-------
def calculate_total(items):
    sum = 0
=======
def calculate_sum(items):
    sum = 0
>>>>>>> REPLACE

<<<<<<< SEARCH
:start_line:4
:end_line:5
-------
        total += item
    return total
=======
        sum += item
    return sum
>>>>>>> REPLACE
```


Usage:
<apply_diff>
<path>File path here</path>
<diff>
Your search/replace content here
You can use multi search/replace block in one diff block, but make sure to include the line numbers for each block.
Only use a single line of '=======' between search and replacement content, because multiple '=======' will corrupt the file.
</diff>
</apply_diff>

*/
// ApplyDiff 实现apply_diff工具功能

type DiffBlock struct {
	StartLine      int    `json:"start_line"`
	EndLine        int    `json:"end_line"`
	OriginContent  string `json:"origin_content"`
	ReplaceContent string `json:"replace_content"`
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
			return nil, fmt.Errorf("malformed diff block: missing content separator")
		}
		metaLines := strings.Split(strings.TrimSpace(metaPart[:contentSep]), "\n")

		var startLine, endLine int
		startLineFound, endLineFound := false, false
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

		// 验证行号
		if !startLineFound {
			return nil, fmt.Errorf("missing start line")
		}

		if startLine <= 0 {
			return nil, fmt.Errorf("start line must be positive: startLine=%d", startLine)
		}

		// 提取原始内容和替换内容
		contentSep = strings.Index(blockPart, "-------")
		rawContent := strings.Trim(blockPart[contentSep+len("-------"):sepIdx], "\n")

		// 如果没有指定结束行号，根据原始内容的行数计算
		if !endLineFound {
			lineCount := len(strings.Split(rawContent, "\n"))
			endLine = startLine + lineCount - 1
			log.Infof("未指定结束行号，根据内容自动计算: startLine=%d, lineCount=%d, endLine=%d", startLine, lineCount, endLine)
		} else if endLine <= 0 {
			return nil, fmt.Errorf("end line must be positive: endLine=%d", endLine)
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
		return fmt.Errorf("file not found: %s", fullPath)
	}

	diffBlocks, err := splitDiffBlock(diff)
	if err != nil {
		return fmt.Errorf("failed to parse diff: %v", err)
	}
	if len(diffBlocks) == 0 {
		return nil
	}

	// 创建临时文件
	tmpPath := fullPath + ".tmp"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	// 打开原文件
	origFile, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("failed to open original file: %v", err)
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
							return fmt.Errorf("failed to write newline: %v", err)
						}
					}
					firstLine = false
					if _, err := writer.WriteString(block.ReplaceContent); err != nil {
						return fmt.Errorf("failed to write replacement content: %v", err)
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
				return fmt.Errorf("failed to write newline: %v", err)
			}
		}
		firstLine = false
		// 不在替换范围内，写入原始行
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("failed to write original content: %v", err)
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// 刷新缓冲区
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush write buffer: %v", err)
	}

	// 关闭文件
	tmpFile.Close()
	origFile.Close()

	// 重命名临时文件
	if err := os.Rename(tmpPath, fullPath); err != nil {
		return fmt.Errorf("failed to rename temp file: %v", err)
	}

	return nil
}
