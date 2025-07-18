package biz

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/cli/config"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/cli/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"

	ic "github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator/pkg/common"
	mc "github.com/muidea/magicCodeAgent/internal/modules/blocks/module/pkg/common"
	pc "github.com/muidea/magicCodeAgent/internal/modules/blocks/project/pkg/common"
	imbtmc "github.com/muidea/magicCodeAgent/internal/modules/blocks/templatemanager/pkg/common"
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

func (s *Cli) queryTemplate(templateName string) (ret *models.Template) {
	loadEvent := event.NewEvent(imbtmc.LoadTemplate, s.ID(), imbtmc.TemplatemanagerModule, nil, templateName)
	result := s.SendEvent(loadEvent)
	templateVal, templateErr := result.Get()
	if templateErr != nil {
		log.Errorf("queryTemplate failed, error:%s", templateErr.Error())
		return
	}
	ret = templateVal.(*models.Template)
	return
}

func (s *Cli) Run() {
	err := config.VerifyVar()
	if err != nil {
		log.Errorf("verify config failed, error:%s", err.Error())
		return
	}

	actionType := config.GetActionType()
	workspacePtr := config.GetCurrentWorkspace()
	templatePtr := s.queryTemplate(workspacePtr.Template)
	if templatePtr == nil {
		log.Errorf("query template failed, templateName:%s", workspacePtr.Template)
		return
	}

	switch actionType {
	case config.ModuleKey:
		s.runModuleAction(workspacePtr, templatePtr)
	case config.InitiatorKey:
		s.runInitiatorAction(workspacePtr, templatePtr)
	case config.ProjectKey:
		s.runProjectAction(workspacePtr, templatePtr)
	default:
		log.Warnf("unknown module type:%s", actionType)
	}
}

func (s *Cli) runModuleAction(workspacePtr *models.Workspace, templatePtr *models.Template) {
	var ev event.Event
	switch config.GetExecuteAction() {
	case config.CreateKey:
		ev = event.NewEvent(mc.CreateNewModule, s.ID(), mc.ModuleModule, nil, workspacePtr)
	case config.DestroyKey:
		ev = event.NewEvent(mc.DestroyModule, s.ID(), mc.ModuleModule, nil, workspacePtr)
	case config.RefreshKey:
		ev = event.NewEvent(mc.RefreshModule, s.ID(), mc.ModuleModule, nil, workspacePtr)
	default:
		log.Warnf("unknown execute action:%s", config.GetExecuteAction())
	}

	ev.SetData("template", templatePtr)
	s.SendEvent(ev)
}

func (s *Cli) runInitiatorAction(workspacePtr *models.Workspace, templatePtr *models.Template) {
	var ev event.Event
	switch config.GetExecuteAction() {
	case config.CreateKey:
		ev = event.NewEvent(ic.CreateNewInitiator, s.ID(), ic.InitiatorModule, nil, workspacePtr)
	case config.DestroyKey:
		ev = event.NewEvent(ic.DestroyInitiator, s.ID(), ic.InitiatorModule, nil, workspacePtr)
	case config.RefreshKey:
		ev = event.NewEvent(ic.RefreshInitiator, s.ID(), ic.InitiatorModule, nil, workspacePtr)
	default:
		log.Warnf("unknown execute action:%s", config.GetExecuteAction())
	}

	ev.SetData("template", templatePtr)
	s.SendEvent(ev)
}

func (s *Cli) runProjectAction(workspacePtr *models.Workspace, templatePtr *models.Template) {
	var ev event.Event
	switch config.GetExecuteAction() {
	case config.CreateKey:
		ev = event.NewEvent(pc.CreateNewProject, s.ID(), pc.ProjectModule, nil, workspacePtr)
	case config.DestroyKey:
		ev = event.NewEvent(pc.DestroyProject, s.ID(), pc.ProjectModule, nil, workspacePtr)
	case config.RefreshKey:
		ev = event.NewEvent(pc.RefreshProject, s.ID(), pc.ProjectModule, nil, workspacePtr)
	default:
		log.Warnf("unknown execute action:%s", config.GetExecuteAction())
	}

	ev.SetData("template", templatePtr)
	s.SendEvent(ev)
}
