package biz

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type ReadFile struct{}

func (s *ReadFile) Name() string {
	return "read_file"
}

func (s *ReadFile) Prompts() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *ReadFile) Execute(param imbtc.ToolParam) (err error) {
	return
}

/*
## read_file
Description: Request to read the contents of a file at the specified path. Use this when you need to examine the contents of an existing file you do not know the contents of, for example to analyze code, review text files, or extract information from configuration files. The output includes line numbers prefixed to each line (e.g. \"1 | const x = 1\"), making it easier to reference specific lines when creating diffs or discussing code. By specifying start_line and end_line parameters, you can efficiently read specific portions of large files without loading the entire file into memory. Automatically extracts raw text from PDF and DOCX files. May not be suitable for other types of binary files, as it returns the raw content as a string.
描述：请求读取指定路径的文件内容。当您需要检查一个现有文件的内容但不知道其内容时使用此功能，例如分析代码、查看文本文件或从配置文件中提取信息。输出包括每行前缀的行号（例如"1 | const x = 1"），使得在创建差异或讨论代码时更容易引用特定行。通过指定start_line和end_line参数，您可以高效地读取大文件的特定部分而无需将整个文件加载到内存中。自动从PDF和DOCX文件中提取原始文本。可能不适合其他类型的二进制文件，因为它返回原始内容作为字符串。
Parameters:
- path: (required) The path of the file to read (relative to the current workspace directory /home/rangh/webspace)
  路径：(必填) 要读取的文件路径（相对于当前工作空间目录/home/rangh/webspace）
- start_line: (optional) The starting line number to read from (1-based). If not provided, it starts from the beginning of the file.
  起始行：(可选) 开始读取的行号(从1开始)。如果未提供，则从文件开头开始。
- end_line: (optional) The ending line number to read to (1-based, inclusive). If not provided, it reads to the end of the file.
  结束行：(可选) 结束读取的行号(从1开始，包含该行)。如果未提供，则读取到文件末尾。
Usage:
<read_file>
<path>File path here</path>
<start_line>Starting line number (optional)</start_line>
<end_line>Ending line number (optional)</end_line>
</read_file>

Examples:

1. Reading an entire file:
<read_file>
<path>frontend-config.json</path>
</read_file>

2. Reading the first 1000 lines of a large log file:
<read_file>
<path>logs/application.log</path>
<end_line>1000</end_line>
</read_file>

3. Reading lines 500-1000 of a CSV file:
<read_file>
<path>data/large-dataset.csv</path>
<start_line>500</start_line>
<end_line>1000</end_line>
</read_file>

4. Reading a specific function in a source file:
<read_file>
<path>src/app.ts</path>
<start_line>46</start_line>
<end_line>68</end_line>
</read_file>

Note: When both start_line and end_line are provided, this tool efficiently streams only the requested lines, making it suitable for processing large files like logs, CSV files, and other large datasets without memory issues.
注意：当同时提供start_line和end_line时，此工具会高效地仅流式传输请求的行，使其适合处理大型文件（如日志、CSV文件和其他大型数据集）而不会出现内存问题。
*/
// ReadFile 读取文件内容
func (s *Tools) ReadFile(path string, startLine, endLine int) (string, error) {
	fullPath := filepath.Join(WorkspaceDir, path)
	file, err := os.Open(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var result strings.Builder
	scanner := bufio.NewScanner(file)
	lineNum := 1

	// 跳过直到开始行
	for lineNum < startLine && scanner.Scan() {
		lineNum++
	}

	// 读取直到结束行
	for (endLine <= 0 || lineNum <= endLine) && scanner.Scan() {
		line := scanner.Text()
		result.WriteString(fmt.Sprintf("%d | %s\n", lineNum, line))
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	return result.String(), nil
}
