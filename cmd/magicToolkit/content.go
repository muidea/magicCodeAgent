package main

import (
	"github.com/muidea/magicToolkit/core/content"
	"log"
)

func init() {
	actionList = append(actionList, &action{
		name:         "createContent",
		defaultValue: false,
		description:  "create project content",
		funcPtr:      createContent,
	})
}

func createContent(curPath string) {
	content := content.New(curPath)
	curErr := content.Create()
	if curErr != nil {
		log.Printf("Create template failed, err:%s", curErr.Error())
		return
	}
}
