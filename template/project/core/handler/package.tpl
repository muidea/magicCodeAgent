{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}
package {{$packageName}}

{{define "serviceFuncDeclare"}}
    {{$serviceName := .Name}}
    // New{{$serviceName}}Handler 新建{{$serviceName}}Handler
    func New{{$serviceName}}Handler(clnt client.Client) {{$serviceName}}Handler {
    handler := &{{toLower $serviceName}}Impl{ {{toLower $serviceName}}Impl{batisClient: clnt} }

    return handler
    }
{{end}}

{{range $serviceList}}
    {{template "serviceFuncDeclare" .}}{{end}}
