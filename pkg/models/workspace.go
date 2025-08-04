package models

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"
)

const (
	EnableStatus  = "enable"  // 启用状态
	DisableStatus = "disable" // 禁用状态
	OfflineStatus = "offline" // 离线状态
)

const (
	UnknownStage = ""
	DesignStage  = "design"
	CodeStage    = "code"
	DebugStage   = "debug"
	PublishStage = "publish"
)

const (
	InnerKey = "inner"
	LocalKey = "local"
)

const TemplatePath = "assets/templates"

func ValidStage(stage string) bool {
	switch stage {
	case DesignStage, CodeStage, DebugStage, PublishStage:
		return true
	default:
		return false
	}
}

func GetStageReferenceStage(stage string) string {
	switch stage {
	case DesignStage:
		return ""
	case CodeStage, DebugStage:
		return DesignStage
	default:
		return ""
	}
}

func GetStageReferenceLabel(stage string) string {
	switch stage {
	case DesignStage:
		return "设计文件列表"
	case CodeStage, DebugStage:
		return "代码列表"
	default:
		return "参考列表"
	}
}

type Workspace struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Template     string   `json:"template"`
	Status       string   `json:"status"`
	Stage        string   `json:"stage"`
	WorkspaceDir string   `json:"workspaceDir"`
	TechStack    []string `json:"techStack"`
	Ignores      []string `json:"ignores"`
	UpdateTime   string   `json:"updateTime"`
}

func (s *Workspace) GetName() string {
	return s.Name
}

func (s *Workspace) GetRootPath() string {
	return s.WorkspaceDir
}

func (s *Workspace) GetWorkingPath() string {
	switch s.Stage {
	case DesignStage:
		return filepath.Join(s.WorkspaceDir, "design")
	default:
		return filepath.Join(s.WorkspaceDir, "code")
	}
}

func (s *Workspace) GetStageWorkingPath(stage string) string {
	switch stage {
	case DesignStage:
		return filepath.Join(s.WorkspaceDir, "design")
	case CodeStage, DebugStage:
		return filepath.Join(s.WorkspaceDir, "code")
	default:
		return ""
	}
}

func (s *Workspace) GetTemplate() string {
	return s.Template
}

func (s Workspace) GetStage() string {
	return s.Stage
}

func (s *Workspace) Available() bool {
	return s.Status == EnableStatus
}

func (s *Workspace) IsSamge(ptr *Workspace) bool {
	if ptr == nil {
		return false
	}

	return s.Name == ptr.Name &&
		s.Template == ptr.Template &&
		s.Status == ptr.Status &&
		s.Stage == ptr.Stage &&
		s.WorkspaceDir == ptr.WorkspaceDir &&
		fu.StringSlicesEqual(s.TechStack, ptr.TechStack)
}

func (s *Workspace) IsInnerTemplate() bool {
	items := strings.Split(s.Template, "/")
	return len(items) > 1 && items[0] == InnerKey
}

func (s *Workspace) FilterTage() fu.StringSet {
	ret := s.TechStack
	items := strings.Split(s.Template, "/")
	switch len(items) {
	case 2:
		ret = append(ret, items[0])
	case 3:
		ret = append(ret, items[1])
	}
	return ret
}

type Template struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	TechStack   []string `json:"techStack"`
	Ignores     []string `json:"ignores"`
	HostRootFS  fs.FS    `json:"-"`
}

func (s *Template) IsWeb() bool {
	return strings.Contains(s.Category, "web")
}

func (s *Template) IsService() bool {
	return strings.Contains(s.Category, "service")
}

/*
├── prompts 目录及内容必选
│   ├── decompose-code-task-block.md
│   ├── decompose-code-task-module.md
│   ├── decompose-code-task-project.md
│   ├── decompose-design-task-block.md
│   ├── decompose-design-task-module.md
│   ├── decompose-design-task-project.md
│   ├── project-struct.md
│   ├── stage-code-status.md
│   └── stage-design-status.md
├── template  目录及内容必选
│   ├── module  目录必选
│   └── project  目录必选
├── template.json 文件必选
└── tools 目录必选

*/
// 校验HostPath下的包含的内容是否符合上述要求，不符合要求返回error
func (s *Template) Validate() error {
	requiredDirs := []string{
		"prompts",
		"template/module",
		"template/project",
		"tools",
	}

	requiredFiles := []string{
		"prompts/decompose-code-task-block.md",
		"prompts/decompose-code-task-module.md",
		"prompts/decompose-code-task-project.md",
		"prompts/decompose-design-task-block.md",
		"prompts/decompose-design-task-module.md",
		"prompts/decompose-design-task-project.md",
		"prompts/project-struct.md",
		"prompts/stage-code-status.md",
		"prompts/stage-design-status.md",
		"template.json",
	}

	var missing []string

	// 检查必选目录
	for _, dir := range requiredDirs {
		finalPath := filepath.Join(TemplatePath, s.Category, s.Name, dir)
		if _, err := fs.Stat(s.HostRootFS, finalPath); os.IsNotExist(err) {
			missing = append(missing, dir)
		}
	}

	// 检查必选文件
	for _, file := range requiredFiles {
		finalPath := filepath.Join(TemplatePath, s.Category, s.Name, file)
		if _, err := fs.Stat(s.HostRootFS, finalPath); os.IsNotExist(err) {
			missing = append(missing, file)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("缺少必要的文件或目录: %v", missing)
	}

	return nil
}

func (s *Template) LoadTools() (ret []string) {
	finalPath := filepath.Join(TemplatePath, s.Category, s.Name, "tools")
	entityList, entityErr := fs.ReadDir(s.HostRootFS, finalPath)
	if entityErr != nil {
		log.Errorf("load tools failed, error:%s", entityErr.Error())
		return
	}
	for _, val := range entityList {
		if val.IsDir() {
			continue
		}

		finalPath = filepath.Join(TemplatePath, s.Category, s.Name, "tools", val.Name())
		fileContent, fileErr := fs.ReadFile(s.HostRootFS, finalPath)
		if fileErr != nil {
			log.Errorf("load tools failed, error:%s", fileErr.Error())
			continue
		}
		ret = append(ret, string(fileContent))
	}

	return
}

func (s *Template) ReadPrompt(name string) (ret string, err error) {
	finalPath := filepath.Join(TemplatePath, s.Category, s.Name, "prompts", name+".md")
	fileContent, fileErr := fs.ReadFile(s.HostRootFS, finalPath)
	if fileErr != nil {
		log.Errorf("load prompt failed, error:%s", fileErr.Error())
		err = fileErr
		return
	}
	ret = string(fileContent)
	return
}
