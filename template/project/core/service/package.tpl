{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}
package {{$packageName}}

import (
{{range $sv := $serviceList}}
    {{toLower $sv.Name}}Biz "github.com/muidea/{{$projectName}}/core/biz/{{toLower $sv.Name}}"{{end}}

{{range $sv := $serviceList}}
    "github.com/muidea/{{$projectName}}/core/service/{{toLower $sv.Name}}"{{end}}
)

{{range $sv :=  $serviceList}}
    {{$serviceName := $sv.Name}}
    func New{{$serviceName}}(
    handler handler.CommonHandler,
    biz *{{toLower $serviceName}}Biz.{{$serviceName}},
    ) *{{toLower $serviceName}}.{{$serviceName}} {
    return {{toLower $serviceName}}.New(handler, biz)
    }
{{end}}