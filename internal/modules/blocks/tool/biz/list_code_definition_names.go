package biz

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

/*
## list_code_definition_names
Description: Request to list definition names (classes, functions, methods, etc.) from source code. This tool can analyze either a single file or all files at the top level of a specified directory. It provides insights into the codebase structure and important constructs, encapsulating high-level concepts and relationships that are crucial for understanding the overall architecture.
Parameters:
- path: (required) The path of the file or directory (relative to the current working directory /home/rangh/webspace) to analyze. When given a directory, it lists definitions from all top-level source files.
Usage:
<list_code_definition_names>
<path>Directory path here</path>
</list_code_definition_names>

Examples:

1. List definitions from a specific file:
<list_code_definition_names>
<path>src/main.ts</path>
</list_code_definition_names>

2. List definitions from all files in a directory:
<list_code_definition_names>
<path>src/</path>
</list_code_definition_names>

*/

  // ListCodeDefinitionNames 实现list_code_definition_names工具功能
  // ListCodeDefinitionNames function implements the list_code_definition_names tool functionality
  func (s *Tool) ListCodeDefinitionNames(path string) ([]string, error) {
  	fullPath := filepath.Join(WorkspaceDir, path)
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
