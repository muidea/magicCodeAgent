{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}
package {{$packageName}}

import (
"github.com/muidea/{{$projectName}}/core/dao"

{{range $sv := $serviceList}}
    {{toLower $sv.Name}}Biz "github.com/muidea/{{$projectName}}/core/biz/{{toLower $sv.Name}}"{{end}}

{{range $sv := $serviceList}}
    "github.com/muidea/{{$projectName}}/core/service/{{toLower $sv.Name}}"{{end}}
)

{{range $sv :=  $serviceList}}
    {{$serviceName := $sv.Name}}
    func New{{$serviceName}}(
    dao dao.Common,
    biz *{{toLower $serviceName}}Biz.{{$serviceName}},
    ) *{{toLower $serviceName}}.{{$serviceName}} {
    return {{toLower $serviceName}}.New(dao, biz)
    }
{{end}}