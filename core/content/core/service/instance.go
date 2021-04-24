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
	serviceName string
	projectName string
}

func (s *Instance) getTemplate() (ret string) {
	ret = `
{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
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
{{$url := $module.UrlHandle $serviceName $v }}
{{$url}}
{{end}}
`

	return
}

func (s *Instance) Create(module *model.Module, forceFlag bool) (err error) {
	if len(module.Action) == 0 && len(module.Child) == 0 {
		return
	}

	filePath := fmt.Sprintf("./core/service/%s", s.packageName)
	existFlag := util.Exists(filePath)
	if !existFlag {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return
		}
	}

	filePath = fmt.Sprintf("%s/%s.go", filePath, strings.ToLower(module.Name))
	if !forceFlag {
		existsFlag := util.Exists(filePath)
		if existsFlag {
			return
		}
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

		data := map[string]interface{}{"ProjectName": s.projectName, "ServiceName": s.serviceName, "PackageName": s.packageName, "Module": module}
		err = templateHandle.Execute(fileHandle, data)
	}()
	if err != nil {
		return
	}

	cmd := exec.Command("gofmt", "-w", filePath)
	err = cmd.Run()
	return
}
