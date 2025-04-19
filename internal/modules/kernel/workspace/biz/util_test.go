package biz

import (
	"testing"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

// TestConvertPathsToTree 测试路径转换功能
func TestConvertPathsToTree(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		expected []*common.DataNode
	}{
		{
			name:     "空路径列表",
			paths:    []string{},
			expected: []*common.DataNode{},
		},
		{
			name:  "单层路径",
			paths: []string{"file1", "file2"},
			expected: []*common.DataNode{
				{Name: "file1", Path: "file1", IsDir: false},
				{Name: "file2", Path: "file2", IsDir: false},
			},
		},
		{
			name:  "多层嵌套路径",
			paths: []string{"dir1/file1", "dir1/file2", "dir2/sub/file3"},
			expected: []*common.DataNode{
				{
					Name:  "dir1",
					Path:  "dir1",
					IsDir: true,
					Children: []*common.DataNode{
						{Name: "file1", Path: "dir1/file1", IsDir: false},
						{Name: "file2", Path: "dir1/file2", IsDir: false},
					},
				},
				{
					Name:  "dir2",
					Path:  "dir2",
					IsDir: true,
					Children: []*common.DataNode{
						{
							Name:  "sub",
							Path:  "dir2/sub",
							IsDir: true,
							Children: []*common.DataNode{
								{Name: "file3", Path: "dir2/sub/file3", IsDir: false},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertPathsToTree("", tt.paths)
			if !compareTree(got, tt.expected) {
				t.Errorf("ConvertPathsToTree() = %v, want %v", got, tt.expected)
			}
		})
	}

	tests = []struct {
		name     string
		paths    []string
		expected []*common.DataNode
	}{
		{
			name:     "空路径列表",
			paths:    []string{"/home/test"},
			expected: []*common.DataNode{},
		},
		{
			name:  "单层路径",
			paths: []string{"/home/test/file1", "/home/test/file2"},
			expected: []*common.DataNode{
				{Name: "file1", Path: "file1", IsDir: false},
				{Name: "file2", Path: "file2", IsDir: false},
			},
		},
		{
			name:  "多层嵌套路径",
			paths: []string{"/home/test/dir1/file1", "/home/test/dir1/file2", "/home/test/dir2/sub/file3"},
			expected: []*common.DataNode{
				{
					Name:  "dir1",
					Path:  "dir1",
					IsDir: true,
					Children: []*common.DataNode{
						{Name: "file1", Path: "dir1/file1", IsDir: false},
						{Name: "file2", Path: "dir1/file2", IsDir: false},
					},
				},
				{
					Name:  "dir2",
					Path:  "dir2",
					IsDir: true,
					Children: []*common.DataNode{
						{
							Name:  "sub",
							Path:  "dir2/sub",
							IsDir: true,
							Children: []*common.DataNode{
								{Name: "file3", Path: "dir2/sub/file3", IsDir: false},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertPathsToTree("/home/test", tt.paths)
			if !compareTree(got, tt.expected) {
				t.Errorf("ConvertPathsToTree() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func compareTree(a, b []*common.DataNode) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name || a[i].Path != b[i].Path || a[i].IsDir != b[i].IsDir {
			return false
		}
		if !compareTree(a[i].Children, b[i].Children) {
			return false
		}
	}
	return true
}
