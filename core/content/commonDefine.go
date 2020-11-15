package content

import (
	"fmt"
	"github.com/muidea/magicToolkit/model"
	"os"
	"text/template"
)

type commonDefine struct {
}

func FormatUrl(url string) (ret string) {
	ret = fmt.Sprintf("\"%s\"", url)
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
			{{$url := FormatUrl $v }}
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
