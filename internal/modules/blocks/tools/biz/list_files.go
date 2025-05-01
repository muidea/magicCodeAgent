package biz

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/muidea/magicCommon/foundation/log"

	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tools/pkg/common"
)

type ListFilesTool struct{}

func (s *ListFilesTool) Name() string { return "list_files" }
func (s *ListFilesTool) Prompts() string {
	promptsVal, promptsErr := imbtc.ReadToolsPrompts(s.Name())
	if promptsErr != nil {
		log.Warnf("read prompts failed, err: %s", promptsErr.Error())
		return ""
	}

	return promptsVal
}

func (s *ListFilesTool) Execute(param imbtc.ToolParam) (imbtc.ToolResult, error) {
	// 参数校验
	path, ok := param["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("缺少或无效的路径参数")
	}

	// 获取可选参数
	recursive, _ := param["recursive"].(bool)

	// 调用工具函数
	result, err := ListFiles(path, recursive)
	if err != nil {
		return nil, fmt.Errorf("列出文件失败: %v", err)
	}

	// 处理返回值
	ret, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("序列化结果失败: %v", err)
	}

	return ret, nil
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
func ListFiles(path string, recursive bool) ([]string, error) {
	fullPath := filepath.Join(WorkspaceDir, path)
	var results []string

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			results = append(results, path)
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
			results = append(results, filepath.Join(fullPath, entry.Name()))
		}
	}

	return results, nil
}
