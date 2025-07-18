package biz

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	fp "github.com/muidea/magicCommon/foundation/path"

	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/workspace/pkg/common"
)

// 判断是否忽略
// .git .tmp 隐藏文件都忽略，不用处理
func isIgnored(path string) bool {
	base := filepath.Base(path)
	return strings.Contains(path, ".git") ||
		strings.Contains(path, ".tmp") ||
		strings.Contains(path, ".umi") ||
		strings.Contains(path, ".umi-production") ||
		strings.Contains(path, "node_modules") ||
		strings.Contains(path, "vendor") ||
		strings.EqualFold(base, ".") ||
		strings.EqualFold(base, "..")
}

func NewBasePathFilter(rootPath string, basePath string) *basePathFilter {
	return &basePathFilter{
		basePath: filepath.Clean(path.Join(rootPath, basePath)),
	}
}

type basePathFilter struct {
	basePath string
}

func (s *basePathFilter) Filter(val interface{}) bool {
	if val == nil {
		return false
	}

	pathVal, pathOK := val.(string)
	if !pathOK {
		return false
	}

	return fp.IsSubPath(s.basePath, pathVal)
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

func formatUsage(ptr *imklc.Usage, loopCount int64) string {
	strBuff := &strings.Builder{}
	strBuff.WriteString("对话次数: " + fmt.Sprintf("%d", loopCount) + "\n")
	strBuff.WriteString("Token用量: \n")
	strBuff.WriteString("  - 总消耗: " + fmt.Sprintf("%d", ptr.TotalTokens) + "\n")
	strBuff.WriteString("  - 输入: " + fmt.Sprintf("%d", ptr.PromptTokens) + "(Hit-" + fmt.Sprintf("%d", ptr.PromptCacheHitToken) + "/Miss-" + fmt.Sprintf("%d", ptr.PromptMissHitToken) + ")\n")
	strBuff.WriteString("  - 输出: " + fmt.Sprintf("%d", ptr.CompletionTokens) + "\n")
	return strBuff.String()
}
