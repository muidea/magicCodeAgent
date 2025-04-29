package models

import "fmt"

type MessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (s *MessageItem) String() string {
	return fmt.Sprintf("Role:%s,Content:%s", s.Role, s.Content)
}

type MessageInfo struct {
	SN            int64        `json:"sn"`
	Item          *MessageItem `json:"item"`
	CreateTime    string       `json:"createTime"`
	ModelProvider string       `json:"modelProvider"`
	Workspace     string       `json:"workspace"`
}

func (s *MessageInfo) String() string {
	return fmt.Sprintf("MessageInfo:{SN:%d,Content:[%v],CreateTime:%s,ModelProvider:%s,Workspace:%s}", s.SN, s.Item, s.CreateTime, s.ModelProvider, s.Workspace)
}
