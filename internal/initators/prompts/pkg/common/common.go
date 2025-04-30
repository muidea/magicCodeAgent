package common

const PromptsInitator = "/internal/initators/prompts"

type Prompts interface {
	Name() string
	Prompts() string
}

type PromptsRegistry interface {
	// 添加自定section的prompts
	AddPrompts(section string, promptsPtr Prompts)
	// 查询自定section的prompts
	GetSectionPrompts(section string) []Prompts
	// 查询执行自定section的prompts内容，以Md格式进行返回
	GetSectionPromptsContent(section string) string
}
