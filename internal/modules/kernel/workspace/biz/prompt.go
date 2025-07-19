package biz

import (
	"html/template"
	"strings"

	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"

	imbtmc "github.com/muidea/magicCodeAgent/internal/modules/blocks/templatemanager/pkg/common"
)

type PromptReader interface {
	ReadPrompt(workspaceConfig *models.Workspace, promptName string, param map[string]string) (string, error)
}

func (s *Workspace) ReadPrompt(workspaceConfig *models.Workspace, promptName string, param map[string]string) (string, error) {
	queryEvent := event.NewEvent(imbtmc.QueryTemplate, s.ID(), imbtmc.TemplatemanagerModule, nil, workspaceConfig.GetTemplate())
	queryResult := s.SendEvent(queryEvent)
	templateVal, templateErr := queryResult.Get()
	if templateErr != nil {
		log.Errorf("queryTemplate failed, name:%s, error:%s", promptName, templateErr.Error())
		return "", templateErr
	}

	templatePtr := templateVal.(*models.Template)
	contentVal, contentErr := templatePtr.ReadPrompt(promptName)
	if contentErr != nil {
		log.Errorf("readPrompt failed, name:%s, error:%s", promptName, contentErr.Error())
		return "", contentErr
	}

	tmplVal, tmplErr := template.New(promptName).Parse(contentVal)
	if tmplErr != nil {
		log.Errorf("parse system info template failed:%s", tmplErr.Error())
		return "", tmplErr
	}

	var builder strings.Builder
	execErr := tmplVal.Execute(&builder, param)
	if execErr != nil {
		log.Errorf("execute system info template failed:%s", execErr.Error())
		return "", execErr
	}

	return builder.String(), nil
}
