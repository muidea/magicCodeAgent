package config

import (
	"embed"
	"os"

	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"
)

var templateFS embed.FS

var currentConfigPath = ""
var currentWorkPath = ""
var listenPort = "8080"
var sysConfig *Config

type Workspace struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Path        string   `json:"path"`
	Ignores     []string `json:"ignores"`
}

type Config struct {
	BindPort  string       `json:"bindPort"`
	Workspace []*Workspace `json:"workspace"`
}

func init() {
	currentWorkPath, _ = os.Getwd()
}

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

func BindPort() string {
	if sysConfig != nil {
		return sysConfig.BindPort
	}

	return listenPort
}

func Workspaces() []*Workspace {
	if sysConfig != nil {
		return sysConfig.Workspace
	}

	return []*Workspace{}
}

func AddWorkspace(workspace *Workspace) {
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
		err := fu.SaveConfig(currentConfigPath, sysConfig)
		if err != nil {
			log.Errorf("save config failed, err:%s", err.Error())
		}
	}
}

func DelWorkspace(workspace *Workspace) {
	if sysConfig != nil {
		newWorkspace := []*Workspace{}
		for _, val := range sysConfig.Workspace {
			if val.Name != workspace.Name {
				newWorkspace = append(newWorkspace, val)
			}
		}

		if len(newWorkspace) != len(sysConfig.Workspace) {
			sysConfig.Workspace = newWorkspace
			err := fu.SaveConfig(currentConfigPath, sysConfig)
			if err != nil {
				log.Errorf("save config failed, err:%s", err.Error())
			}
		}
	}
}

func SetEmbedTemplateFS(fs embed.FS) {
	templateFS = fs
}

func GetTemplateFS() embed.FS {
	return templateFS
}

func EndpointName() string {
	return "magicToolkit"
}

func CurrentWorkPath() string {
	return currentWorkPath
}
