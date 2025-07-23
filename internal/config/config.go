// Package config 提供magicCodeAgent的配置管理功能
package config

import (
	"embed"
	"os"
	"path"
	"time"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"

	"github.com/muidea/magicCodeAgent/pkg/models"
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

// Config 定义系统配置结构
type Config struct {
	BindPort     string              `json:"bindPort"`  // 绑定端口
	Workspace    []*models.Workspace `json:"workspace"` // 工作区列表
	Models       []*models.ModelInfo `json:"models"`    // AI模型列表
	AuthorConfig *models.AuthorInfo  `json:"authorConfig"`
	WorkspaceDir string              `json:"workspaceDir"`
	DebugFlag    bool                `json:"debugFlag"`
}

// init 初始化当前工作目录路径
func init() {
	currentWorkPath, _ = os.Getwd()

	defaultConfigPath := RootConfigDir()
	if !fp.Exist(defaultConfigPath) {
		os.MkdirAll(defaultConfigPath, 0755)
	}
	configPath := path.Join(defaultConfigPath, "config.json")
	_ = LoadConfig(configPath)
}

func RootConfigDir() string {
	return path.Join(fu.GetHomeDir(), ".magicCodeAgent")
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
	envBindPort, envOK := os.LookupEnv("BindPort")
	if envOK {
		return envBindPort
	}

	if sysConfig != nil {
		return sysConfig.BindPort
	}

	return listenPort
}

// Workspaces 获取工作区列表
// 返回: 工作区切片
func Workspaces() []*models.Workspace {
	if sysConfig != nil {
		return sysConfig.Workspace
	}

	return []*models.Workspace{}
}

// AddWorkspace 添加工作区配置
// workspace: 工作区配置
// 返回: 错误信息
func AddWorkspace(workspace *models.Workspace) (err error) {
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

func UpdateWorkspace(workspace *models.Workspace) (err error) {
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
				if !models.ValidStage(workspace.Stage) {
					workspace.Stage = sysConfig.Workspace[idx].Stage
				}
				if workspace.WorkspaceDir == "" {
					workspace.WorkspaceDir = sysConfig.Workspace[idx].WorkspaceDir
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
		newWorkspace := []*models.Workspace{}
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
				sysConfig.Workspace[idx].Status = models.EnableStatus
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
				sysConfig.Workspace[idx].Status = models.DisableStatus
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

func RefreshWorkspace(workspaceName string) (ret *models.Workspace, err error) {
	if sysConfig != nil {
		for idx := range sysConfig.Workspace {
			if sysConfig.Workspace[idx].Name == workspaceName {
				if !fp.Exist(sysConfig.Workspace[idx].WorkspaceDir) {
					sysConfig.Workspace[idx].Status = models.OfflineStatus
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

func GetWorkspace(workspaceName string) (ret *models.Workspace) {
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
func Models() []*models.ModelInfo {
	if sysConfig != nil {
		return sysConfig.Models
	}

	return []*models.ModelInfo{}
}

func AuthorConfig() *models.AuthorInfo {
	if sysConfig != nil {
		return sysConfig.AuthorConfig
	}

	return nil
}

func UpdateConfig(authConfigPtr *models.AuthorInfo, models []*models.ModelInfo) {
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

// CurrentWorkingPath 获取当前工作目录路径
// 返回: 路径字符串
func CurrentWorkingPath() string {
	return currentWorkPath
}

func DebugFlag() bool {
	if sysConfig != nil {
		return sysConfig.DebugFlag || debugWorkspace
	}

	return debugWorkspace
}
