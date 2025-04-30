package biz

import (
	"fmt"
	"net/http"
	"net/url"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/pkg/common"
	imkdc "github.com/muidea/magicCodeAgent/internal/modules/kernel/debug/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Debug struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Debug {
	ptr := &Debug{
		Base: biz.New(imkdc.DebugModule, eventHub, backgroundRoutine),
	}

	return ptr
}

func (s *Debug) Run() (err *cd.Error) {
	filter := &common.MessageFilter{
		PageIndex: 1,
		PageSize:  100,
	}

	for {
		messageList := s.filterMessage(filter)
		for _, message := range messageList {
			s.traceContent(message)
		}
		if len(messageList) < filter.PageSize {
			break
		}
		filter.PageIndex++
	}

	return
}

func (s *Debug) traceContent(message *models.MessageInfo) {
	// TODO 根据todo补充逻辑实现
}

func (s *Debug) filterMessage(filter *common.MessageFilter) []*models.MessageInfo {
	result := &imbtc.MessageInfoListResult{}

	vals := url.Values{}
	vals = filter.Encode(vals)
	serverUrl := fmt.Sprintf("http://127.0.0.1:%s%s%s", config.BindPort(), common.ApiVersion, imbtc.FilterMessage)
	urlVal, urlErr := url.ParseRequestURI(serverUrl)
	if urlErr != nil {
		log.Errorf("filterMessage url parse failed, url:%s, err:%s", serverUrl, urlErr)
		return nil
	}
	urlVal.RawQuery = vals.Encode()
	_, httpErr := net.HTTPGet(http.DefaultClient, urlVal.String(), result)
	if httpErr != nil {
		return nil
	}

	return result.Values
}
