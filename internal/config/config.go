package config

import (
	"embed"
	"os"

	"github.com/muidea/magicCommon/foundation/log"
	fu "github.com/muidea/magicCommon/foundation/util"
)

var templateFS embed.FS

var currentWorkPath = ""
var listenPort = "8080"
var sysConfig *Config

type Config struct {
	BindPort  string   `json:"bindPort"`
	Workspace []string `json:"workspace"`
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

	sysConfig = cfg
	return
}

func BindPort() string {
	if sysConfig != nil {
		return sysConfig.BindPort
	}

	return listenPort
}

func Workspace() []string {
	if sysConfig != nil {
		return sysConfig.Workspace
	}

	return []string{}
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
