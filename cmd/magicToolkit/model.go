package main

import (
	"fmt"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicToolkit/core/model"
)

func init() {
	actionList = append(actionList, &action{
		name:         "initModel",
		defaultValue: false,
		description:  "initialize model",
		funcPtr:      InitModel,
	})
}

func InitModel(curPath string) {
	projectDefine, projectErr := loadProjectDefine(curPath)
	if projectErr != nil {
		fmt.Printf("load project define failed, err:%s", projectErr.Error())
		return
	}

	ptr := model.New(curPath)
	curErr := ptr.Create(projectDefine.ProjectName)
	if curErr != nil {
		log.Errorf("init models failed, err:%s", curErr.Error())
		return
	}
}
