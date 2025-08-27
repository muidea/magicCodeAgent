package common

import (
	"html/template"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/pkg/models"
	"github.com/muidea/magicCodeAgent/pkg/util"
	fu "github.com/muidea/magicCommon/foundation/util"
)

const PromptsInitiator = "/internal/initiators/prompts"

const innerSystemPromptsBasePath = "assets/system/prompts"

const (
	Global = "*"
)

type PromptsRegistryHelper interface {
	GetPromptsRegistry() (ret PromptsRegistry)
}

type Prompt interface {
	Name() string
	Prompt() string
	Tags() fu.StringSet
}

type PromptsRegistry interface {
	Reload(workspacePtr *models.Workspace)
	AddSection(sectionPrompt Prompt)
	AddPrompts(sectionPrompt string, prompt Prompt)
	GetStagePromptsContent(stage string) string
	ReadSectionPromptContent(sectionPrompt string, param map[string]string) string
}

func ReadInnerSystemPrompts(promptName string) (ret string, err error) {
	fileName := promptName + ".md"
	ret, err = util.GetTemplateContent(config.GetTemplateFS(), innerSystemPromptsBasePath, "", fileName)
	return
}

func ReadInnerSystemPromptsWithParam(promptName string, param map[string]string) (ret string, err error) {
	fileName := promptName + ".md"
	contentVal, contentErr := util.GetTemplateContent(config.GetTemplateFS(), innerSystemPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get inner system prompts content failed:%s", contentErr.Error())
		return
	}

	if len(param) == 0 {
		ret = contentVal
		return
	}

	tmpl, err := template.New(promptName).Parse(contentVal)
	if err != nil {
		log.Errorf("parse prompt %s template failed:%s", promptName, err.Error())
		ret = contentVal
		return
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, param)
	if err != nil {
		log.Errorf("execute prompt %s template failed:%s", promptName, err.Error())
		ret = contentVal
		return
	}

	ret = builder.String()
	return
}
