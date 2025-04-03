package biz

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
	"github.com/muidea/magicToolkit/internal/modules/base/biz"
	"github.com/muidea/magicToolkit/internal/modules/module/pkg/common"
)

type Module struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Module {
	ptr := &Module{
		Base: biz.New(common.ModuleModule, eventHub, backgroundRoutine),
	}

	ptr.SubscribeFunc(common.CreateNewModule, ptr.create)
	ptr.SubscribeFunc(common.RefreshModule, ptr.refresh)
	ptr.SubscribeFunc(common.DestroyModule, ptr.destroy)

	return ptr
}

func (s *Module) create(ev event.Event, re event.Result) {

}

func (s *Module) refresh(ev event.Event, re event.Result) {

}

func (s *Module) destroy(ev event.Event, re event.Result) {

}
