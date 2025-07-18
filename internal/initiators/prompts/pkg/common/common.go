package common

import (
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

type PromptsRegistryHeper interface {
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
