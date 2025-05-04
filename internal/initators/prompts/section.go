package prompts

import (
	"strings"

	"github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
)

type Section struct {
	name        string
	description string
	prompts     []common.Prompt
}

func (s *Section) Name() string {
	return s.name
}

func (s *Section) Prompt() string {
	return s.description
}

func (s *Section) AddPrompt(prompt common.Prompt) {
	s.prompts = append(s.prompts, prompt)
}

func (s *Section) GetPrompts() []common.Prompt {
	return s.prompts
}

func (s *Section) GetPromptsContent() string {
	var builder strings.Builder
	if s.description != "" {
		builder.WriteString(s.description)
		builder.WriteString("\n\n")
	}

	builder.WriteString("## ")
	builder.WriteString(s.name)
	builder.WriteString(" list\n\n")
	for _, prompt := range s.prompts {
		builder.WriteString(prompt.Prompt())
		builder.WriteString("\n\n")
	}

	return builder.String()
}
