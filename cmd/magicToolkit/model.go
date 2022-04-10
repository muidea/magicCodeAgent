package main

import (
	"flag"
	"fmt"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicToolkit/core/model"
)

var bothFlag = false

func init() {
	actionList = append(actionList, &action{
		name:         "initModel",
		defaultValue: false,
		description:  "initialize model",
		funcPtr:      InitModel,
	})

	flag.BoolVar(&bothFlag, "bothFlag", bothFlag, "init model both backend & frontend")
}

func InitModel(curPath string) {
	projectDefine, projectErr := loadProjectDefine(curPath)
	if projectErr != nil {
		fmt.Printf("load project define failed, err:%s", projectErr.Error())
		return
	}

	ptr := model.New(curPath, bothFlag)
	curErr := ptr.Create(projectDefine.ProjectName)
	if curErr != nil {
		log.Errorf("init models failed, err:%s", curErr.Error())
		return
	}
}
