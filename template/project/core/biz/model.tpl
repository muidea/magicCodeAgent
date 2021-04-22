{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$packageName := .PackageName}}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$lowName := $module.Lower $moduleName}}
package {{$packageName}}

import (
casModel "github.com/muidea/magicCas/model"
{{if eq $serviceName $moduleName}}
    "github.com/muidea/{{$projectName}}/core/handler"
{{end}}
"github.com/muidea/{{$projectName}}/model"

commonDef "github.com/muidea/magicCommon/def"
)

{{if eq $serviceName $moduleName}}
    type {{$serviceName}} struct {
    {{$lowName}}Handler handler.{{$serviceName}}Handler
    }

    func New{{$serviceName}}({{$lowName}}Handler handler.{{$serviceName}}Handler) *{{$serviceName}} {
    return &{{$serviceName}}{ {{$lowName}}Handler: {{$lowName}}Handler }
    }
{{end}}

func (s *{{$serviceName}}) Filter{{$moduleName}}(filter *commonDef.Filter) (ret []*model.{{$moduleName}}, total int64, err error) {
return
}

func (s *{{$serviceName}}) Query{{$moduleName}}ByID(id int) (ret *model.{{$moduleName}}, err error) {
return
}

func (s *{{$serviceName}}) Create{{$moduleName}}(
{{$lowName}}Ptr *model.{{$moduleName}},
creater *casModel.Entity,
namespace string) (ret *model.{{$moduleName}}, err error) {
return
}

func (s *{{$serviceName}}) Update{{$moduleName}}(
{{$lowName}}Ptr *model.{{$moduleName}},
updater *casModel.Entity,
namespace string) (ret *model.{{$moduleName}}, err error) {
return
}

func (s *{{$serviceName}}) Delete{{$moduleName}}(id int) (ret *model.{{$moduleName}}, err error) {
return
}
