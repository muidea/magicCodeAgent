package main

import (
	"flag"
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

var batisService = "http://localhost:8080"

func init() {
	actionList = append(actionList, &action{
		name:         "createContent",
		defaultValue: false,
		description:  "create project content",
		funcPtr:      createContent,
	})

	flag.StringVar(&batisService, "BatisService", batisService, "magicBatis service address.")
}

func createContent(curPath string) {
	modelList, modelErr := queryModel(projectName, batisService)
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

func queryModel(projectName, batisService string) (ret []*remote.Object, err error) {
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
