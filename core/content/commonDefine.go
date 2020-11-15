package content

import (
	"fmt"
	"github.com/muidea/magicToolkit/model"
	"os"
	"strings"
	"text/template"
)

type commonDefine struct {
}

func FormatUrl(name, url string) (ret string) {
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

func getTemplate() (ret string) {
	ret = `
		{{$projectName := .ProjectName}}
		{{$moduleName := .Module.Name}}
		{{$actionList := .Module.Action}}
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
		)

		type {{$moduleName}}Result struct {
			commonDef.Result
			{{$moduleName}} *model.{{$moduleName}}View
		}

		type {{$moduleName}}ListResult struct {
			commonDef.Result
			Total    int64                
			{{$moduleName}}s []*model.{{$moduleName}}View 
		}
	`
	return
}

func (s *commonDefine) Create(module *model.Module) (err error) {
	fileHandle, fileErr := os.Create(fmt.Sprintf("%s.go", module.Name))
	if fileErr != nil {
		err = fileErr
		return
	}
	defer fileHandle.Close()

	funcMap := template.FuncMap{"FormatUrl": FormatUrl}
	templateHandle, templateErr := template.New(module.Name).Funcs(funcMap).Parse(getTemplate())
	if templateErr != nil {
		err = templateErr
		return
	}

	data := map[string]interface{}{"ProjectName": "Test", "Module": module}
	err = templateHandle.Execute(fileHandle, data)
	return
}
