package common

import (
	cd "github.com/muidea/magicCommon/def"
)

const TracerModule = "/internal/modules/blocks/tracer"

const (
	FilterMessage = "/message/filter"
	AppendMessage = "/message/append"
)

type MessageInfo struct {
	SN         int64  `json:"sn"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	Owner      string `json:"owner"`
}

type MessageFilter struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type MessageInfoList struct {
	cd.Result
	Values []*MessageInfo `json:"values"`
}

type Tracer interface {
	AppendMessage(msg *MessageInfo)
}
