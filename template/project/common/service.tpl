{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceName := .Service.Name}}
{{$childList := .Service.Child}}
package {{$packageName}}

import (
    commonDef "github.com/muidea/magicCommon/def"
    "github.com/muidea/{{$projectName}}/model"
)

const (
{{range $mk, $mv := $childList}}
    {{$modelName := $mv.Name}}
    {{range $ak, $av := $mv.Action}}
        {{formatAction $serviceName $serviceName $av}}URL = "{{toLower $serviceName}}/{{toLower $av}}"{{end}}
{{end}}
)

{{range $mk, $mv := $childList}}
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