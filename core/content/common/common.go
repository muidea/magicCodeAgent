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

func (s *Define) FormatUrl(name, url string) (ret string) {
	switch url {
	case "query/all/":
		ret = fmt.Sprintf("QueryAll%sURL = \"/api/v1/%s/query/all/\"", strings.Title(name), strings.ToLower(name))
	case "query/":
		ret = fmt.Sprintf("Filter%sURL = \"/api/v1/%s/query/\"", strings.Title(name), strings.ToLower(name))
	case "query/:id":
		ret = fmt.Sprintf("Query%sURL = \"/api/v1/%s/query/:id\"", strings.Title(name), strings.ToLower(name))
	case "create/":
		ret = fmt.Sprintf("Create%sURL = \"/api/v1/%s/create/\"", strings.Title(name), strings.ToLower(name))
	case "update/:id":
		ret = fmt.Sprintf("Update%sURL = \"/api/v1/%s/update/:id\"", strings.Title(name), strings.ToLower(name))
	case "delete/:id":
		ret = fmt.Sprintf("Delete%sURL = \"/api/v1/%s/delete/:id\"", strings.Title(name), strings.ToLower(name))
	}

	return
}

func (s *Define) FormatJsonName(name string) string {
	return strings.ToLower(name)
}

func (s *Define) getTemplate() (ret string) {
	ret = `
		{{$projectName := .ProjectName}}
		{{$moduleName := .Module.Name}}
		{{$actionList := .Module.Action}}
		{{$childList := .Module.Child}}
		package common
		
		import (
			commonDef "github.com/muidea/magicCommon/def"
			"github.com/muidea/{{$projectName}}/model"
		)

		const (
			{{range $k, $v := $actionList}}
			{{$url := FormatUrl $moduleName $v }}
			{{$url}}
			{{end}}

			{{range $mk, $mv := $childList}}
			{{$subModuleName := $mv.Name}}
			{{$subModuleAction := $mv.Action}}
			{{range $ak, $av := $subModuleAction}}
			{{$url := FormatUrl $subModuleName $av }}
			{{$url}}
			{{end}}
			{{end}}
		)

		{{$jsonName := FormatJsonName $moduleName}}
		type {{$moduleName}}Result struct {
			commonDef.Result
			{{$moduleName}} *model.{{$moduleName}}View ` + "`" + "json:\"{{$jsonName}}\"" + "`" + `
		}

		type {{$moduleName}}ListResult struct {
			commonDef.Result
			Total    int64                 ` + "`" + "json:\"total\"" + "`" + `
			{{$moduleName}} []*model.{{$moduleName}}View ` + "`" + "json:\"{{$jsonName}}\"" + "`" + ` 
		}

		{{range $mk, $mv := $childList}}
		{{$subModuleName := $mv.Name}}
		{{$jsonName := FormatJsonName $subModuleName}}
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

	fileName := fmt.Sprintf("./common/%s.go", strings.ToLower(module.Name))
	func() {
		fileHandle, fileErr := os.Create(fileName)
		if fileErr != nil {
			err = fileErr
			return
		}
		defer fileHandle.Close()

		funcMap := template.FuncMap{"FormatUrl": s.FormatUrl, "FormatJsonName": s.FormatJsonName}
		templateHandle, templateErr := template.New(module.Name).Funcs(funcMap).Parse(s.getTemplate())
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

	cmd := exec.Command("gofmt", "-w", fileName)
	err = cmd.Run()
	return
}
