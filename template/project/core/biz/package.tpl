{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}
package {{$packageName}}

import (
"github.com/muidea/magicCommon/event"
"github.com/muidea/magicCommon/task"
"github.com/muidea/{{$projectName}}/core/handler"
{{range $sv := $serviceList}}
"github.com/muidea/{{$projectName}}/core/biz/{{toLower $sv.Name}}"{{end}}
)

{{define "serviceFuncDeclare"}}
    {{$serviceName := .Name}}
    func New{{$serviceName}}(
    handler handler.{{$serviceName}}Handler,
    eventHub event.Hub,
    backgroundRoutine *task.BackgroundRoutine,
    ) *{{toLower $serviceName}}.{{$serviceName}} {
    return {{toLower $serviceName}}.New(handler, eventHub, backgroundRoutine)
    }
{{end}}

{{range $serviceList}}
    {{template "serviceFuncDeclare" .}}{{end}}
