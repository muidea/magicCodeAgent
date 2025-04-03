package config

import "os"

var currentWorkPath = ""

func init() {
	currentWorkPath, _ = os.Getwd()
}

func EndpointName() string {
	return "magicToolkit"
}

func CurrentWorkPath() string {
	return currentWorkPath
}
