package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/muidea/magicToolkit/core/project"
)

var projectName = ""
var projectType = ""

func init() {
	actionList = append(actionList, &action{
		name:         "createProject",
		defaultValue: false,
		description:  "create new project,and initialize",
		funcPtr:      createProject,
	})

	flag.StringVar(&projectName, "projectName", projectName, "project name")
	flag.StringVar(&projectType, "projectType", projectType, "project type,[default,lite]")
}

func createProject(curPath string) {
	if projectName == "" {
		fmt.Printf("must be input project name")
		return
	}

	project := project.New(curPath)
	curProject, curErr := project.Create(projectName, projectType)
	if curErr != nil {
		log.Printf("Create project failed, err:%s", curErr.Error())
		return
	}
	curErr = saveProjectDefine(curPath, curProject)
	if curErr != nil {
		log.Printf("save project define failed, err:%s", curErr.Error())
		return
	}
}
