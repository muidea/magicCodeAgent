package biz

import (
	"os"

	cd "github.com/muidea/magicCommon/def"
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
	loadEvent := event.NewEvent(imbtmc.QueryTemplate, s.ID(), imbtmc.TemplatemanagerModule, nil, templateName)
	result := s.SendEvent(loadEvent)
	templateVal, templateErr := result.Get()
	if templateErr != nil {
		log.Errorf("queryTemplate failed, name:%s, error:%s", templateName, templateErr.Error())
		return
	}
	ret = templateVal.(*models.Template)
	return
}

func (s *Cli) Run() (err *cd.Error) {
	defer func() {
		if err != nil {
			log.Errorf("run failed, error:%s", err.Error())
			os.Exit(1)
		}
	}()

	err = config.VerifyVar()
	if err != nil {
		log.Errorf("verify config failed, error:%s", err.Error())
		return
	}

	actionType := config.GetActionType()
	workspacePtr := config.GetCurrentWorkspace()
	templatePtr := s.queryTemplate(workspacePtr.GetTemplate())
	if templatePtr == nil {
		log.Errorf("query template failed, templateName:%s", workspacePtr.GetTemplate())
		err = cd.NewError(cd.Unexpected, "query template failed")
		return
	}

	switch actionType {
	case config.ModuleKey:
		err = s.runModuleAction(workspacePtr, templatePtr)
	case config.InitiatorKey:
		err = s.runInitiatorAction(workspacePtr, templatePtr)
	case config.ProjectKey:
		err = s.runProjectAction(workspacePtr, templatePtr)
	default:
		log.Errorf("unknown module type:%s", actionType)
		err = cd.NewError(cd.Unexpected, "unknown module type")
	}

	return
}

func (s *Cli) runModuleAction(workspacePtr *models.Workspace, templatePtr *models.Template) *cd.Error {
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
	result := s.SendEvent(ev)
	return result.Error()
}

func (s *Cli) runInitiatorAction(workspacePtr *models.Workspace, templatePtr *models.Template) *cd.Error {
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
	result := s.SendEvent(ev)
	return result.Error()
}

func (s *Cli) runProjectAction(workspacePtr *models.Workspace, templatePtr *models.Template) *cd.Error {
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
	result := s.SendEvent(ev)
	return result.Error()
}
