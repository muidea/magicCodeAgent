package biz

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
{{$packageName := .PackageName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowName := $module.Lower $moduleName}}
package {{$packageName}}

import (
	casModel "github.com/muidea/magicCas/model"
{{if eq $serviceName $moduleName}}
	"github.com/muidea/{{$projectName}}/core/handler"
{{end}}
	"github.com/muidea/{{$projectName}}/model"

	commonDef "github.com/muidea/magicCommon/def"
)

{{if eq $serviceName $moduleName}}
type {{$serviceName}} struct {
	{{$lowName}}Handler handler.{{$serviceName}}Handler
}

func New{{$serviceName}}({{$lowName}}Handler handler.{{$serviceName}}Handler) *{{$serviceName}} {
	return &{{$serviceName}}{ {{$lowName}}Handler: {{$lowName}}Handler }
}
{{end}}

func (s *{{$serviceName}}) Filter{{$moduleName}}(filter *commonDef.Filter) (ret []*model.{{$moduleName}}, total int64, err error) {
	return
}

func (s *{{$serviceName}}) Query{{$moduleName}}ByID(id int) (ret *model.{{$moduleName}}, err error) {
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
	updater *casModel.Entity,
	namespace string) (ret *model.{{$moduleName}}, err error) {
	return
}

func (s *{{$serviceName}}) Delete{{$moduleName}}(id int) (ret *model.{{$moduleName}}, err error) {
	return
}
`

	return
}

func (s *Instance) Create(module *model.Module, forceFlag bool) (err error) {
	if len(module.Action) == 0 && len(module.Child) == 0 {
		return
	}

	filePath := fmt.Sprintf("./core/biz/%s", s.packageName)
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
	if !forceFlag {
		existsFlag, existsErr := util.Exists(filePath)
		if existsErr != nil {
			err = fmt.Errorf("unexpect, fileName:%s", filePath)
			return
		}
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

		data := map[string]interface{}{"ProjectName": s.projectName, "PackageName": s.packageName, "ServiceName": s.serviceName, "Module": module}
		err = templateHandle.Execute(fileHandle, data)
	}()
	if err != nil {
		return
	}

	cmd := exec.Command("gofmt", "-w", filePath)
	err = cmd.Run()
	return
}
