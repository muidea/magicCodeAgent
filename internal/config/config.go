// Package config 提供magicCodeAgent的配置管理功能
package config

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	"github.com/muidea/magicCodeAgent/pkg/common"
)

// templateFS 用于嵌入模板文件系统
var templateFS embed.FS

// 全局配置变量
var (
	currentConfigPath = ""               // 当前配置文件路径
	currentWorkPath   = ""               // 当前工作目录路径
	currentEndpoint   = "magicCodeAgent" // 默认端点名称
	currentBuildTime  = ""
	currentGitCommit  = ""
	currentVersion    = "0.0.1"
	listenPort        = "8090" // 默认监听端口
	sysConfig         *Config  // 系统配置实例
	debugWorkspace    = true
)

// 状态常量定义
const (
	EnableStatus  = "enable"  // 启用状态
	DisableStatus = "disable" // 禁用状态
	OfflineStatus = "offline" // 离线状态
)

// Workspace 定义工作区配置结构
type Workspace struct {
	Name        string       `json:"name"`        // 工作区名称
	Description string       `json:"description"` // 工作区描述
	Type        string       `json:"type"`
	Status      string       `json:"status"` // 工作区状态(enable/disable/offline)
	Stage       string       `json:"stage"`
	Path        string       `json:"path"`      // 工作区路径
	TechStack   fu.StringSet `json:"techStack"` // 技术栈列表
	Ignores     fu.StringSet `json:"ignores"`   // 忽略文件/目录列表
	UpdateTime  string       `json:"updateTime"`
}

func (s *Workspace) GetName() string {
	return s.Name
}

func (s *Workspace) GetStage() string {
	return s.Stage
}

func (s *Workspace) GetWorkspaceDir() string {
	return s.Path
}

func (s *Workspace) GetWorkspaceType() string {
	frontendTechs := fu.StringSet{"javascript", "html", "css", "react", "vue", "umijs", "antd", "typescript"}

	for _, tech := range s.TechStack {
		lowerTech := strings.ToLower(tech)
		for _, frontend := range frontendTechs {
			if strings.Contains(lowerTech, frontend) {
				return common.Web
			}
		}
	}

	return common.Service
}

func (s *Workspace) GetReferenceStagePath(stageVal string) (string, string, string) {
	previousLabel := ""
	previousStage := ""
	switch stageVal {
	case common.DesignStage:
		previousStage = common.DesignStage
		previousLabel = "设计文件列表"
	case common.CodeStage:
		previousLabel = "代码列表"
		previousStage = common.CodeStage
	default:
		return "", "", ""
	}

	return previousStage, previousLabel, filepath.Join(s.Path, previousStage)
}

func (s *Workspace) GetStageWorkPath() string {
	switch s.Stage {
	case common.CodeStage, common.TestStage:
		return filepath.Join(s.Path, common.CodeStage)
	}
	return filepath.Join(s.Path, s.Stage)
}

func (s *Workspace) GetStagePath(stage string) string {
	return filepath.Join(s.Path, stage)
}

func (s *Workspace) Available() bool {
	return s.Status == EnableStatus
}

// ModelInfo 定义AI模型配置结构
type ModelInfo struct {
	Name     string `json:"name"`     // 模型名称
	BaseUrl  string `json:"baseUrl"`  // 模型基础URL
	APIKey   string `json:"apiKey"`   // 模型API密钥
	Provider string `json:"provider"` // 模型提供商
	Model    string `json:"model"`
}

type AuthorConfigInfo struct {
	Author string `json:"author"`
	EMail  string `json:"email"`
}

// Config 定义系统配置结构
type Config struct {
	BindPort     string            `json:"bindPort"`  // 绑定端口
	Workspace    []*Workspace      `json:"workspace"` // 工作区列表
	Models       []*ModelInfo      `json:"models"`    // AI模型列表
	AuthorConfig *AuthorConfigInfo `json:"authorConfig"`
	WorkspaceDir string            `json:"workspaceDir"`
	DebugFlag    bool              `json:"debugFlag"`
}

// init 初始化当前工作目录路径
func init() {
	currentWorkPath, _ = os.Getwd()

	defaultConfigPath := path.Join(fu.GetHomeDir(), ".magicCodeAgent")
	if !fp.Exist(defaultConfigPath) {
		os.MkdirAll(defaultConfigPath, 0755)
	}
	configPath := path.Join(defaultConfigPath, "config.json")
	_ = LoadConfig(configPath)
}

