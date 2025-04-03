package main

import (
	"flag"
	"fmt"

	log "github.com/cihub/seelog"

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

	ptr := module.New(curPath)
	curModule, curErr := ptr.Create(projectDefine.ProjectName, moduleName)
	if curErr != nil {
		log.Errorf("Create module failed, err:%s", curErr.Error())
		return
	}

	found := false
	for _, val := range projectDefine.Module {
		if val.Name == curModule.Name {
			found = true
		}
	}
	if !found {
		projectDefine.Module = append(projectDefine.Module, curModule)
	}
	projectErr = saveProjectDefine(curPath, projectDefine)
	if projectErr != nil {
		fmt.Printf("save project define failed, err:%s", projectErr.Error())
		return
	}
}
