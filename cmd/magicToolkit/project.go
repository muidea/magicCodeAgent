package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/muidea/magicToolkit/core/project"
)

var projectName = ""

func init() {
	actionList = append(actionList, &action{
		name:         "createProject",
		defaultValue: false,
		description:  "create new project,and initialize",
		funcPtr:      createProject,
	})

	flag.StringVar(&projectName, "name", projectName, "project name")
}

func createProject(curPath string) {
	if projectName == "" {
		fmt.Printf("must be input project name")
		return
	}

	core := project.NewCore(curPath)
	curErr := core.Create(projectName)
	if curErr != nil {
		log.Printf("Create template failed, err:%s", curErr.Error())
		return
	}
}