// LoadConfig 加载配置文件
// cfgFile: 配置文件路径
// 返回: 错误信息
func LoadConfig(cfgFile string) (err error) {
	if cfgFile == "" {
		return
	}

	cfg := &Config{}
	if !fp.Exist(cfgFile) {
		cfg.BindPort = listenPort
		err = fu.SaveConfig(cfgFile, cfg)
		if err != nil {
			log.Errorf("save config failed, err:%s", err.Error())
			return
		}
	} else {
		err = fu.LoadConfig(cfgFile, cfg)
		if err != nil {
			log.Errorf("load config failed, err:%s", err.Error())
			return
		}
	}

	currentConfigPath = cfgFile
	sysConfig = cfg
	return
}

func SetApplicationInfo(endpointName, buildTime, gitCommit, version string) {
	currentEndpoint = endpointName
	currentBuildTime = buildTime
	currentGitCommit = gitCommit
	currentVersion = version
}

// SetEndpoint 设置当前端点名称
// endpointName: 端点名称
func SetEndpoint(endpointName string) {
	currentEndpoint = endpointName
}

func EnableDebug() {
	debugWorkspace = true
}

// BindPort 获取绑定端口
// 返回: 端口号字符串
func BindPort() string {
	if sysConfig != nil {
		return sysConfig.BindPort
	}

	return listenPort
}

// Workspaces 获取工作区列表
// 返回: 工作区切片
func Workspaces() []*Workspace {
	if sysConfig != nil {
		return sysConfig.Workspace
	}

	return []*Workspace{}
}

// AddWorkspace 添加工作区配置
// workspace: 工作区配置
// 返回: 错误信息
func AddWorkspace(workspace *Workspace) (err error) {
	if sysConfig != nil {
		existFlag := false
		for idx := range sysConfig.Workspace {
			if sysConfig.Workspace[idx].Name == workspace.Name {
				sysConfig.Workspace[idx] = workspace
				existFlag = true
			}
		}
		if !existFlag {
			sysConfig.Workspace = append(sysConfig.Workspace, workspace)
		}
		err = fu.SaveConfig(currentConfigPath, sysConfig)
		if err != nil {
			log.Errorf("save config failed, err:%s", err.Error())
		}
	}

	return
}

func UpdateWorkspace(workspace *Workspace) (err error) {
	if sysConfig != nil {
		for idx := range sysConfig.Workspace {
			if sysConfig.Workspace[idx].Name == workspace.Name {
				if workspace.Name == "" {
					workspace.Name = sysConfig.Workspace[idx].Name
				}
				if workspace.Description == "" {
					workspace.Description = sysConfig.Workspace[idx].Description
				}
				if workspace.Status == "" {
					workspace.Status = sysConfig.Workspace[idx].Status
				}
				if !common.ValidStage(workspace.Stage) {
					workspace.Stage = sysConfig.Workspace[idx].Stage
				}
				if workspace.Path == "" {
					workspace.Path = sysConfig.Workspace[idx].Path
				}
				if workspace.TechStack == nil {
					workspace.TechStack = sysConfig.Workspace[idx].TechStack
				}
				if workspace.Ignores == nil {
					workspace.Ignores = sysConfig.Workspace[idx].Ignores
				}
				workspace.UpdateTime = time.Now().Format(fu.CSTLayout)
				sysConfig.Workspace[idx] = workspace
				err = fu.SaveConfig(currentConfigPath, sysConfig)
				if err != nil {
					log.Errorf("save config failed,err:%s", err.Error())
				}
			}
		}
	}
	return
}

// DelWorkspace 删除工作区配置
// workspace: 工作区配置
// 返回: 错误信息
func DelWorkspace(workspaceName string) (err error) {
	if sysConfig != nil {
		newWorkspace := []*Workspace{}
		for _, val := range sysConfig.Workspace {
			if val.Name != workspaceName {
				newWorkspace = append(newWorkspace, val)
			}
		}

		if len(newWorkspace) != len(sysConfig.Workspace) {
			sysConfig.Workspace = newWorkspace
			err = fu.SaveConfig(currentConfigPath, sysConfig)
			if err != nil {
				log.Errorf("save config failed, err:%s", err.Error())
			}
		}
	}

	return
}

