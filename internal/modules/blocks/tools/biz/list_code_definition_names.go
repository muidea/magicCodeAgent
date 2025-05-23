package biz

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initiators/prompts/pkg/common"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type ListCodeDefinitionNamesTool struct{}

func (s *ListCodeDefinitionNamesTool) Name() string { return "list_code_definition_names" }
func (s *ListCodeDefinitionNamesTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *ListCodeDefinitionNamesTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *ListCodeDefinitionNamesTool) Execute(workspaceName, workspaceDir, stage string, param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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
	fullPath := filepath.Join(workspaceDir, stage)
	if !fp.Exist(fullPath) {
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			return
		}
	}
	listResult, listErr := ListCodeDefinitionNames(fullPath, path)
	if listErr != nil {
		err = listErr
		log.Errorf("failed to list code definitions: %v", err)
		return
	}
	byteVal, byteErr := json.MarshalIndent(listResult, "", "  ")
	if byteErr != nil {
		err = byteErr
		return
	}

	result = string(byteVal)
	return
}

/*
## list_code_definition_names
## 列出代码定义名称

Description: Request to list definition names (classes, functions, methods, etc.) from source code. This tool can analyze either a single file or all files at the top level of a specified directory. It provides insights into the codebase structure and important constructs, encapsulating high-level concepts and relationships that are crucial for understanding the overall architecture.
描述：请求列出源代码中的定义名称（类、函数、方法等）。此工具可以分析单个文件或指定目录顶层的所有文件。它提供了代码库结构和重要构造的洞察，封装了对理解整体架构至关重要的高级概念和关系。

Parameters:
参数：
- path: (required) The path of the file or directory (relative to the current working directory /home/rangh/webspace) to analyze. When given a directory, it lists definitions from all top-level source files.
- path: (必需) 要分析的文件或目录的路径（相对于当前工作目录 /home/rangh/webspace）。当给定目录时，它会列出所有顶层源文件中的定义。

Usage:
用法：
<list_code_definition_names>
<path>Directory path here</path>
</list_code_definition_names>

Examples:
示例：

1. List definitions from a specific file:
1. 列出特定文件中的定义：
<list_code_definition_names>
<path>src/main.ts</path>
</list_code_definition_names>

2. List definitions from all files in a directory:
2. 列出目录中所有文件的定义：
<list_code_definition_names>
<path>src/</path>
</list_code_definition_names>
*/

// ListCodeDefinitionNames 实现list_code_definition_names工具功能
// ListCodeDefinitionNames function implements the list_code_definition_names tool functionality
func ListCodeDefinitionNames(workspaceDir, path string) ([]string, error) {
	fullPath := filepath.Join(workspaceDir, path)
	var definitions []string

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	if info.IsDir() {
		// 处理目录
		pkgs, err := parser.ParseDir(fset, fullPath, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				ast.Inspect(file, func(n ast.Node) bool {
					switch x := n.(type) {
					case *ast.FuncDecl:
						definitions = append(definitions, x.Name.Name)
					case *ast.TypeSpec:
						definitions = append(definitions, x.Name.Name)
					case *ast.ValueSpec:
						for _, name := range x.Names {
							definitions = append(definitions, name.Name)
						}
					}
					return true
				})
			}
		}
	} else {
		// 处理单个文件
		file, err := parser.ParseFile(fset, fullPath, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				definitions = append(definitions, x.Name.Name)
			case *ast.TypeSpec:
				definitions = append(definitions, x.Name.Name)
			case *ast.ValueSpec:
				for _, name := range x.Names {
					definitions = append(definitions, name.Name)
				}
			}
			return true
		})
	}

	return definitions, nil
}
