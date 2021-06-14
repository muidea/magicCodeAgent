{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$apiVersion := .Service.ApiVersion}}
{{$serviceName := .Service.Name}}
{{$childList := .Service.Child}}
package {{$packageName}}

{{if isEntity .Service}}
import (
    commonDef "github.com/muidea/magicCommon/def"
    "github.com/muidea/{{$projectName}}/model"
)
{{end}}

const (
{{range $mk, $mv := $childList}}
    {{$modelName := $mv.Name}}
    {{range $ak, $av := $mv.Action}}
        {{formatAction $serviceName $modelName $av}}URL = "/{{$apiVersion}}/{{formatUrl $serviceName $modelName $av}}"
    {{- end -}}
{{end}}
)

{{if isEntity .Service}}
    {{range $mk, $mv := $childList}}
        {{$modelName := $mv.Name}}
        type {{$modelName}}View struct {
            // TODO
        }

        func (s *{{$modelName}}View) From{{$modelName}}(ptr *model.{{$modelName}}) {
            // TODO
        }

        type {{$modelName}}Lite struct {
            // TODO
        }

        func (s *{{$modelName}}Lite) From{{$modelName}}(ptr *model.{{$modelName}}) {
            // TODO
        }

        func (s *{{$modelName}}Lite) To{{$modelName}}()(ret *model.{{$modelName}}) {
            // TODO
            return
        }

        type {{$modelName}}Param struct {
            // TODO
        }

        func (s *{{$modelName}}Param) To{{$modelName}}() (ret *model.{{$modelName}}) {
            // TODO
            return
        }

        {{$jsonName := toCamelCase $modelName}}
        type {{$modelName}}Result struct {
        commonDef.Result
        {{$modelName}} *model.{{$modelName}} `json:"{{$jsonName}}"`
        }

        type {{$modelName}}ListResult struct {
        commonDef.Result
        {{$modelName}} []*{{$modelName}}Lite `json:"{{$jsonName}}"`
        }

        type {{$modelName}}StatisticResult struct {
        commonDef.Result
        Total    int64                 `json:"total"`
        {{$modelName}} []*model.{{$modelName}} `json:"{{$jsonName}}"`
        }
    {{end}}
{{end}}
