{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}
package {{$packageName}}

import "github.com/muidea/magicBatis/client"

{{define "serviceFuncDefine"}}
    {{$serviceName := .Name}}
    // New{{$serviceName}} 新建{{$serviceName}}
    func New{{$serviceName}}(clnt client.Client) {{$serviceName}} {
    return &{{toLower $serviceName}}{ common{batisClient: clnt} }
    }
{{end}}

{{range $sv := $serviceList}}
    {{if isEntity $sv}}
        {{with $sv}}
            {{template "serviceFuncDefine" .}}
        {{end -}}
    {{end -}}
{{end -}}