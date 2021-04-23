{{$apiVersion := .ApiVersion}}
{{$projectName := .ProjectName}}
{{$packageName := .PackageName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowModuleName := $module.Lower $moduleName}}
{{$actionList := $module.Action}}
{{$childList := $module.Child}}
package {{$packageName}}

import (
    "fmt"
    "github.com/muidea/magicBatis/client"
    "github.com/muidea/magicBatis/def"
    "github.com/muidea/{{$projectName}}/model"
)

type {{$moduleName}}Handler interface {

}

type {{$lowModuleName}}Handler struct {
    commonHandler
}

func New{{$moduleName}}Handler(clnt client.Client) {{$moduleName}}Handler {
    handler := &{{$lowModuleName}}Handler{}
    handler.batisClnt = clnt

    return handler
}

