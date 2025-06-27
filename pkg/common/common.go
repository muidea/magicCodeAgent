package common

import (
	"net/http"
	"net/url"
	"strconv"

	imklc "github.com/muidea/magicCodeAgent/internal/modules/kernel/llm/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type FlowID struct{}

func (s *FlowID) String() string {
	return "FlowID"
}

type ContentSink struct{}

func (s *ContentSink) String() string {
	return "ContentSink"
}

type ContentSinker interface {
	OnClose()
	OnRecv(event, content string, usagePtr *imklc.Usage)
}

type CreateParam struct {
	ProjectName        string `json:"projectName omitempty"`
	ProjectType        string `json:"projectType"`
	ProjectDescription string `json:"projectDescription"`
	ModuleName         string `json:"moduleName"`
	WorkingPath        string `json:"workingPath"`
	EnableService      bool   `json:"enableService"`
	EnableDao          bool   `json:"enableDao"`
}

type MessageSink interface {
	OnTerminate(msgPtr *models.MessageInfo)
	OnSendMessage(msgPtr *models.MessageInfo)
	OnRecvMessage(msgPtr *models.MessageInfo, usagePtr *imklc.Usage)
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

func (s *MessageFilter) Encode(vals url.Values) url.Values {
	if s.Workspace != "" {
		vals.Add("workspace", s.Workspace)
	}
	if s.PageIndex > 0 {
		vals.Add("pageIndex", strconv.Itoa(s.PageIndex))
	}
	if s.PageSize > 0 {
		vals.Add("pageSize", strconv.Itoa(s.PageSize))
	}
	return vals
}
