package biz

import (
	"os"
	"path/filepath"
)

/*
## list_files
Description: Request to list files and directories within the specified directory. If recursive is true, it will list all files and directories recursively. If recursive is false or not provided, it will only list the top-level contents. Do not use this tool to confirm the existence of files you may have created, as the user will let you know if the files were created successfully or not.
Parameters:
- path: (required) The path of the directory to list contents for (relative to the current workspace directory /home/rangh/webspace)
- recursive: (optional) Whether to list files recursively. Use true for recursive listing, false or omit for top-level only.
Usage:
<list_files>
<path>Directory path here</path>
<recursive>true or false (optional)</recursive>
</list_files>

Example: Requesting to list all files in the current directory
<list_files>
<path>.</path>
<recursive>false</recursive>
</list_files>
*/

// ListFiles 实现list_files工具功能
func (s *Tool) ListFiles(path string, recursive bool) ([]string, error) {
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
