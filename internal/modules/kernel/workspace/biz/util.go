package biz

import (
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

type HybridEvent struct {
	EventType string // "filesystem" 或 "git"
	Path      string
	Operation string // create/update/delete 或 commit
	Timestamp time.Time
}

func (s HybridEvent) String() string {
	return s.EventType + " " + s.Operation + " " + s.Path + " " + s.Timestamp.String()
}

type EventHandler func(HybridEvent)

// 判断是否忽略
// .git .tmp 隐藏文件都忽略，不用处理
func isIgnored(path string) bool {
	base := filepath.Base(path)
	return strings.Contains(path, ".git") ||
		strings.HasSuffix(path, ".tmp") ||
		strings.HasPrefix(base, ".")
}

func NewParentPath(rootPath string, parentPath string) *parentPathFilter {
	return &parentPathFilter{
		parentPath: filepath.Clean(path.Join(rootPath, parentPath)),
	}
}

type parentPathFilter struct {
	parentPath string
}

func (s *parentPathFilter) Filter(val interface{}) bool {
	if val == nil {
		return false
	}

	pathVal, pathOK := val.(string)
	if !pathOK {
		return false
	}

	return strings.HasPrefix(pathVal, s.parentPath)
}

// ConvertPathsToTree 将文件路径列表转换为树形结构
func ConvertPathsToTree(basePath string, paths []string) []*common.DataNode {
	root := &common.DataNode{Name: "", Path: "", IsDir: true}
	for _, p := range paths {
		if basePath != "" {
			p = strings.TrimPrefix(p, basePath)
			if p == "" {
				continue
			}
			p = p[1:]
		}

		parts := strings.Split(p, "/")
		current := root
		for i, part := range parts {
			found := false
			for _, child := range current.Children {
				if child.Name == part {
					current = child
					found = true
					break
				}
			}
			if !found {
				newNode := &common.DataNode{
					Name:  part,
					Path:  strings.Join(parts[:i+1], "/"),
					IsDir: i < len(parts)-1,
				}
				current.Children = append(current.Children, newNode)
				current = newNode
			}
		}
	}

	return root.Children
}
