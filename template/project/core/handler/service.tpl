{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$apiVersion := .Service.ApiVersion}}
{{$serviceName := .Service.Name}}
{{$childList := .Service.Child}}
package {{$packageName}}

import (
"github.com/muidea/magicBatis/def"
"github.com/muidea/{{$projectName}}/model"
)

// {{$serviceName}}Handler {{toLower $serviceName}} handler
type {{$serviceName}}Handler interface {
    {{range $mv := $childList}}
        {{$modelName := $mv.Name}}
        {{range $url := $mv.Action}}
            {{formatAction $serviceName $modelName $url}}({{formatParam $serviceName $modelName $url}})({{formatResult $serviceName $modelName $url}}){{end}}
    {{end}}
}

type {{toLower $serviceName}}Handler struct {
    commonHandler
}

{{range $mv := $childList}}
    {{$modelName := $mv.Name}}
    {{range $url := $mv.Action}}
        func (s *{{toLower $serviceName}}Handler){{formatAction $serviceName $modelName $url}}({{formatParam $serviceName $modelName $url}})({{formatResult $serviceName $modelName $url}}){
            return
        }
    {{end}}
{{end}}