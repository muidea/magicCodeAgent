package biz

import (
	"os"
	"path/filepath"
	"regexp"
)

/*
## search_files
Description: Request to perform a regex search across files in a specified directory, providing context-rich results. This tool searches for patterns or specific content across multiple files, displaying each match with encapsulating context.
Parameters:
- path: (required) The path of the directory to search in (relative to the current workspace directory /home/rangh/webspace). This directory will be recursively searched.
- regex: (required) The regular expression pattern to search for. Uses Rust regex syntax.
- file_pattern: (optional) Glob pattern to filter files (e.g., '*.ts' for TypeScript files). If not provided, it will search all files (*).
Usage:
<search_files>
<path>Directory path here</path>
<regex>Your regex pattern here</regex>
<file_pattern>file pattern here (optional)</file_pattern>
</search_files>

Example: Requesting to search for all .ts files in the current directory
<search_files>
<path>.</path>
<regex>.*</regex>
<file_pattern>*.ts</file_pattern>
</search_files>

*/

// SearchFiles 实现search_files工具功能
// SearchFiles function implements the search_files tool functionality
func (s *Tool) SearchFiles(path, regexPattern, filePattern string) ([]string, error) {
	fullPath := filepath.Join(WorkspaceDir, path)
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
