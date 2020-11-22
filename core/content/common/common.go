package common

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

type Define struct {
	projectName string
}

func New(projectName string) *Define {
	return &Define{projectName: projectName}
}

func (s *Define) getTemplate() (ret string) {
	ret = `
		{{$projectName := .ProjectName}}
		{{$module := .Module}}
		{{$moduleName := $module.Name}}
		{{$actionList := $module.Action}}
		{{$childList := $module.Child}}
		package common
		
		import (
			commonDef "github.com/muidea/magicCommon/def"
			"github.com/muidea/{{$projectName}}/model"
		)

		const (
			{{range $k, $v := $actionList}}
			{{$url := $module.FormatUrl $v }}
			{{$url}}
			{{end}}

			{{range $mk, $mv := $childList}}
			{{$subModuleName := $mv.Name}}
			{{$subModuleAction := $mv.Action}}
			{{range $ak, $av := $subModuleAction}}
			{{$url := $mv.FormatUrl $av }}
			{{$url}}
			{{end}}
			{{end}}
		)

		{{if not $module.IsPackage }}
		{{$jsonName := $module.Lower $moduleName}}
		type {{$moduleName}}Result struct {
			commonDef.Result
			{{$moduleName}} *model.{{$moduleName}}View ` + "`" + "json:\"{{$jsonName}}\"" + "`" + `
		}

		type {{$moduleName}}ListResult struct {
			commonDef.Result
			Total    int64                 ` + "`" + "json:\"total\"" + "`" + `
			{{$moduleName}} []*model.{{$moduleName}}View ` + "`" + "json:\"{{$jsonName}}\"" + "`" + ` 
		}
		{{end}}

		{{range $mk, $mv := $childList}}
		{{if not $mv.IsPackage}}
		{{$subModuleName := $mv.Name}}
		{{$jsonName := $mv.Lower $subModuleName}}
		type {{$subModuleName}}Result struct {
			commonDef.Result
			{{$subModuleName}} *model.{{$subModuleName}}View ` + "`" + "json:\"{{$jsonName}}\"" + "`" + `
		}

		type {{$subModuleName}}ListResult struct {
			commonDef.Result
			Total    int64                 ` + "`" + "json:\"total\"" + "`" + `
			{{$subModuleName}} []*model.{{$subModuleName}}View ` + "`" + "json:\"{{$jsonName}}\"" + "`" + ` 
		}
		{{end}}
		{{end}}
	`
	return
}

func (s *Define) Create(module *model.Module) (err error) {
	if len(module.Action) == 0 && len(module.Child) == 0 {
		return
	}

	existFlag, existErr := util.Exists("./common")
	if existErr != nil {
		err = existErr
		return
	}
	if !existFlag {
		err = os.Mkdir("./common", os.ModePerm)
		if err != nil {
			return
		}
	}

	filePath := fmt.Sprintf("./common/%s.go", strings.ToLower(module.Name))
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

		data := map[string]interface{}{"ProjectName": s.projectName, "Module": module}
		err = templateHandle.Execute(fileHandle, data)
	}()
	if err != nil {
		return
	}

	cmd := exec.Command("gofmt", "-w", filePath)
	err = cmd.Run()
	return
}
