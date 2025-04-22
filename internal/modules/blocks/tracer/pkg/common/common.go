package common

import (
	"net/http"
	"strconv"

	cd "github.com/muidea/magicCommon/def"
)

const TracerModule = "/internal/modules/blocks/tracer"

const (
	FilterMessage = "/message/filter"
	AppendMessage = "/message/append"
)

type MessageInfo struct {
	SN            int64  `json:"sn"`
	Type          string `json:"type"`
	Content       string `json:"content"`
	CreateTime    string `json:"createTime"`
	ModelProvider string `json:"modelProvider"`
	Workspace     string `json:"workspace"`
}

type MessageFilter struct {
	Workspace string `json:"workspace"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

func (s *MessageFilter) Decode(req *http.Request) {
	query := req.URL.Query()

	if workspace := query.Get("workspace"); workspace != "" {
		s.Workspace = workspace
	}

	// 解析pageIndex
	if pageIndex := query.Get("pageIndex"); pageIndex != "" {
		if idx, err := strconv.Atoi(pageIndex); err == nil && idx > 0 {
			s.PageIndex = idx
		}
	} else {
		s.PageIndex = 1 // 默认第一页
	}

	// 解析pageSize
	if pageSize := query.Get("pageSize"); pageSize != "" {
		if size, err := strconv.Atoi(pageSize); err == nil && size > 0 {
			if size > 100 {
				size = 100 // 限制最大100条/页
			}
			s.PageSize = size
		}
	} else {
		s.PageSize = 20 // 默认每页20条
	}
}

type MessageInfoList struct {
	cd.Result
	Values []*MessageInfo `json:"values"`
}

type Tracer interface {
	AppendMessage(msg *MessageInfo)
}
