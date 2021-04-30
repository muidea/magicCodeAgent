{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$apiVersion := .Service.ApiVersion}}
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
        {{formatAction $serviceName $serviceName $av}}URL = "/{{$apiVersion}}/{{toLower $serviceName}}/{{toLower $av}}"{{end}}
{{end}}
)

{{range $mk, $mv := $childList}}
    {{$modelName := $mv.Name}}
    type {{$modelName}}Lite struct {
    }

    func (s *{{$modelName}}Lite) From{{$modelName}}(ptr *model.{{$modelName}}) {
    }

    type {{$modelName}}Param struct {
    }

    func (s *{{$modelName}}Param) To{{$modelName}}() *model.{{$modelName}} {
        return &model.{{$modelName}}{}
    }

    type {{$modelName}}View struct {
    }

    func (s *{{$modelName}}View) From{{$modelName}}(ptr *model.{{$modelName}}) {
    }

    func (s *{{$modelName}}View) To{{$modelName}}() *model.{{$modelName}} {
        return &model.{{$modelName}}{}
    }

    {{$jsonName := toLower $modelName}}
    type {{$modelName}}Result struct {
    commonDef.Result
    {{$modelName}} *{{$modelName}}View `json:"{{$jsonName}}"`
    }

    type {{$modelName}}ListResult struct {
    commonDef.Result
    {{$modelName}} []*{{$modelName}}Lite `json:"{{$jsonName}}"`
    }

    type {{$modelName}}StatisticResult struct {
    commonDef.Result
    Total    int64                 `json:"total"`
    {{$modelName}} []*{{$modelName}}View `json:"{{$jsonName}}"`
    }
{{end}}