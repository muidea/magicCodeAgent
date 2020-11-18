package service

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

type Instance struct {
	packageName string
	projectName string
}

func (s *Instance) getTemplate() (ret string) {
	ret = `
	{{$projectName := .ProjectName}}
	{{$packageName := .PackageName }}
	{{$module := .Module}}
	{{$moduleName := $module.Name}}
	{{$actionList := $module.Action}}
	{{$childList := $module.Child}}
	package {{$packageName}}
	
	import (
		"encoding/json"
		"log"
		"net/http"
		"strconv"
	
		commonDef "github.com/muidea/magicCommon/def"
		"github.com/muidea/magicCommon/foundation/net"
	
		"github.com/muidea/{{$projectName}}/common"
		"github.com/muidea/{{$projectName}}/model"
	)
	
	{{range $k, $v := $actionList}}
	{{$url := $module.UrlHandle $v }}
	{{$url}}
	{{end}}
	`

	return
}

func (s *Instance) Create(module *model.Module) (err error) {
	if len(module.Action) == 0 && len(module.Child) == 0 {
		return
	}

	filePath := fmt.Sprintf("./core/service/%s", s.packageName)
	existFlag, existErr := util.Exists(filePath)
	if existErr != nil {
		err = existErr
		return
	}
	if !existFlag {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return
		}
	}

	filePath = fmt.Sprintf("%s/%s.go", filePath, strings.ToLower(module.Name))
	existsFlag, existsErr := util.Exists(filePath)
	if existsErr != nil {
		err = fmt.Errorf("unexpect, fileName:%s", filePath)
		return
	}
	if existsFlag {
		return
	}

	func() {
		fileHandle, fileErr := os.Create(filePath)
		if fileErr != nil {
			err = fileErr
			return
		}
		defer fileHandle.Close()

		templateHandle, templateErr := template.New(module.Name).Parse(s.getTemplate())
		if templateErr != nil {
			err = templateErr
			return
		}

		data := map[string]interface{}{"ProjectName": s.projectName, "PackageName": module.LowerName(), "Module": module}
		err = templateHandle.Execute(fileHandle, data)
	}()
	if err != nil {
		return
	}

	cmd := exec.Command("gofmt", "-w", filePath)
	err = cmd.Run()
	return
}
