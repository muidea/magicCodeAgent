package biz

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/task"

	iconfig "github.com/muidea/magicToolkit/internal/config"
	"github.com/muidea/magicToolkit/internal/modules/base/biz"
	"github.com/muidea/magicToolkit/internal/modules/executor/config"
	"github.com/muidea/magicToolkit/internal/modules/executor/pkg/common"
	tc "github.com/muidea/magicToolkit/pkg/common"

	ic "github.com/muidea/magicToolkit/internal/modules/initator/pkg/common"
	mc "github.com/muidea/magicToolkit/internal/modules/module/pkg/common"
	pc "github.com/muidea/magicToolkit/internal/modules/project/pkg/common"
)

type Executor struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Executor {
	ret := &Executor{
		Base: biz.New(common.ExecutorModule, eventHub, backgroundRoutine),
	}
	return ret
}

func (s *Executor) Run() {
	err := config.VerifyVar()
	if err != nil {
		return
	}

	moduleType := config.GetModuleType()
	switch moduleType {
	case config.ModuleVal:
		s.runModuleAction()
	case config.InitatorVal:
		s.runInitatorAction()
	case config.ProjectVal:
		s.runProjectAction()
	default:
		log.Warnf("unknown module type:%s", moduleType)
	}
}

func (s *Executor) runModuleAction() {
	param := &tc.CreateParam{
		EnableDao:     config.GetEnableDao(),
		EnableService: config.GetEnableService(),
		Name:          config.GetModuleName(),
		WorkingPath:   iconfig.CurrentWorkPath(),
	}
	var ev event.Event
	switch config.GetExecuteAction() {
	case config.CreateVal:
		ev = event.NewEvent(mc.CreateNewModule, s.ID(), mc.ModuleModule, nil, param)
	case config.DestroyVal:
		ev = event.NewEvent(mc.DestroyModule, s.ID(), mc.ModuleModule, nil, param)
	case config.RefreshVal:
		ev = event.NewEvent(mc.RefreshModule, s.ID(), mc.ModuleModule, nil, param)
	default:
		log.Warnf("unknown execute action:%s", config.GetExecuteAction())
	}

	s.SendEvent(ev)
}

func (s *Executor) runInitatorAction() {
	param := &tc.CreateParam{
		EnableDao:     config.GetEnableDao(),
		EnableService: config.GetEnableService(),
		Name:          config.GetModuleName(),
		WorkingPath:   iconfig.CurrentWorkPath(),
	}
	var ev event.Event
	switch config.GetExecuteAction() {
	case config.CreateVal:
		ev = event.NewEvent(ic.CreateNewInitator, s.ID(), ic.InitatorModule, nil, param)
	case config.DestroyVal:
		ev = event.NewEvent(ic.DestroyInitator, s.ID(), ic.InitatorModule, nil, param)
	case config.RefreshVal:
		ev = event.NewEvent(ic.RefreshInitator, s.ID(), ic.InitatorModule, nil, param)
	default:
		log.Warnf("unknown execute action:%s", config.GetExecuteAction())
	}

	s.SendEvent(ev)
}

func (s *Executor) runProjectAction() {
	param := &tc.CreateParam{
		EnableDao:     config.GetEnableDao(),
		EnableService: config.GetEnableService(),
		Name:          config.GetModuleName(),
		WorkingPath:   iconfig.CurrentWorkPath(),
	}
	var ev event.Event
	switch config.GetExecuteAction() {
	case config.CreateVal:
		ev = event.NewEvent(pc.CreateNewProject, s.ID(), pc.ProjectModule, nil, param)
	case config.DestroyVal:
		ev = event.NewEvent(pc.DestroyProject, s.ID(), pc.ProjectModule, nil, param)
	case config.RefreshVal:
		ev = event.NewEvent(pc.RefreshProject, s.ID(), pc.ProjectModule, nil, param)
	default:
		log.Warnf("unknown execute action:%s", config.GetExecuteAction())
	}

	s.SendEvent(ev)
}
