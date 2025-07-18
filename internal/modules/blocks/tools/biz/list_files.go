package biz

import (
	"encoding/json"
	"fmt"
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

type ListFilesTool struct {
	workspacePtr *models.Workspace
}

func (s *ListFilesTool) Name() string { return "list_files" }
func (s *ListFilesTool) Prompt() string {
	promptsVal, promptsErr := imbtc.ReadInnerToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("failed to read prompts, error: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *ListFilesTool) Tags() fu.StringSet { return fu.StringSet{iipc.Global} }

func (s *ListFilesTool) Execute(param imbtc.ToolParam) (ret *imbtc.ToolResult, err error) {
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

	fullPath := s.workspacePtr.GetWorkingPath()
	if !fp.Exist(fullPath) {
		err = fmt.Errorf("workspace not exist: %s", fullPath)
		return
	}

	path, ok := param["path"].(string)
	if !ok || path == "" {
		err = fmt.Errorf("missing or invalid path parameter")
		return
	}
	recursive, _ := param["recursive"].(bool)
	listResult, listErr := ListFiles(fullPath, path, recursive)
	if listErr != nil {
		log.Errorf("ListFiles failed, fullPath:%s, path:%s error: %v", fullPath, path, listErr)
		result = fmt.Sprintf("list files failed, path:%s, error:%s", path, listErr.Error())
		return
	}

	byteVal, byteErr := json.MarshalIndent(listResult, "", "  ")
	if byteErr != nil {
		result = fmt.Sprintf("list files failed, path:%s, error:%s", path, byteErr.Error())
		return
	}

	result = fmt.Sprintf("list files successed, path:%s\n---\n%s", path, string(byteVal))
	return
}

/*
## list_files
## 列出文件

Description: Request to list files and directories within the specified directory. If recursive is true, it will list all files and directories recursively. If recursive is false or not provided, it will only list the top-level contents. Do not use this tool to confirm the existence of files you may have created, as the user will let you know if the files were created successfully or not.
描述：请求列出指定目录中的文件和目录。如果recursive为true，将递归列出所有文件和目录。如果recursive为false或未提供，则仅列出顶层内容。请勿使用此工具确认您可能已创建的文件是否存在，因为用户会告知您文件是否创建成功。

Parameters 参数:
- path: (required) The path of the directory to list contents for (relative to the current workspace directory /home/rangh/webspace)
- path: (必填) 要列出内容的目录路径（相对于当前工作空间目录 /home/rangh/webspace）
- recursive: (optional) Whether to list files recursively. Use true for recursive listing, false or omit for top-level only.
- recursive: (可选) 是否递归列出文件。使用true进行递归列出，false或省略则仅列出顶层内容。

Usage 用法:
<list_files>
<path>Directory path here 在此处填写目录路径</path>
<recursive>true or false (optional) true或false（可选）</recursive>
</list_files>

Example 示例: Requesting to list all files in the current directory
示例：请求列出当前目录中的所有文件
<list_files>
<path>.</path>
<recursive>false</recursive>
</list_files>
*/

// ListFiles 实现list_files工具功能
func ListFiles(workspaceDir, path string, recursive bool) ([]string, error) {
	fullPath := filepath.Join(workspaceDir, path)

	// 检查路径是否存在且是目录
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}

	var results []string

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对于fullPath的路径
		relPath, _ := filepath.Rel(fullPath, path)
		if relPath == "." {
			return nil // 跳过根目录自身
		}

		if info.IsDir() {
			// 对于目录，添加斜杠后缀
			results = append(results, filepath.Base(path)+"/")
		} else {
			// 对于文件，只添加文件名
			results = append(results, filepath.Base(path))
		}
		return nil
	}

	if recursive {
		err := filepath.Walk(fullPath, walkFunc)
		if err != nil {
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				results = append(results, entry.Name()+"/")
			} else {
				results = append(results, entry.Name())
			}
		}
	}

	return results, nil
}
