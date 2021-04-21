{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$packageName := .PackageName }}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$actionList := $module.Action}}
{{$childList := $module.Child}}
{{$lowName := $module.Lower $moduleName}}
package biz

import (
"github.com/muidea/{{$projectName}}/core/handler"
{{range $mk, $mv := $childList}}
    {{$subModuleName := $mv.Name}}
    {{$lowName := $mv.Lower $subModuleName}}
    "github.com/muidea/{{$projectName}}/core/biz/{{$lowName}}"{{end}}
)

{{range $mk, $mv := $childList}}
    {{$subModuleName := $mv.Name}}
    {{$lowName := $mv.Lower $subModuleName}}
    func New{{$subModuleName}}(handler handler.{{$subModuleName}}Handler) *{{$lowName}}.{{$subModuleName}} {
    return {{$lowName}}.New{{$subModuleName}}(handler)
    }{{end}}