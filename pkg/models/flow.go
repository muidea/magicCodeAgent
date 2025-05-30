package models

import (
	"fmt"
	"strings"

	fu "github.com/muidea/magicCommon/foundation/util"
)

const (
	FLOW_STATUS_WAIT = iota
	FLOW_STATUS_FINISH
)

const (
	RootLevel    = "root"
	ProjectLevel = "project"
	ModuleLevel  = "module"
	BlockLevel   = "block"
)

const (
	LevelGreater = 1
	LevelEqual   = 0
	LevelLess    = -1
)

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
	Id        string           `json:"id"`
	Title     string           `json:"title"`
	Message   []*MessageItem   `json:"message,omitempty"`
	Status    int              `json:"status"`
	AddFiles  fu.StringSortSet `json:"addFiles,omitempty"`
	ModFiles  fu.StringSortSet `json:"modFiles,omitempty"`
	DelFiles  fu.StringSortSet `json:"delFiles,omitempty"`
	Parent    *Flow            `json:"-"`
	Childrens []*Flow          `json:"childrens"`
	Level     string           `json:"level"`
}

func (s *Flow) String() string {
	if s.Parent == nil {
		var stringBuff strings.Builder
		stringBuff.WriteString(fmt.Sprintf("Flow %s, title:%s, level:%s, status:%d, parent is null", s.Id, s.Title, s.Level, s.Status))
		for _, val := range s.Childrens {
			stringBuff.WriteString(fmt.Sprintf("\n - children:%s, title:%s, level:%s, status:%d, parent:%s", val.Id, val.Title, val.Level, val.Status, val.Parent.Id))
		}

		return stringBuff.String()
	}

	var stringBuff strings.Builder
	stringBuff.WriteString(fmt.Sprintf("Flow %s, title:%s, level:%s, status:%d, parent is %s", s.Id, s.Title, s.Level, s.Status, s.Parent.Id))
	for _, val := range s.Childrens {
		stringBuff.WriteString(fmt.Sprintf("\n - children:%s, title:%s, level:%s, status:%d, parent:%s", val.Id, val.Title, val.Level, val.Status, val.Parent.Id))
	}
	return stringBuff.String()
}

func (s *Flow) AppendMessage(m *MessageItem) {
	if s.Message == nil {
		s.Message = make([]*MessageItem, 0)
	}
	s.Message = append(s.Message, m)
}

func (s *Flow) GetChildrens() []*Flow {
	return s.Childrens
}

func (s *Flow) GetParent() *Flow {
	return s.Parent
}

func (s *Flow) IsFinish() bool {
	return s.Status == FLOW_STATUS_FINISH
}

func (s *Flow) CanFinish() bool {
	// 如果Flow没有子节点，并且是最终level，直接认为是可以结束的
	if len(s.Childrens) == 0 {
		return s.Level == BlockLevel
	}

	// 这里要检查其子节点是否可以结束
	// 任意一个子节点不能结束，则认为该Flow不能结束
	for _, val := range s.Childrens {
		if !val.CanFinish() {
			return false
		}
	}

	// 所有的子节点都能结束，则认为该Flow可以结束
	return true
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

func (s *Flow) GetLevel() string {
	return s.Level
}
