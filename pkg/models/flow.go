package models

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"
)

const (
	StatusInit = iota
	StatusDoing
	StatusFinish
	StatusDone
)

const (
	RootLevel    = ""
	ProjectLevel = "project"
	ModuleLevel  = "module"
	BlockLevel   = "block"
)

const (
	LevelGreater = 1
	LevelEqual   = 0
	LevelLess    = -1
)

var designStatus = map[int]string{
	0: "未设计",
	1: "正在设计",
	2: "设计文档已输出",
	3: "设计完成",
}

var codeStatus = map[int]string{
	0: "未开发",
	1: "正在开发",
	2: "代码已输出",
	3: "开发完成",
}

// 比较l1, l3的层级关系
// 如果l1 > l2, 返回1，表示l1比l2高一个层级
// 如果l1 < l2, 返回-1，表示l1比l2低一个层级
// 如果l1 == l2, 返回0，表示l1和l2在同一层级
// Project > Module > Block > Point
func CompareLevel(l1, l2 string) int {
	levelMap := map[string]int{
		RootLevel:    5,
		ProjectLevel: 4,
		ModuleLevel:  3,
		BlockLevel:   2,
	}

	l1Level := levelMap[l1]
	l2Level := levelMap[l2]

	switch {
	case l1Level > l2Level:
		return 1
	case l1Level < l2Level:
		return -1
	default:
		return 0
	}
}

func GetChildLevel(level string) string {
	switch level {
	case RootLevel:
		return ProjectLevel
	case ProjectLevel:
		return ModuleLevel
	case ModuleLevel:
		return BlockLevel
	default:
		return ""
	}
}

type Flow struct {
	Id            string           `json:"id"`
	Name          string           `json:"name"`
	Level         string           `json:"level"`
	Stage         string           `json:"stage"`
	Status        int              `json:"status"`
	Path          string           `json:"path"`
	Message       []*MessageItem   `json:"-"`
	MessageResume bool             `json:"-"`
	AddFiles      fu.StringSortSet `json:"addFiles,omitempty"`
	ModFiles      fu.StringSortSet `json:"modFiles,omitempty"`
	DelFiles      fu.StringSortSet `json:"delFiles,omitempty"`
	TodoList      TodoList         `json:"todoList"`
	Parent        *Flow            `json:"-"`
	Childrens     []*Flow          `json:"childrens"`
}

func (s *Flow) GetName() string {
	return s.Name
}

func (s *Flow) GetLevel() string {
	return s.Level
}

func (s *Flow) GetStage() string {
	return s.Stage
}

func (s *Flow) GetPath() string {
	return s.Path
}

func (s *Flow) GetStatus() string {
	getStatusStr := func() string {
		switch s.Stage {
		case DesignStage:
			return designStatus[s.Status]
		default:
			return codeStatus[s.Status]
		}
	}
	switch s.Level {
	case ProjectLevel:
		return fmt.Sprintf("%s, 层级：项目（`%s`), 状态：%s", s.Name, s.Level, getStatusStr())
	case ModuleLevel:
		return fmt.Sprintf("%s, 层级：业务模块(`%s`), 状态：%s", s.Name, s.Level, getStatusStr())
	case BlockLevel:
		return fmt.Sprintf("%s, 层级：功能模块(`%s`), 状态：%s", s.Name, s.Level, getStatusStr())
	default:
		return ""
	}
}

func (s *Flow) String() string {
	if s.Parent == nil {
		var stringBuff strings.Builder
		stringBuff.WriteString(fmt.Sprintf("id:%s, name:%s, level:%s, path:%s, stage:%s, status:%d, parent is null", s.Id, s.Name, s.Level, s.Path, s.Stage, s.Status))
		return stringBuff.String()
	}

	var stringBuff strings.Builder
	stringBuff.WriteString(fmt.Sprintf("id:%s, name:%s, level:%s, path:%s, stage:%s, status:%d, parent is %s", s.Id, s.Name, s.Level, s.Path, s.Stage, s.Status, s.Parent.Id))
	return stringBuff.String()
}

