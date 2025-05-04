package prompts

import (
	"strings"

	"github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/util"
	"github.com/muidea/magicCommon/foundation/log"
)

type promptsRegistry struct {
	sections []*Section
}

func NewPromptsRegistry() *promptsRegistry {
	return &promptsRegistry{
		sections: []*Section{},
	}
}

func (s *promptsRegistry) findSection(name string) (*Section, bool) {
	for _, section := range s.sections {
		if section.name == name {
			return section, true
		}
	}
	return nil, false
}

func (s *promptsRegistry) AddSection(sectionPrompt common.Prompt) {
	if _, ok := s.findSection(sectionPrompt.Name()); ok {
		log.Warnf("duplicate section prompts:%s", sectionPrompt.Name())
		return
	}
	s.sections = append(s.sections, &Section{
		name:        sectionPrompt.Name(),
		description: sectionPrompt.Prompt(),
		prompts:     []common.Prompt{},
	})
}

func (s *promptsRegistry) AddPrompts(sectionName string, prompt common.Prompt) {
	sectionObj, ok := s.findSection(sectionName)
	if !ok {
		log.Warnf("section prompts not found:%s", sectionName)
		return
	}
	sectionObj.AddPrompt(prompt)
}

func (s *promptsRegistry) GetSectionPrompts(sectionName string) []common.Prompt {
	sectionObj, ok := s.findSection(sectionName)
	if !ok {
		log.Warnf("section prompts not found:%s", sectionName)
		return nil
	}
	return sectionObj.GetPrompts()
}

func (s *promptsRegistry) GetSectionPromptsContent(sectionName string) string {
	sectionObj, ok := s.findSection(sectionName)
	if !ok {
		return ""
	}
	return sectionObj.GetPromptsContent()
}

func (s *promptsRegistry) GetPromptsContent() string {
	var builder strings.Builder
	builder.WriteString(s.getRolePrompts())
	builder.WriteString("\n\n")

	for _, section := range s.sections {
		builder.WriteString(section.GetPromptsContent())
		builder.WriteString("\n\n")
	}

	return builder.String()
}

func (s *promptsRegistry) getRolePrompts() string {
	fileName := "role.md"
	contentVal, contentErr := util.GetTemplateContent(common.SystemPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get role prompts content failed:%s", contentErr.Error())
		return "You are magicAssistant, a highly skilled software engineer with extensive knowledge in many programming languages, frameworks, design patterns, and best practices."
	}

	return contentVal
}
