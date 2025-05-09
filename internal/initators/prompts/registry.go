package prompts

import (
	"fmt"
	"strings"
	"text/template"

	fu "github.com/muidea/magicCommon/foundation/util"

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

func (s *promptsRegistry) GetPromptsContent(modeName string) string {
	var builder strings.Builder
	builder.WriteString(s.getRolePrompts(modeName))
	builder.WriteString("\n\n")

	for _, section := range s.sections {
		builder.WriteString(section.GetPromptsContent())
		builder.WriteString("\n\n")
	}

	builder.WriteString(s.getSystemInformationPrompts())
	builder.WriteString("\n\n")

	builder.WriteString(s.getObjecttivePrompts())
	builder.WriteString("\n\n")

	builder.WriteString(s.getCustomInstructionsPrompts())
	builder.WriteString("\n\n")

	return builder.String()
}

func (s *promptsRegistry) getRolePrompts(modeName string) string {
	fileName := fmt.Sprintf("role_%s.md", modeName)
	contentVal, contentErr := util.GetTemplateContent(common.SectionPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get role prompts content failed:%s", contentErr.Error())
		return "You are magicAssistant, a highly skilled software engineer with extensive knowledge in many programming languages, frameworks, design patterns, and best practices."
	}

	return contentVal
}

func (s *promptsRegistry) getSystemInformationPrompts() string {
	fileName := "system-info.md"
	contentVal, contentErr := util.GetTemplateContent(common.SectionPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get system information prompts content failed:%s", contentErr.Error())
		return "Operating System: Linux 5.15\nDefault Shell: /bin/bash"
	}

	param := map[string]string{
		"OsName":      fu.GetOsName(),
		"SystemShell": fu.GetDefaultShell(),
	}

	tmpl, err := template.New("systemInfo").Parse(contentVal)
	if err != nil {
		log.Errorf("parse system info template failed:%s", err.Error())
		return contentVal
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, param)
	if err != nil {
		log.Errorf("execute system info template failed:%s", err.Error())
		return contentVal
	}

	return builder.String()
}

func (s *promptsRegistry) getObjecttivePrompts() string {
	fileName := "objective.md"
	contentVal, contentErr := util.GetTemplateContent(common.SectionPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get objective prompts content failed:%s", contentErr.Error())
		return ""
	}
	return contentVal
}

func (s *promptsRegistry) getCustomInstructionsPrompts() string {
	fileName := "custom-instructions.md"
	contentVal, contentErr := util.GetTemplateContent(common.SectionPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get custom prompts content failed:%s", contentErr.Error())
		return ""
	}
	return contentVal
}
