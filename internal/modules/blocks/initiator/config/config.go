package config

import (
	"errors"
	"flag"
)

var initiatorName string
var initiatorPath string

func init() {
	flag.StringVar(&initiatorName, "initiatorName", "", "initiator name，must start with letter and contain only alphanumeric chars")
	flag.StringVar(&initiatorPath, "initiatorPath", "", "initiator path, relative path based on the project root directory")
}

func Validate() error {
	if initiatorName == "" || initiatorPath == "" {
		return errors.New("initiatorName and initiatorPath must be specified")
	}
	return nil
}

func GetName() string {
	return initiatorName
}

func GetPath() string {
	return initiatorPath
}
