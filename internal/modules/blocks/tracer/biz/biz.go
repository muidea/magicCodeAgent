package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/config"
	imbtc "github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/cache"
	"github.com/muidea/magicCodeAgent/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"
)

type Tracer struct {
	biz.Base

	messageInfoCache *cache.MessageSequenceCache
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Tracer {
	ptr := &Tracer{
		Base:             biz.New(imbtc.TracerModule, eventHub, backgroundRoutine),
		messageInfoCache: cache.NewMessageSequenceCache(config.MaxCacheItemSize()),
	}

	ptr.SubscribeFunc(imbtc.AppendMessage, ptr.appendMessage)
	ptr.SubscribeFunc(imbtc.FilterMessage, ptr.filterMessage)

	return ptr
}

func (s *Tracer) Run() (err *cd.Error) {
	return
}

func (s *Tracer) appendMessage(ev event.Event, re event.Result) {
	msgVal := ev.Data()
	if msgVal == nil {
		return
	}

	msgPtr, msgOK := msgVal.(*models.MessageInfo)
	if !msgOK {
		return
	}
	msgPtr = s.AppendMessage(msgPtr)
	if re != nil {
		re.Set(msgPtr, nil)
	}
}

func (s *Tracer) filterMessage(ev event.Event, re event.Result) {
	filterVal := ev.Data()
	if filterVal == nil {
		return
	}
	filterPtr, filterOK := filterVal.(*common.MessageFilter)
	if !filterOK {
		return
	}
	ret, err := s.FilterMessage(filterPtr)
	if re != nil {
		re.Set(ret, err)
	}
}

func (s *Tracer) AppendMessage(msg *models.MessageInfo) *models.MessageInfo {
	return s.messageInfoCache.Add(msg)
}

func (s *Tracer) FilterMessage(filter *common.MessageFilter) (ret []*models.MessageInfo, err *cd.Error) {
	ret, _ = s.messageInfoCache.FilterMessage(filter)
	return
}
