package models

import "fmt"

type MessageInfo struct {
	SN            int64  `json:"sn"`
	Type          string `json:"type"`
	Content       string `json:"content"`
	CreateTime    string `json:"createTime"`
	ModelProvider string `json:"modelProvider"`
	Workspace     string `json:"workspace"`
}

func (s *MessageInfo) String() string {
	return fmt.Sprintf("MessageInfo:{SN:%d,Type:%s,Content:%s,CreateTime:%s,ModelProvider:%s,Workspace:%s}", s.SN, s.Type, s.Content, s.CreateTime, s.ModelProvider, s.Workspace)
}