// EnableWorkspace 启用工作区
// workspace: 工作区配置
// 返回: 错误信息
func EnableWorkspace(workspaceName string) (err error) {
	if sysConfig != nil {
		for idx := range sysConfig.Workspace {
			if sysConfig.Workspace[idx].Name == workspaceName {
				sysConfig.Workspace[idx].Status = EnableStatus
				err = fu.SaveConfig(currentConfigPath, sysConfig)
				if err != nil {
					log.Errorf("save config failed, err:%s", err.Error())
				}

				return
			}
		}
	}

	return
}

// DisableWorkspace 禁用工作区
// workspace: 工作区配置
// 返回: 错误信息
func DisableWorkspace(workspaceName string) (err error) {
	if sysConfig != nil {
		for idx := range sysConfig.Workspace {
			if sysConfig.Workspace[idx].Name == workspaceName {
				sysConfig.Workspace[idx].Status = DisableStatus
				err = fu.SaveConfig(currentConfigPath, sysConfig)
				if err != nil {
					log.Errorf("save config failed, err:%s", err.Error())
				}

				return
			}
		}
	}

	return
}

func RefreshWorkspace(workspaceName string) (ret *Workspace, err error) {
	if sysConfig != nil {
		for idx := range sysConfig.Workspace {
			if sysConfig.Workspace[idx].Name == workspaceName {
				if !fp.Exist(sysConfig.Workspace[idx].Path) {
					sysConfig.Workspace[idx].Status = OfflineStatus
				}
				sysConfig.Workspace[idx].UpdateTime = time.Now().Format(fu.CSTLayout)
				err = fu.SaveConfig(currentConfigPath, sysConfig)
				if err != nil {
					log.Errorf("save config failed, err:%s", err.Error())
					return
				}
				ret = sysConfig.Workspace[idx]
				return
			}
		}
	}
	return
}

func GetWorkspace(workspaceName string) (ret *Workspace) {
	if sysConfig != nil {
		for idx := range sysConfig.Workspace {
			if sysConfig.Workspace[idx].Name == workspaceName {
				ret = sysConfig.Workspace[idx]
				return
			}
		}
	}

	return
}

// Models 获取AI模型列表
// 返回: 模型切片
func Models() []*ModelInfo {
	if sysConfig != nil {
		return sysConfig.Models
	}

	return []*ModelInfo{}
}

func AuthorConfig() *AuthorConfigInfo {
	if sysConfig != nil {
		return sysConfig.AuthorConfig
	}

	return nil
}

func UpdateConfig(authConfigPtr *AuthorConfigInfo, models []*ModelInfo) {
	if sysConfig == nil {
		return
	}
	if authConfigPtr != nil {
		sysConfig.AuthorConfig = authConfigPtr
	}
	if models != nil {
		sysConfig.Models = models
	}

	err := fu.SaveConfig(currentConfigPath, sysConfig)
	if err != nil {
		log.Errorf("save config failed, err:%s", err.Error())
	}
}

func GetWorkspaceDir() string {
	return sysConfig.WorkspaceDir
}

// SetEmbedTemplateFS 设置嵌入模板文件系统
// fs: 文件系统实例
func SetEmbedTemplateFS(fs embed.FS) {
	templateFS = fs
}

// GetTemplateFS 获取嵌入模板文件系统
// 返回: 文件系统实例
func GetTemplateFS() embed.FS {
	return templateFS
}

// EndpointName 获取当前端点名称
// 返回: 端点名称字符串
func EndpointName() string {
	return currentEndpoint
}

func BuildTime() string {
	return currentBuildTime
}

func GitCommit() string {
	return currentGitCommit
}

func Version() string {
	return currentVersion
}

// CurrentWorkPath 获取当前工作目录路径
// 返回: 路径字符串
func CurrentWorkPath() string {
	return currentWorkPath
}

func DebugFlag() bool {
	if sysConfig != nil {
		return sysConfig.DebugFlag || debugWorkspace
	}

	return debugWorkspace
}

var validServiceModuleParentPaths = []string{
	filepath.Join("internal", "modules"),
	filepath.Join("internal", "modules", "blocks"),
	filepath.Join("internal", "modules", "kernel"),
}

var validWebModuleParentPaths = []string{
	filepath.Join("src", "pages"),
}

