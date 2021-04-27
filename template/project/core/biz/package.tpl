{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}
package {{$packageName}}

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
