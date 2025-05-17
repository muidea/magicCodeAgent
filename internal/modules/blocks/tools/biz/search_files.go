package biz

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type SearchFilesTool struct{}

func (s *SearchFilesTool) Name() string { return "search_files" }
func (s *SearchFilesTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *SearchFilesTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *SearchFilesTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult) {
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

	path, ok := param["path"].(string)
	if !ok || path == "" {
		err = fmt.Errorf("missing or invalid path parameter")
		return
	}

	regex, ok := param["regex"].(string)
	if !ok || regex == "" {
		err = fmt.Errorf("missing or invalid regex parameter")
		return
	}

	// 获取可选参数
	filePattern := "*"
	if filePatternParam, ok := param["file_pattern"].(string); ok && filePatternParam != "" {
		filePattern = filePatternParam
	}
	fullPath := filepath.Join(workspaceDir, stage)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}
	searchResult, searchErr := SearchFiles(fullPath, path, regex, filePattern)
	if searchErr != nil {
		err = searchErr
		log.Errorf("SearchFiles failed, fullPath:%s, path:%s err: %v", fullPath, path, err)
		return
	}

	byteVal, byteErr := json.MarshalIndent(searchResult, "", "  ")
	if byteErr != nil {
		err = byteErr
		return
	}

	result = string(byteVal)
	return
}

/*
## search_files
## 搜索文件

Description: Request to perform a regex search across files in a specified directory, providing context-rich results. This tool searches for patterns or specific content across multiple files, displaying each match with encapsulating context.
描述：在指定目录中执行正则表达式搜索，提供带上下文的结果。此工具可以跨多个文件搜索模式或特定内容，显示每个匹配项及其上下文。

Parameters 参数:
- path: (required) The path of the directory to search in (relative to the current workspace directory /home/rangh/webspace). This directory will be recursively searched.
- path: (必填) 要搜索的目录路径（相对于当前工作空间目录 /home/rangh/webspace）。将递归搜索此目录。

- regex: (required) The regular expression pattern to search for. Uses Rust regex syntax.
- regex: (必填) 要搜索的正则表达式模式。使用Rust正则表达式语法。

- file_pattern: (optional) Glob pattern to filter files (e.g., '*.ts' for TypeScript files). If not provided, it will search all files (*).
- file_pattern: (可选) 用于过滤文件的通配符模式（例如，'*.ts'用于TypeScript文件）。如果未提供，将搜索所有文件(*)。

Usage 用法:
<search_files>
<path>Directory path here 在此处填写目录路径</path>
<regex>Your regex pattern here 在此处填写正则表达式模式</regex>
<file_pattern>file pattern here (optional) 在此处填写文件模式（可选）</file_pattern>
</search_files>

Example 示例: Requesting to search for all .ts files in the current directory
示例：请求搜索当前目录中的所有.ts文件
<search_files>
<path>.</path>
<regex>.*</regex>
<file_pattern>*.ts</file_pattern>
</search_files>
*/

// SearchFiles 实现search_files工具功能
// SearchFiles function implements the search_files tool functionality
func SearchFiles(workspaceDir, path, regexPattern, filePattern string) ([]string, error) {
	fullPath := filepath.Join(workspaceDir, path)
	var results []string

	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if filePattern != "" {
				matched, _ := filepath.Match(filePattern, filepath.Base(path))
				if !matched {
					return nil
				}
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			regex, err := regexp.Compile(regexPattern)
			if err != nil {
				return err
			}

			if regex.Match(content) {
				results = append(results, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