func (s *Flow) AppendMessage(m *MessageItem) {
	if s.Message == nil {
		s.Message = make([]*MessageItem, 0)
	}
	s.Message = append(s.Message, m)
	s.MessageResume = false
}

func (s *Flow) ResumeFlow(m *MessageItem) {
	if !s.MessageResume {
		s.Message = append(s.Message, m)
		s.MessageResume = true
		return
	}
	log.Infof("resume flow:%s, level:%s, path:%s, status:%d, parent:%s", s.GetName(), s.GetLevel(), s.GetPath(), s.GetStatus(), s.Parent.GetName())
}

func (s *Flow) GetChildrens() []*Flow {
	return s.Childrens
}

func (s *Flow) GetParent() *Flow {
	return s.Parent
}

func (s *Flow) IsFinish() bool {
	return s.Status == StatusDone
}

func (s *Flow) IsLeaf() bool {
	return s.Level == BlockLevel
}

func (s *Flow) CanFinish() (ret bool) {
	if s.IsFinish() {
		ret = true
		return
	}

	// 如果Flow没有子节点，并且是最终level，并且已全部完成，直接认为是可以结束的
	if len(s.Childrens) == 0 {
		ret = s.IsLeaf() && s.Status >= StatusFinish
		return
	}

	// 这里要检查其子节点是否可以结束
	// 任意一个子节点不能结束，则认为该Flow不能结束
	for _, val := range s.Childrens {
		if !val.CanFinish() {
			ret = false
			return
		}
	}

	// 所有的子节点都能结束，则认为该Flow可以结束
	if s.Status >= StatusFinish {
		ret = true
	}
	return
}

func (s *Flow) GetFlowPath() []*Flow {
	nodes := []*Flow{}
	nodes = append(nodes, s)
	pPtr := s.Parent
	for {
		if pPtr == nil {
			break
		}
		nodes = append(nodes, pPtr)
		pPtr = pPtr.Parent
	}

	// 反转slice使父节点在前
	flowPath := make([]*Flow, len(nodes))

	for i, j := 0, len(nodes)-1; j >= 0; i, j = i+1, j-1 {
		flowPath[i] = nodes[j]
	}

	return flowPath
	//return []*Flow{s}
}

func (s *Flow) AddFile(filePath string) {
	s.AddFiles = s.AddFiles.Add(filePath)
	s.DelFiles = s.DelFiles.Remove(filePath)
}

func (s *Flow) ModFile(filePath string) {
	s.ModFiles = s.ModFiles.Add(filePath)
}

func (s *Flow) DelFile(filePath string) {
	s.DelFiles = s.DelFiles.Add(filePath)
	s.AddFiles = s.AddFiles.Remove(filePath)
}

func (s *Flow) GetTodoList() TodoList {
	return s.TodoList
}

func (s *Flow) UpdateTodoList(todoList TodoList) {
	s.TodoList = todoList
}

func (s *Flow) GetLatestUnFinishFlow() *Flow {
	if s.CanFinish() {
		return nil
	}

	for _, child := range s.Childrens {
		if child.CanFinish() {
			continue
		}

		if unfinished := child.GetLatestUnFinishFlow(); unfinished != nil {
			return unfinished
		}
	}

	return s
}

// 打印 Flow 树结构
func DumpFlowStatusTree(root *Flow, buf *bytes.Buffer) {
	// 打印 root 本身，不加前缀
	buf.WriteString(fmt.Sprintf("%s\n", root.GetStatus()))

	// 打印子节点
	for i, child := range root.Childrens {
		isLast := i == len(root.Childrens)-1
		dumpFlowStatus(child, buf, "", isLast)
	}

	buf.WriteString("\n")
}

func dumpFlowStatus(node *Flow, buf *bytes.Buffer, prefix string, isLast bool) {
	connector := "├── "
	childPrefix := prefix + "│   "
	if isLast {
		connector = "└── "
		childPrefix = prefix + "    "
	}
	buf.WriteString(fmt.Sprintf("%s%s%s\n", prefix, connector, node.GetStatus()))

	for i, child := range node.Childrens {
		isLastChild := i == len(node.Childrens)-1
		dumpFlowStatus(child, buf, childPrefix, isLastChild)
	}
}
