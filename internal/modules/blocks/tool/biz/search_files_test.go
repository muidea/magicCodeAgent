package biz

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchFiles(t *testing.T) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "search_files_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 保存原始的WorkspaceDir并在测试结束后恢复
	originalWorkspaceDir := WorkspaceDir
	defer func() { WorkspaceDir = originalWorkspaceDir }()
	// 设置测试的WorkspaceDir为临时目录
	WorkspaceDir = tempDir

	// 创建测试文件结构
	testFiles := map[string]string{
		"file1.txt":        "This is a test file with some content",
		"file2.txt":        "Another test file with different content",
		"file3.log":        "This is a log file",
		"subdir/file4.txt": "A file in subdirectory",
	}

	// 创建测试文件
	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		// 确保目录存在
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		// 写入文件内容
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}

	// 创建测试工具实例
	tool := &Tool{}

	// 测试用例
	tests := []struct {
		name        string
		path        string
		regex       string
		filePattern string
		expected    int // 期望匹配的文件数量
		isError     bool
	}{
		{
			name:        "匹配所有txt文件",
			path:        ".",
			regex:       ".*",
			filePattern: "*.txt",
			expected:    3,
			isError:     false,
		},
		{
			name:        "匹配包含test的文件",
			path:        ".",
			regex:       "test",
			filePattern: "",
			expected:    2,
			isError:     false,
		},
		{
			name:        "匹配log文件",
			path:        ".",
			regex:       ".*",
			filePattern: "*.log",
			expected:    1,
			isError:     false,
		},
		{
			name:        "递归搜索子目录",
			path:        ".",
			regex:       "subdir",
			filePattern: "",
			expected:    0, // 文件内容中没有"subdir"字样
			isError:     false,
		},
		{
			name:        "搜索子目录中的文件",
			path:        ".",
			regex:       "A file in subdirectory",
			filePattern: "",
			expected:    1,
			isError:     false,
		},
		{
			name:        "无效的正则表达式",
			path:        ".",
			regex:       "[", // 无效的正则表达式
			filePattern: "",
			expected:    0,
			isError:     true,
		},
		{
			name:        "无效的路径",
			path:        "non_existent_dir",
			regex:       ".*",
			filePattern: "",
			expected:    0,
			isError:     true,
		},
		{
			name:        "无匹配文件",
			path:        ".",
			regex:       "nonexistent content",
			filePattern: "",
			expected:    0,
			isError:     false,
		},
		{
			name:        "无效的文件模式",
			path:        ".",
			regex:       ".*",
			filePattern: "[*.invalid", // 无效的文件模式
			expected:    0,
			isError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := tool.SearchFiles(tt.path, tt.regex, tt.filePattern)

			if tt.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, len(results), "Expected %d files, got %d", tt.expected, len(results))
			}
		})
	}
}
