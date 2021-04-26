{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$apiVersion := .Service.ApiVersion}}
{{$serviceName := .Service.Name}}
{{$childList := .Service.Child}}
package {{$packageName}}

import (
    "fmt"
    "github.com/muidea/magicBatis/client"
    "github.com/muidea/magicBatis/def"
    "github.com/muidea/{{$projectName}}/model"
)

type {{$serviceName}}Handler interface {

}

type {{toLower $serviceName}}Handler struct {
    commonHandler
}

func New{{$serviceName}}Handler(clnt client.Client) {{$serviceName}}Handler {
    handler := &{{toLower $serviceName}}Handler{}
    handler.batisClnt = clnt

    return handler
}

{{range $mk, $mv := $childList}}
    {{$modelName := $mv.Name}}
    {{range $ak, $av := $mv.Action}}
        func (s *{{toLower $serviceName}}Handler){{formatAction $serviceName $serviceName $av}}
    {{end}}
{{end}}
