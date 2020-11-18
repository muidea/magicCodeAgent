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

type Biz struct {
	packageName string
	serviceName string
	projectName string
}

func (s *Biz) getTemplate() (ret string) {
	ret = `
{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowName := $module.Lower $serviceName}}
package biz

import (
	"fmt"
	"github.com/muidea/magicBatis/def"
	casModel "github.com/muidea/magicCas/model"
	"github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/{{$projectName}}/core/handler"
	"github.com/muidea/{{$projectName}}/model"
	"time"

	commonDef "github.com/muidea/magicCommon/def"
)

type {{$serviceName}} struct {
	{{$lowName}}Handler handler.{{$serviceName}}Handler
}

func New{{$serviceName}}({{$lowName}}Handler handler.{{$serviceName}}Handler) *{{$serviceName}} {
	return &{{$serviceName}}{{{$lowName}}Handler: {{$lowName}}Handler}
}

func (s *{{$serviceName}}) Filter{{$moduleName}}(filter *commonDef.Filter) (ret []*model.{{$moduleName}}, total int64, err error) {
	return
}

func (s *{{$serviceName}}) Query{{$moduleName}}ByID(id int) (ret *model.{{{{$moduleName}}}}, err error) {
	return
}

func (s *{{$serviceName}}) Create{{$moduleName}}(
	{{$lowName}}Ptr *model.{{$moduleName}},
	creater *casModel.Entity,
	namespace string) (ret *model.{{$moduleName}}, err error) {
	return
}

func (s *{{$serviceName}}) Update{{$moduleName}}(
	{{$lowName}}Ptr *model.{{$moduleName}},
	creater *casModel.Entity,
	namespace string) (ret *model.{{$moduleName}}, err error) {
	return
}

func (s *{{$serviceName}}) Delete{{$moduleName}}(id int) (ret *model.{{$moduleName}}, err error) {
	return
}
`

	return
}

func (s *Biz) Create(module *model.Module) (err error) {
	if len(module.Action) == 0 && len(module.Child) == 0 {
		return
	}

	filePath := fmt.Sprintf("./core/service/%s/biz", s.packageName)
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

		data := map[string]interface{}{"ProjectName": s.projectName, "ServiceName": s.serviceName, "Module": module}
		err = templateHandle.Execute(fileHandle, data)
	}()
	if err != nil {
		return
	}

	cmd := exec.Command("gofmt", "-w", filePath)
	err = cmd.Run()
	return
}
