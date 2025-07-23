package biz

import (
	"os"
	"strings"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/base/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/cli/config"
	"github.com/muidea/magicCodeAgent/internal/modules/kernel/cli/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/models"

	imbcc "github.com/muidea/magicCodeAgent/internal/modules/blocks/commander/pkg/common"
	imbic "github.com/muidea/magicCodeAgent/internal/modules/blocks/initiator/pkg/common"
	imbmc "github.com/muidea/magicCodeAgent/internal/modules/blocks/module/pkg/common"
	imbpc "github.com/muidea/magicCodeAgent/internal/modules/blocks/project/pkg/common"
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
	case config.CommanderKey:
		err = s.runCommandAction(workspacePtr, templatePtr)
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
		ev = event.NewEvent(imbmc.CreateNewModule, s.ID(), imbmc.ModuleModule, nil, workspacePtr)
	case config.DestroyKey:
		ev = event.NewEvent(imbmc.DestroyModule, s.ID(), imbmc.ModuleModule, nil, workspacePtr)
	case config.RefreshKey:
		ev = event.NewEvent(imbmc.RefreshModule, s.ID(), imbmc.ModuleModule, nil, workspacePtr)
	default:
		log.Warnf("unknown execute action:%s", config.GetExecuteAction())
		return cd.NewError(cd.Unexpected, "unknown execute action")
	}

	ev.SetData("template", templatePtr)
	result := s.SendEvent(ev)
	return result.Error()
}

func (s *Cli) runInitiatorAction(workspacePtr *models.Workspace, templatePtr *models.Template) *cd.Error {
	var ev event.Event
	switch config.GetExecuteAction() {
	case config.CreateKey:
		ev = event.NewEvent(imbic.CreateNewInitiator, s.ID(), imbic.InitiatorModule, nil, workspacePtr)
	case config.DestroyKey:
		ev = event.NewEvent(imbic.DestroyInitiator, s.ID(), imbic.InitiatorModule, nil, workspacePtr)
	case config.RefreshKey:
		ev = event.NewEvent(imbic.RefreshInitiator, s.ID(), imbic.InitiatorModule, nil, workspacePtr)
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
		ev = event.NewEvent(imbpc.CreateNewProject, s.ID(), imbpc.ProjectModule, nil, workspacePtr)
	case config.DestroyKey:
		ev = event.NewEvent(imbpc.DestroyProject, s.ID(), imbpc.ProjectModule, nil, workspacePtr)
	case config.RefreshKey:
		ev = event.NewEvent(imbpc.RefreshProject, s.ID(), imbpc.ProjectModule, nil, workspacePtr)
	default:
		log.Warnf("unknown execute action:%s", config.GetExecuteAction())
		return cd.NewError(cd.Unexpected, "unknown execute action")
	}

	ev.SetData("template", templatePtr)
	result := s.SendEvent(ev)
	return result.Error()
}

func (s *Cli) runCommandAction(workspacePtr *models.Workspace, _ *models.Template) *cd.Error {
	var commandEvent event.Event
	switch config.GetExecuteAction() {
	case config.ExecuteKey:
		paramPtr, paramErr := s.parseCommand(workspacePtr)
		if paramErr != nil {
			log.Errorf("parse command failed, error:%s", paramErr.Error())
			return paramErr
		}
		commandEvent = event.NewEvent(imbcc.CommandRun, s.ID(), imbcc.CommanderModule, nil, workspacePtr.Name)
		commandEvent.SetData("command", paramPtr)
	default:
		log.Warnf("unknown execute action:%s", config.GetExecuteAction())
		return cd.NewError(cd.Unexpected, "unknown execute action")
	}

	commandResult := s.SendEvent(commandEvent)
	resultVal, resultErr := commandResult.Get()
	if resultErr != nil {
		log.Errorf("run command failed, error:%s", resultErr.Error())
		return resultErr
	}
	commanderResult := resultVal.(*imbcc.CommanderResult)
	if commanderResult.StatusMessage != "" {
		os.Stdout.WriteString(commanderResult.StatusMessage)
	}
	if commanderResult.ErrorMessage != "" {
		os.Stderr.WriteString(commanderResult.ErrorMessage)
	}
	return nil
}

func (s *Cli) parseCommand(workspacePtr *models.Workspace) (ret *imbcc.CommanderParam, err *cd.Error) {
	command := config.GetExecuteCommand()
	if command == "" {
		err = cd.NewError(cd.InvalidParameter, "command is empty")
		return
	}

	items := strings.Split(strings.TrimSpace(command), " ")
	ret = &imbcc.CommanderParam{
		WorkspaceDir: workspacePtr.GetRootPath(),
		Command:      items[0],
		Args:         items[1:],
	}
	return
}
