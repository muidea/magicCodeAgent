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

	module := module.New(curPath)
	curErr := module.Create("", moduleName)
	if curErr != nil {
		log.Printf("Create module failed, err:%s", curErr.Error())
		return
	}
}
