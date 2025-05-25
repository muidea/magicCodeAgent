package models

import (
	"fmt"

	fu "github.com/muidea/magicCommon/foundation/util"
)

const (
	FLOW_STATUS_WAIT = iota
	FLOW_STATUS_FINISH
)

type Flow struct {
	Id        string           `json:"id"`
	Title     string           `json:"title"`
	Message   []*MessageItem   `json:"message"`
	Status    int              `json:"status"`
	AddFiles  fu.StringSortSet `json:"addFiles"`
	ModFiles  fu.StringSortSet `json:"modFiles"`
	DelFiles  fu.StringSortSet `json:"delFiles"`
	Parent    *Flow            `json:"parent"`
	Childrens []*Flow          `json:"childrens"`
	Level     string           `json:"level"`
}

func (s *Flow) String() string {
	return fmt.Sprintf("Flow:{%s}, level:%s, status:%d", s.Title, s.Level, s.Status)
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

func (s *Flow) GetFlowPath() []*Flow {
	nodes := make([]*Flow, 0)
	if s.Parent != nil {
		nodes = append(nodes, s.Parent.GetFlowPath()...)
	}
	nodes = append(nodes, s)

	// 反转slice使父节点在前
	flowPath := make([]*Flow, len(nodes))
	for i, j := 0, len(nodes)-1; j >= 0; i, j = i+1, j-1 {
		flowPath[i] = nodes[j]
	}

	return flowPath
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

func (s *Flow) SetLevel(level string) {
	s.Level = level
}
