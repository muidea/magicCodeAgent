{{$projectName := .ProjectName}}
{{$packageName := .PackageName}}
{{$module := .Module}}
{{$apiVersion := .ApiVersion}}
{{$moduleName := $module.Name}}
{{$childList := $module.Child}}
package {{$packageName}}

import (
    commonDef "github.com/muidea/magicCommon/def"
    "github.com/muidea/{{$projectName}}/model"
)

const (
{{range $mk, $mv := $childList}}
    {{$subModuleName := $mv.Name}}
    {{$subModuleAction := $mv.Action}}
    {{range $ak, $av := $subModuleAction}}
        {{$url := $mv.FormatUrl $apiVersion $av $moduleName }}{{$url}}{{end}}{{end}}
)

{{range $mk, $mv := $childList}}
    {{if not $mv.IsPackage}}
        {{$subModuleName := $mv.Name}}
        {{$jsonName := $mv.Lower $subModuleName}}
        type {{$subModuleName}}Result struct {
        commonDef.Result
        {{$subModuleName}} *model.{{$subModuleName}}View `json:"{{$jsonName}}"`
        }

        type {{$subModuleName}}ListResult struct {
        commonDef.Result
        Total    int64                 `json:"total"`
        {{$subModuleName}} []*model.{{$subModuleName}}View `json:"{{$jsonName}}"`
        }
    {{end}}
{{end}}