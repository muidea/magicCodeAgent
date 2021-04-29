{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}
package {{$packageName}}

import "github.com/muidea/magicBatis/client"

{{define "serviceFuncDeclare"}}
    {{$serviceName := .Name}}
    // New{{$serviceName}} 新建{{$serviceName}}
    func New{{$serviceName}}(clnt client.Client) {{$serviceName}} {
    handler := &{{toLower $serviceName}}{ common{batisClient: clnt} }

    return handler
    }
{{end}}

{{range $serviceList}}
    {{template "serviceFuncDeclare" .}}{{end}}
