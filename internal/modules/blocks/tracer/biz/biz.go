package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/config"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/tracer/pkg/common"
)

type Tracer struct {
	biz.Base

	messageInfoCache *MessageSequenceCache
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Tracer {
	ptr := &Tracer{
		Base:             biz.New(common.TracerModule, eventHub, backgroundRoutine),
		messageInfoCache: NewMessageSequenceCache(config.MaxCacheItemSize()),
	}

	ptr.SubscribeFunc(common.AppendMessage, ptr.appendMessage)
	ptr.SubscribeFunc(common.FilterMessage, ptr.filterMessage)

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

	messagePtr, messageOK := msgVal.(*common.MessageInfo)
	if !messageOK {
		return
	}
	s.AppendMessage(messagePtr)
	if re != nil {
		re.Set(nil, nil)
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

func (s *Tracer) AppendMessage(msg *common.MessageInfo) {
	s.messageInfoCache.Add(msg)
}

func (s *Tracer) FilterMessage(filter *common.MessageFilter) (ret []*common.MessageInfo, err *cd.Error) {
	ret, _ = s.messageInfoCache.FilterMessage(filter)
	return
}
