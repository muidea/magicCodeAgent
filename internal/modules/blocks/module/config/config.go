package config

import (
	"errors"
	"flag"
)

var moduleName string
var modulePath string

func init() {
	flag.StringVar(&moduleName, "moduleName", "", "module name，must start with letter and contain only alphanumeric chars")
	flag.StringVar(&modulePath, "modulePath", "", "module path, relative path based on the project root directory")
}

func Validate() error {
	if moduleName == "" {
		return errors.New("moduleName must be specified")
	}
	return nil
}

func GetName() string {
	return moduleName
}

func GetPath() string {
	return modulePath
}
