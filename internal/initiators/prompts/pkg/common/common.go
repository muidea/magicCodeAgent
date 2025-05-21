package common

import (
	fu "github.com/muidea/magicCommon/foundation/util"
)

const PromptsInitiator = "/internal/initiators/prompts"

const SectionPromptsBasePath = "assets/prompts/section"

const (
	Global  = "*"
	Service = "service"
	Web     = "web"
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
	AddSection(sectionPrompt Prompt)
	// 添加自定section的prompts
	AddPrompts(sectionPrompt string, prompt Prompt)
	GetPromptsContent(stage string, filterTags fu.StringSet) string
	ReadSectionPromptContent(sectionPrompt string) string
}
