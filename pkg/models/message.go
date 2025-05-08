package models

import "fmt"

type MessageItem struct {
	Content []string `json:"content"`
}

func (s *MessageItem) String() string {
	return fmt.Sprintf("Content:%v", s.Content)
}

type MessageItemList []MessageItem

type MessageInfo struct {
	SN         int64       `json:"sn"`
	Role       string      `json:"role"`
	Item       MessageItem `json:"item"`
	CreateTime string      `json:"createTime"`
	Workspace  string      `json:"workspace"`
}

func (s *MessageInfo) String() string {
	return fmt.Sprintf("SN:%d,Role:%s,Workspace:%s", s.SN, s.Role, s.Workspace)
}
