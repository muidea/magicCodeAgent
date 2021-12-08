package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/muidea/magicToolkit/core/module"
)

var moduleName = ""

func init() {
	actionList = append(actionList, &action{
		name:         "createModule",
		defaultValue: false,
		description:  "create new project module",
		funcPtr:      createModule,
	})

	flag.StringVar(&moduleName, "moduleName", moduleName, "project module name")
}

func createModule(curPath string) {
	if moduleName == "" {
		fmt.Printf("must be input module name")
		return
	}

	projectDefine, projectErr := loadProjectDefine(curPath)
	if projectErr != nil {
		fmt.Printf("load project define failed, err:%s", projectErr.Error())
		return
	}

	module := module.New(curPath)
	curModule, curErr := module.Create(projectDefine.ProjectName, moduleName)
	if curErr != nil {
		log.Printf("Create module failed, err:%s", curErr.Error())
		return
	}

	projectDefine.Module = append(projectDefine.Module, curModule)
	projectErr = saveProjectDefine(curPath, projectDefine)
	if projectErr != nil {
		fmt.Printf("save project define failed, err:%s", projectErr.Error())
		return
	}
}
