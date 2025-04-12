package biz

import (
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/project/pkg/common"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
)

type Project struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Project {
	ret := &Project{
		Base: biz.New(common.ProjectModule, eventHub, backgroundRoutine),
	}

	ret.SubscribeFunc(common.CreateNewProject, ret.create)
	ret.SubscribeFunc(common.RefreshProject, ret.refresh)
	ret.SubscribeFunc(common.DestroyProject, ret.destroy)

	return ret
}

func (s *Project) create(ev event.Event, re event.Result) {

}

func (s *Project) refresh(ev event.Event, re event.Result) {

}

func (s *Project) destroy(ev event.Event, re event.Result) {

}
