package biz

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/task"

	icfg "github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/cli/config"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/cli/pkg/common"
	tc "github.com/muidea/magicCodeAgent/pkg/common"

	ic "github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator/pkg/common"
	mc "github.com/muidea/magicCodeAgent/internal/modules/blocks/module/pkg/common"
	pc "github.com/muidea/magicCodeAgent/internal/modules/blocks/project/pkg/common"
)

type Cli struct {
	biz.Base
}

func New(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine,
) *Cli {
	ret := &Cli{
		Base: biz.New(common.CliModule, eventHub, backgroundRoutine),
	}
	return ret
}

func (s *Cli) Run() {
	err := config.VerifyVar()
	if err != nil {
		log.Errorf("verify config failed,err:%s", err.Error())
		return
	}

	actionType := config.GetActionType()
	switch actionType {
	case config.ModuleVal:
		s.runModuleAction()
	case config.InitiatorVal:
		s.runInitiatorAction()
	case config.ProjectVal:
		s.runProjectAction()
	default:
		log.Warnf("unknown module type:%s", actionType)
	}
}

func (s *Cli) runModuleAction() {
	param := &tc.CreateParam{
		EnableDao:     config.GetEnableDao(),
		EnableService: config.GetEnableService(),
		ModuleName:    config.GetModuleName(),
		WorkingPath:   icfg.CurrentWorkPath(),
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

func (s *Cli) runInitiatorAction() {
	param := &tc.CreateParam{
		EnableDao:     config.GetEnableDao(),
		EnableService: config.GetEnableService(),
		ModuleName:    config.GetModuleName(),
		WorkingPath:   icfg.CurrentWorkPath(),
	}
	var ev event.Event
	switch config.GetExecuteAction() {
	case config.CreateVal:
		ev = event.NewEvent(ic.CreateNewInitiator, s.ID(), ic.InitiatorModule, nil, param)
	case config.DestroyVal:
		ev = event.NewEvent(ic.DestroyInitiator, s.ID(), ic.InitiatorModule, nil, param)
	case config.RefreshVal:
		ev = event.NewEvent(ic.RefreshInitiator, s.ID(), ic.InitiatorModule, nil, param)
	default:
		log.Warnf("unknown execute action:%s", config.GetExecuteAction())
	}

	s.SendEvent(ev)
}

func (s *Cli) runProjectAction() {
	param := &tc.CreateParam{
		ProjectName:        config.GetProjectName(),
		ProjectDescription: config.GetProjectDescription(),
		EnableDao:          config.GetEnableDao(),
		EnableService:      config.GetEnableService(),
		WorkingPath:        icfg.CurrentWorkPath(),
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
