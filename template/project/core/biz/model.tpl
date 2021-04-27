{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$model := .Model}}
{{$modelName := $model.Name}}
package {{$packageName}}

{{if eq $serviceName $modelName}}
type {{$serviceName}} struct {
{{toLower $serviceName}}Handler handler.{{$serviceName}}Handler
}

func New(
{{toLower $serviceName}}Handler handler.{{$serviceName}}Handler,
eventHub event.Hub,
backgroundRoutine *task.BackgroundRoutine,
) *{{$serviceName}} {
    {{toLower $serviceName}} := &{{$serviceName}}{}

    return {{toLower $serviceName}}
}

{{end}}

{{range $url := $model.Action}}
    func (s *{{$serviceName}}) {{formatAction $serviceName $modelName $url}}({{formatParam $serviceName $modelName $url}})({{formatResult $serviceName $modelName $url}}) {
    return s.{{toLower $serviceName}}Handler.{{formatAction $serviceName $modelName $url}}()
    }
{{end}}
