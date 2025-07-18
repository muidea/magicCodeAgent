package config

import (
	"errors"
	"flag"
)

var projectName string
var projectDescription string

func init() {
	flag.StringVar(&projectName, "projectName", "", "project name，must start with letter and contain only alphanumeric chars")
	flag.StringVar(&projectDescription, "projectDescription", "", "project description")
}

func Validate() error {
	if projectName == "" {
		return errors.New("projectName must be specified")
	}
	return nil
}

func GetName() string {
	return projectName
}

func GetDescription() string {
	return projectDescription
}
