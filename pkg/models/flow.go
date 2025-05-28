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
