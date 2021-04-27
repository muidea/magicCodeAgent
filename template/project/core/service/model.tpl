{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$model := .Model}}
{{$modelName := $model.Name}}
package {{$packageName}}

import (
"github.com/muidea/{{$projectName}}/core/biz/{{toLower $serviceName}}"
"github.com/muidea/{{$projectName}}/core/handler"
)

{{if eq $serviceName $modelName}}
    type {{$serviceName}} struct {
    commonHandler handler.CommonHandler
    {{toLower $serviceName}}Biz {{toLower $serviceName}}.{{$serviceName}}
    }

    func New(
    commonHandler handler.CommonHandler,
    {{toLower $serviceName}}Biz {{toLower $serviceName}}.{{$serviceName}},
    ) *{{$serviceName}} {
    {{toLower $serviceName}} := &{{$serviceName}}{ commonHandler:commonHandler, {{toLower $serviceName}}Biz:{{toLower $serviceName}}Biz }

    return {{toLower $serviceName}}
    }

{{end}}
