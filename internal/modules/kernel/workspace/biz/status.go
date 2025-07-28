package biz

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/muidea/magicCodeAgent/pkg/models"
)

type FlowStatus struct {
	Name      string        `json:"name"`
	Level     string        `json:"level"`
	Status    string        `json:"status"`
	FlowPtr   *models.Flow  `json:"flowPtr"`
	Childrens []*FlowStatus `json:"childrens"`
}

func (s *FlowStatus) String() string {
	switch s.Level {
	case models.ProjectLevel:
		return fmt.Sprintf("%s, 层级：项目（`%s`), 状态：%s", s.Name, s.Level, s.Status)
	case models.ModuleLevel:
		return fmt.Sprintf("%s, 层级：业务模块(`%s`), 状态：%s", s.Name, s.Level, s.Status)
	case models.BlockLevel:
		return fmt.Sprintf("%s, 层级：功能模块(`%s`), 状态：%s", s.Name, s.Level, s.Status)
	default:
		return ""
	}
}

func (s *FlowStatus) StartFlow(ptr *models.Flow) {
	flowStack := []*models.Flow{}
	for ptr != nil {
		flowStack = append(flowStack, ptr)
		ptr = ptr.GetParent()
		if ptr.Level == models.RootLevel {
			break
		}
	}
	idx := len(flowStack) - 1
	curPtr := s
	for {
		flow := flowStack[idx]
		if flow.GetLevel() != curPtr.Level || flow.GetName() != curPtr.Name {
			break
		}

		idx--
		if idx < 0 {
			curPtr.Status = "正在开发"
			curPtr.FlowPtr = ptr
			break
		}
		flow = flowStack[idx]
		for _, val := range curPtr.Childrens {
			if flow.GetLevel() == val.Level && flow.GetName() == val.Name {
				curPtr = val
				break
			}
		}
	}
}

func (s *FlowStatus) FinishFlow(ptr *models.Flow) {
	flowStack := []*models.Flow{}
	for ptr != nil {
		flowStack = append(flowStack, ptr)
		ptr = ptr.GetParent()
		if ptr.Level == models.RootLevel {
			break
		}
	}
	idx := len(flowStack) - 1
	curPtr := s
	for {
		flow := flowStack[idx]
		if flow.GetLevel() != curPtr.Level || flow.GetName() != curPtr.Name {
			break
		}

		idx--
		if idx < 0 {
			curPtr.Status = "已完成"
			break
		}
		flow = flowStack[idx]
		for _, val := range curPtr.Childrens {
			if flow.GetLevel() == val.Level && flow.GetName() == val.Name {
				curPtr = val
				break
			}
		}
	}
}

func NewFlowStatus(flowList []string) *FlowStatus {
	root := &FlowStatus{
		Name:      "root",
		Level:     "root",
		Status:    "未开发",
		Childrens: []*FlowStatus{},
	}

	for _, path := range flowList {
		trimmed := strings.TrimSuffix(path, "-项目功能.md")
		trimmed = strings.TrimSuffix(trimmed, "-业务模块.md")
		trimmed = strings.TrimSuffix(trimmed, "-功能模块.md")
		parts := strings.Split(trimmed, "/")

		insertFlowStatus(root, parts, 0)
	}

	return root.Childrens[0]
}

func insertFlowStatus(current *FlowStatus, parts []string, depth int) {
	if depth >= len(parts) {
		return
	}

	name := parts[depth]
	level := ""
	switch depth {
	case 0:
		level = "project"
	case 1:
		level = "module"
	default:
		level = "block"
	}

	// 查找是否已存在该子节点
	var child *FlowStatus
	for _, c := range current.Childrens {
		if c.Name == name {
			child = c
			break
		}
	}

	// 不存在就创建新节点
	if child == nil {
		child = &FlowStatus{
			Name:      name,
			Level:     level,
			Status:    "未开发",
			Childrens: []*FlowStatus{},
		}
		current.Childrens = append(current.Childrens, child)
	}

	// 递归处理下一层
	insertFlowStatus(child, parts, depth+1)
}

// 打印 FlowStatus 树结构
func DumpFlowStatusTree(root *FlowStatus, buf *bytes.Buffer) {
	// 打印 root 本身，不加前缀
	buf.WriteString(fmt.Sprintf("%s\n", root))

	// 打印子节点
	for i, child := range root.Childrens {
		isLast := i == len(root.Childrens)-1
		dumpNode(child, buf, "", isLast)
	}

	buf.WriteString("\n")
}

func dumpNode(node *FlowStatus, buf *bytes.Buffer, prefix string, isLast bool) {
	connector := "├── "
	childPrefix := prefix + "│   "
	if isLast {
		connector = "└── "
		childPrefix = prefix + "    "
	}
	buf.WriteString(fmt.Sprintf("%s%s%s\n", prefix, connector, node))

	for i, child := range node.Childrens {
		isLastChild := i == len(node.Childrens)-1
		dumpNode(child, buf, childPrefix, isLastChild)
	}
}
