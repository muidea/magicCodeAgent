package biz

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/modules/base/biz"
	"github.com/muidea/magicToolkit/internal/modules/module/pkg/common"
	"github.com/muidea/magicToolkit/pkg/util"
)

type Initator struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Initator {
	ptr := &Initator{
		Base: biz.New(common.ModuleModule, eventHub, backgroundRoutine),
	}

	ptr.SubscribeFunc(common.CreateNewModule, ptr.create)
	ptr.SubscribeFunc(common.RefreshModule, ptr.refresh)
	ptr.SubscribeFunc(common.DestroyModule, ptr.destroy)

	return ptr
}

func (s *Initator) create(ev event.Event, re event.Result) {
	paramPtr, paramOK := ev.Data().(*common.CreateModuleParam)
	if !paramOK {
		if re != nil {
			re.Set(nil, cd.NewError(cd.Unexpected, "illegal create module param"))
		}
		return
	}

	err := util.VerifyModuleWorkPath(paramPtr.WorkPath)
	if err != nil {
		if re != nil {
			re.Set(nil, err)
		}
		return
	}
}

func (s *Initator) refresh(ev event.Event, re event.Result) {

}

func (s *Initator) destroy(ev event.Event, re event.Result) {

}
