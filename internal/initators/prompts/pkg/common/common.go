package common

const PromptsInitator = "/internal/initators/prompts"

const SectionPromptsBasePath = "assets/prompts/section"

type PromptsRegistryHeper interface {
	GetPromptsRegistry() (ret PromptsRegistry)
}

type Prompt interface {
	Name() string
	Prompt() string
}

type PromptsRegistry interface {
	AddSection(sectionPrompt Prompt)
	// 添加自定section的prompts
	AddPrompts(sectionPrompt string, prompt Prompt)
	// 查询自定section的prompts
	GetSectionPrompts(sectionPrompt string) []Prompt
	// 查询执行自定section的prompts内容，以Md格式进行返回
	GetSectionPromptsContent(sectionPrompt string) string
	GetPromptsContent(stage string) string
	ReadSectionPromptContent(sectionPrompt string, sectionParam map[string]string) string
}
