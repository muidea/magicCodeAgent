package prompts

import (
	"strings"

	"github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
)

type promptsRegistry struct {
	promptsMap map[string][]common.Prompts
}

func NewPromptsRegistry() *promptsRegistry {
	return &promptsRegistry{
		promptsMap: make(map[string][]common.Prompts),
	}
}

func (s *promptsRegistry) AddPrompts(section string, promptsPtr common.Prompts) {
	if _, ok := s.promptsMap[section]; !ok {
		s.promptsMap[section] = []common.Prompts{}
	}
	s.promptsMap[section] = append(s.promptsMap[section], promptsPtr)
}

func (s *promptsRegistry) GetSectionPrompts(section string) []common.Prompts {
	return s.promptsMap[section]
}

func (s *promptsRegistry) GetSectionPromptsContent(section string) string {
	promptsList := s.GetSectionPrompts(section)
	if len(promptsList) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("## ")
	builder.WriteString(section)
	builder.WriteString("\n\n")

	for _, prompt := range promptsList {
		builder.WriteString("### ")
		builder.WriteString(prompt.Name())
		builder.WriteString("\n\n")
		builder.WriteString(prompt.Prompts())
		builder.WriteString("\n\n")
	}

	return builder.String()
}
