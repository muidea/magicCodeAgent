package prompts

import (
	"fmt"
	"strings"
	"text/template"

	fu "github.com/muidea/magicCommon/foundation/util"

	iipc "github.com/muidea/magicCodeAgent/internal/initators/prompts/pkg/common"
	"github.com/muidea/magicCodeAgent/pkg/common"
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

func (s *promptsRegistry) AddSection(sectionPrompt iipc.Prompt) {
	if _, ok := s.findSection(sectionPrompt.Name()); ok {
		log.Warnf("duplicate section prompts:%s", sectionPrompt.Name())
		return
	}
	s.sections = append(s.sections, &Section{
		name:        sectionPrompt.Name(),
		description: sectionPrompt.Prompt(),
		prompts:     []iipc.Prompt{},
	})
}

func (s *promptsRegistry) AddPrompts(sectionName string, prompt iipc.Prompt) {
	sectionObj, ok := s.findSection(sectionName)
	if !ok {
		log.Warnf("section prompts not found:%s", sectionName)
		return
	}
	sectionObj.AddPrompt(prompt)
}

func (s *promptsRegistry) GetSectionPrompts(sectionName string) []iipc.Prompt {
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

func (s *promptsRegistry) GetPromptsContent(modeName, workspaceDir string) string {
	var builder strings.Builder
	builder.WriteString(s.getRolePrompts(modeName))
	builder.WriteString("\n\n")

	for _, section := range s.sections {
		builder.WriteString(section.GetPromptsContent())
		builder.WriteString("\n\n")
	}

	if modeName == common.DesignStage {
		builder.WriteString(s.getMarkdownPrompts())
		builder.WriteString("\n\n")
	}

	builder.WriteString(s.getSystemInformationPrompts())
	builder.WriteString("\n\n")

	builder.WriteString(s.getObjecttivePrompts())
	builder.WriteString("\n\n")

	builder.WriteString(s.getCustomInstructionsPrompts())
	builder.WriteString("\n\n")

	builder.WriteString(s.getModelsPrompts())
	builder.WriteString("\n\n")

	builder.WriteString(s.getRulesPrompts(workspaceDir))
	builder.WriteString("\n\n")

	return builder.String()
}

func (s *promptsRegistry) ReadSectionPromptContent(sectionPrompt string) string {
	fileName := fmt.Sprintf("%s.md", sectionPrompt)
	contentVal, contentErr := util.GetTemplateContent(iipc.SectionPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("Failed to read section prompt content for %s: %v", sectionPrompt, contentErr)
		return ""
	}
	return contentVal
}

func (s *promptsRegistry) getRolePrompts(modeName string) string {
	fileName := fmt.Sprintf("role_%s.md", modeName)
	contentVal, contentErr := util.GetTemplateContent(iipc.SectionPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get role prompts content failed:%s", contentErr.Error())
		return "You are magicAssistant, a highly skilled software engineer with extensive knowledge in many programming languages, frameworks, design patterns, and best practices."
	}

	return contentVal
}

func (s *promptsRegistry) getSystemInformationPrompts() string {
	fileName := "system-info.md"
	contentVal, contentErr := util.GetTemplateContent(iipc.SectionPromptsBasePath, "", fileName)
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
	contentVal, contentErr := util.GetTemplateContent(iipc.SectionPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get objective prompts content failed:%s", contentErr.Error())
		return ""
	}
	return contentVal
}

func (s *promptsRegistry) getCustomInstructionsPrompts() string {
	fileName := "custom-instructions.md"
	contentVal, contentErr := util.GetTemplateContent(iipc.SectionPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get custom prompts content failed:%s", contentErr.Error())
		return ""
	}
	return contentVal
}

func (s *promptsRegistry) getMarkdownPrompts() string {
	fileName := "markdown.md"
	contentVal, contentErr := util.GetTemplateContent(iipc.SectionPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get markdown prompts content failed:%s", contentErr.Error())
		return ""
	}
	return contentVal
}

func (s *promptsRegistry) getModelsPrompts() string {
	fileName := "models.md"
	contentVal, contentErr := util.GetTemplateContent(iipc.SectionPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get models prompts content failed:%s", contentErr.Error())
		return ""
	}
	return contentVal
}

func (s *promptsRegistry) getRulesPrompts(workspaceDir string) string {
	fileName := "rules.md"
	contentVal, contentErr := util.GetTemplateContent(iipc.SectionPromptsBasePath, "", fileName)
	if contentErr != nil {
		log.Errorf("get rules prompts content failed:%s", contentErr.Error())
		return ""
	}
	param := map[string]string{
		"WorkspaceDir": workspaceDir,
	}

	tmpl, err := template.New("rules").Parse(contentVal)
	if err != nil {
		log.Errorf("parse rules template failed:%s", err.Error())
		return contentVal
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, param)
	if err != nil {
		log.Errorf("execute rules template failed:%s", err.Error())
		return contentVal
	}

	return builder.String()
}
