{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}
package {{$packageName}}

{{range $mk, $mv := $serviceList}}
    {{commonServiceFuncCreate $mv}}{{end}}
