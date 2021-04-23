package main

import (
	"fmt"
	"github.com/muidea/magicBatis/def"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicOrm/provider/remote"
	"github.com/muidea/magicToolkit/core/content"
	"log"
	"net/http"
	"strings"
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
	if projectName == "" {
		fmt.Printf("must be input project name")
		return
	}

	modelList, modelErr := queryObject(projectName, batisService)
	if modelErr != nil {
		log.Printf("query project model failed, err:%s", modelErr.Error())
		return
	}

	content := content.New(curPath, modelList)
	curErr := content.Create(projectName, forceFlag)
	if curErr != nil {
		log.Printf("Create project content failed, err:%s", curErr.Error())
		return
	}
}

func queryObject(projectName, batisService string) (ret []*remote.Object, err error) {
	result := &def.QueryModelResult{}
	url := strings.Join([]string{batisService, "/model/query/"}, "")
	url = fmt.Sprintf("%s?source=%s", url, projectName)

	httpClient := &http.Client{}
	_, err = net.HTTPGet(httpClient, url, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("register model failed, reason:%s", result.Reason)
		return
	}

	ret = result.Model
	return
}
