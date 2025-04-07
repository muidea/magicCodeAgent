package config

import (
	"embed"
	"os"
)

var templateFS embed.FS

var currentWorkPath = ""

func init() {
	currentWorkPath, _ = os.Getwd()
}

func BindPort() string {
	return ""
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