// VerifyModuleWorkPath 校验工作路径是否合法，并返回路径中ProjectName
// 1.工作路径的父级路径必须是以下几种：xxx/[ProjectName]/internal/modules,xxx/[ProjectName]/internal/modules/blocks,xxx/xxx/[ProjectName]/internal/modules/kernel
// 2. 其他路径非法
// 3. 工作路径不能是空
// 4. 工作路径不能是相对路径
// 5. 非法路径举例： /home/xxx/internal/modules/xxx/
func VerifyModuleWorkPath(projectType, pathVal string) (prefixPath string, err *cd.Error) {
	if projectType == common.Web {
		prefixPath, err = verifyWorkPath(pathVal, validWebModuleParentPaths)
		return
	}
	prefixPath, err = verifyWorkPath(pathVal, validServiceModuleParentPaths)
	return
}

var validInitiatorParentPaths = []string{
	filepath.Join("internal", "initiators"),
}

// VerifyInitiatorWorkPath 校验工作路径是否合法
// 1.工作路径的父级路径必须是以下几种：xxx/xxx/internal/initiators
// 2. 其他路径非法
// 3. 工作路径不能是空
// 4. 工作路径不能是相对路径
// 5. 非法路径举例： /home/xxx/internal/modules/xxx/
func VerifyInitiatorWorkPath(pathVal string) (prefixPath string, err *cd.Error) {
	prefixPath, err = verifyWorkPath(pathVal, validInitiatorParentPaths)
	return
}

func verifyWorkPath(pathVal string, validParentPath []string) (prefixPath string, err *cd.Error) {
	if pathVal == "" {
		return "", cd.NewError(cd.IllegalParam, "work path cannot be empty")
	}

	// 检查相对路径
	if !filepath.IsAbs(pathVal) {
		return "", cd.NewError(cd.IllegalParam, "work path must be absolute path")
	}

	// 标准化路径
	normalizedPath := filepath.Clean(pathVal)
	for _, validPath := range validParentPath {
		// 检查路径是否以validPath结尾
		if strings.HasSuffix(normalizedPath, validPath) {
			prefixPath = validPath
			return
		}
	}

	return "", cd.NewError(cd.IllegalParam,
		fmt.Sprintf("invalid work path: %s, parent path must be one of: %v", normalizedPath, validParentPath))
}

// VerifyProjectWorkPath 校验项目工作目录
// 1. 当前目录下无子目录或子文件
// 2. 当天目录下同时存在如下文件go.mod,go.sum,Makefile, 及application,internal,pkg,vendor目录
func VerifyProjectWorkPath(projectType, pathVal string) *cd.Error {
	// 空路径视为合法
	if pathVal == "" {
		return cd.NewError(cd.IllegalParam, "project work path cannot be empty")
	}

	// 检查目录是否为空
	entries, readErr := os.ReadDir(pathVal)
	if readErr != nil {
		return cd.NewError(cd.IllegalParam, fmt.Sprintf("failed to read directory: %v", readErr))
	}
	if len(entries) == 0 {
		return nil
	}

	// 检查相对路径
	if !filepath.IsAbs(pathVal) {
		return cd.NewError(cd.IllegalParam, "work path must be absolute path")
	}

	if projectType == common.Service {
		// 检查必需文件
		requiredFiles := []string{"go.mod", "go.sum", "Makefile"}
		for _, file := range requiredFiles {
			filePath := filepath.Join(pathVal, file)
			if !fp.Exist(filePath) {
				return cd.NewError(cd.IllegalParam, fmt.Sprintf("required file %s not found in project", file))
			}
		}

		// 检查必需目录
		requiredDirs := []string{"application", "internal", "pkg", "vendor"}
		for _, dir := range requiredDirs {
			dirPath := filepath.Join(pathVal, dir)
			if !fp.Exist(dirPath) {
				return cd.NewError(cd.IllegalParam, fmt.Sprintf("required directory %s not found in project", dir))
			}
		}
		return nil
	}

	if projectType == common.Web {
		requiredDirs := []string{"package.json", "yarn.lock", "src", "node_modules", "tsconfig.json"}
		for _, dir := range requiredDirs {
			dirPath := filepath.Join(pathVal, dir)
			if !fp.Exist(dirPath) {
				return cd.NewError(cd.IllegalParam, fmt.Sprintf("required directory %s not found in project", dir))
			}
		}
		return nil
	}

	return nil
}
