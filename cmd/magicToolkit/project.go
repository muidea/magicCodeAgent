package main

import (
	"flag"
	"fmt"
	"log"
	"path"

	"github.com/muidea/magicToolkit/core/project"
	"github.com/muidea/magicToolkit/core/util"
)

var projectType = ""

func init() {
	actionList = append(actionList, &action{
		name:         "createProject",
		defaultValue: false,
		description:  "create new project,and initialize",
		funcPtr:      createProject,
	})

	flag.StringVar(&projectType, "projectType", projectType, "project type")
}

func createProject(curPath string) {
	if projectName == "" {
		fmt.Printf("must be input project name")
		return
	}

	definePath := path.Join(curPath, "doc", "define.json")
	modelDefine, modelErr := util.GetModelDefine(projectName, definePath)
	if modelErr != nil {
		log.Printf("get project model define failed, err:%s", modelErr.Error())
		return
	}

	objectList, objectErr := util.QueryObject(projectName, batisService)
	if objectErr != nil {
		log.Printf("query project object failed, err:%s", objectErr.Error())
		return
	}

	project := project.New(curPath, modelDefine, objectList)
	curErr := project.Create(projectName, projectType)
	if curErr != nil {
		log.Printf("Create project failed, err:%s", curErr.Error())
		return
	}
}
