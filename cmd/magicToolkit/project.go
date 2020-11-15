package main

import (
	"fmt"
	"log"

	"github.com/muidea/magicToolkit/core/project"
)

func init() {
	actionList = append(actionList, &action{
		name:         "createProject",
		defaultValue: false,
		description:  "create new project,and initialize",
		funcPtr:      createProject,
	})
}

func createProject(curPath string) {
	if projectName == "" {
		fmt.Printf("must be input project name")
		return
	}

	project := project.New(curPath)
	curErr := project.Create(projectName)
	if curErr != nil {
		log.Printf("Create project failed, err:%s", curErr.Error())
		return
	}
}
