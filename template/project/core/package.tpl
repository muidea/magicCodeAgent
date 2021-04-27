{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}
package {{$packageName}}

{{define "serviceFuncDeclare"}}
    {{$serviceName := .Name}}
    {{$modelList := .Child}}
    {{range $modelList}}
        {{$modelName := .Name}}
        func New{{$modelName}}Param() *{{$modelName}}Param {
        return &{{$modelName}}Param{}
        }

        func New{{$modelName}}Result() *{{$modelName}}Result {
        return &{{$modelName}}Result{}
        }

        func New{{$modelName}}ListResult() *{{$modelName}}ListResult {
        return &{{$modelName}}ListResult{}
        }{{end}}
{{end}}

{{range $serviceList}}
    {{template "serviceFuncDeclare" .}}{{end}}
