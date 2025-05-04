package common

const PromptsInitator = "/internal/initators/prompts"

const SystemPromptsBasePath = "assets/template/prompts/system"

type PromptsRegistryHeper interface {
	GetPromptsRegistry() (ret PromptsRegistry)
}

type Prompt interface {
	Name() string
	Prompt() string
}

type PromptsRegistry interface {
	AddSection(prompts Prompt)
	// 添加自定section的prompts
	AddPrompts(section string, prompts Prompt)
	// 查询自定section的prompts
	GetSectionPrompts(section string) []Prompt
	// 查询执行自定section的prompts内容，以Md格式进行返回
	GetSectionPromptsContent(section string) string
}
