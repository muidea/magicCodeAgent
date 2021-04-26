package biz

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

type Biz struct {
	projectName string
}

func New(projectName string) *Biz {
	return &Biz{projectName: projectName}
}

func (s *Biz) getTemplate() (ret string) {
	ret = `
{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$packageName := .PackageName }}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$actionList := $module.Action}}
{{$childList := $module.Child}}
{{$lowName := $module.Lower $moduleName}}
package biz

import (
	"github.com/muidea/{{$projectName}}/core/handler"
	{{range $mk, $mv := $childList}}
	{{$subModuleName := $mv.Name}}
	{{$lowName := $mv.Lower $subModuleName}}
	"github.com/muidea/{{$projectName}}/core/biz/{{$lowName}}"
	{{end}}
)

{{range $mk, $mv := $childList}}
{{$subModuleName := $mv.Name}}
{{$lowName := $mv.Lower $subModuleName}}
func New{{$subModuleName}}(handler handler.{{$subModuleName}}Handler) *{{$lowName}}.{{$subModuleName}} {
	return {{$lowName}}.New{{$subModuleName}}(handler)
}
{{end}}
`

	return
}

func (s *Biz) Create(module *model.Module, forceFlag bool) (err error) {
	if len(module.Action) == 0 && len(module.Child) == 0 {
		return
	}

	filePath := fmt.Sprintf("./core/biz")
	existFlag := util.Exists(filePath)
	if !existFlag {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return
		}
	}

	filePath = fmt.Sprintf("%s/%s.go", filePath, "biz")
	existsFlag := util.Exists(filePath)
	if !existsFlag || forceFlag {
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

			data := map[string]interface{}{"ProjectName": s.projectName, "ServiceName": module.Name, "PackageName": util.ToLower(module.Name), "Module": module}
			err = templateHandle.Execute(fileHandle, data)
		}()
		if err != nil {
			return
		}

		cmd := exec.Command("gofmt", "-w", filePath)
		err = cmd.Run()
		if err != nil {
			return
		}
	}

	if !module.IsPackage() {
		instance := &Instance{projectName: s.projectName, serviceName: module.Name, packageName: util.ToLower(module.Name)}
		err = instance.Create(module, forceFlag)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	for _, val := range module.Child {
		if val.IsPackage() {
			continue
		}

		instance := &Instance{projectName: s.projectName, serviceName: module.Name, packageName: util.ToLower(module.Name)}
		err = instance.Create(val, forceFlag)
		if err != nil {
			return
		}
	}

	return
}
