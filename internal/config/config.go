// Package config 提供magicCodeAgent的配置管理功能
package config

import (
	"embed"
	"os"
	"time"

	"github.com/muidea/magicCommon/foundation/log"
	fp "github.com/muidea/magicCommon/foundation/path"
	fu "github.com/muidea/magicCommon/foundation/util"
)

// templateFS 用于嵌入模板文件系统
var templateFS embed.FS

// 全局配置变量
var (
	currentConfigPath = ""               // 当前配置文件路径
	currentWorkPath   = ""               // 当前工作目录路径
	currentEndpoint   = "magicCodeAgent" // 默认端点名称
	listenPort        = "8080"           // 默认监听端口
	sysConfig         *Config            // 系统配置实例
)

// 状态常量定义
const (
	EnableStatus  = "enable"  // 启用状态
	DisableStatus = "disable" // 禁用状态
	OfflineStatus = "offline" // 离线状态
)

// Workspace 定义工作区配置结构
type Workspace struct {
	Name        string   `json:"name"`        // 工作区名称
	Description string   `json:"description"` // 工作区描述
	Status      string   `json:"status"`      // 工作区状态(enable/disable/offline)
	Path        string   `json:"path"`        // 工作区路径
	TechStack   []string `json:"techStack"`   // 技术栈列表
	Ignores     []string `json:"ignores"`     // 忽略文件/目录列表
	UpdateTime  string   `json:"updateTime"`
}

// Model 定义AI模型配置结构
type Model struct {
	Name     string `json:"name"`     // 模型名称
	BaseUrl  string `json:"baseUrl"`  // 模型基础URL
	APIKey   string `json:"apiKey"`   // 模型API密钥
	Provider string `json:"provider"` // 模型提供商
}

// Config 定义系统配置结构
type Config struct {
	BindPort  string       `json:"bindPort"`  // 绑定端口
	Workspace []*Workspace `json:"workspace"` // 工作区列表
	Models    []*Model     `json:"models"`    // AI模型列表
}

// init 初始化当前工作目录路径
func init() {
	currentWorkPath, _ = os.Getwd()
}

// LoadConfig 加载配置文件
// cfgFile: 配置文件路径
// 返回: 错误信息
func LoadConfig(cfgFile string) (err error) {
	if cfgFile == "" {
		return
	}

	cfg := &Config{}
	err = fu.LoadConfig(cfgFile, cfg)
	if err != nil {
		log.Errorf("load config failed, err:%s", err.Error())
		return
	}

	currentConfigPath = cfgFile
	sysConfig = cfg
	return
}

// SetEndpoint 设置当前端点名称
// endpointName: 端点名称
func SetEndpoint(endpointName string) {
	currentEndpoint = endpointName
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
				if workspace.Status == "" {
					workspace.Status = sysConfig.Workspace[idx].Status
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
func Models() []*Model {
	if sysConfig != nil {
		return sysConfig.Models
	}

	return []*Model{}
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

// CurrentWorkPath 获取当前工作目录路径
// 返回: 路径字符串
func CurrentWorkPath() string {
	return currentWorkPath
}
