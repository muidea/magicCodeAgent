package project

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicToolkit/internal/modules/project/biz"
	"github.com/muidea/magicToolkit/internal/modules/project/pkg/common"
)

func init() {
	module.Register(New())
}

func New() *Project {
	return &Project{}
}

type Project struct {
	bizPtr *biz.Project
}

func (s *Project) ID() string {
	return common.ProjectModule
}

func (s *Project) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {

	s.bizPtr = biz.New(eventHub, backgroundRoutine)
	return
}

func (s *Project) Run() (err *cd.Error) {
	return
}

func (s *Project) Teardown() {
}
