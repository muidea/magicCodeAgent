package models

import "fmt"

type MessageItem struct {
	Role    string   `json:"role"`
	Content []string `json:"content,omitempty"`
}

func (s *MessageItem) String() string {
	return fmt.Sprintf("Role:%s, Content item len:%d", s.Role, len(s.Content))
}

type MessageItemList []MessageItem

type MessageInfo struct {
	SN         int64       `json:"sn"`
	Item       MessageItem `json:"item"`
	CreateTime string      `json:"createTime"`
	Workspace  string      `json:"workspace"`
}

func (s *MessageInfo) String() string {
	return fmt.Sprintf("SN:%d,Role:%s,Workspace:%s", s.SN, s.Item.Role, s.Workspace)
}
