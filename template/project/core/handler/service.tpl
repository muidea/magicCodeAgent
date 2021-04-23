{{$apiVersion := .ApiVersion}}
{{$projectName := .ProjectName}}
{{$packageName := .PackageName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$actionList := $module.Action}}
{{$childList := $module.Child}}
package {{$packageName}}

import (
commonDef "github.com/muidea/magicCommon/def"
"github.com/muidea/{{$projectName}}/model"
)

const (
{{range $k, $v := $actionList}}
    {{$url := $module.FormatUrl $apiVersion $v "" }}{{$url}}{{end}}

{{range $mk, $mv := $childList}}
    {{$subModuleName := $mv.Name}}
    {{$subModuleAction := $mv.Action}}
    {{range $ak, $av := $subModuleAction}}
        {{$url := $mv.FormatUrl $apiVersion $av $moduleName }}{{$url}}{{end}}{{end}}
)

{{if not $module.IsPackage }}
    {{$jsonName := $module.Lower $moduleName}}
    type {{$moduleName}}Result struct {
    commonDef.Result
    {{$moduleName}} *model.{{$moduleName}}View `json:"{{$jsonName}}"`
    }

    type {{$moduleName}}ListResult struct {
    commonDef.Result
    Total    int64                 `json:"total"`
    {{$moduleName}} []*model.{{$moduleName}}View `json:"{{$jsonName}}"`
    }
{{end}}

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