{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$modelName := .Model.Name}}
package {{$packageName}}

import (
"net/http"

commonSession "github.com/muidea/magicCommon/session"
engine "github.com/muidea/magicEngine"

"github.com/muidea/{{$projectName}}/core/biz/{{toLower $serviceName}}"
"github.com/muidea/{{$projectName}}/core/dao"
)

{{if eq $serviceName $modelName}}
    type {{$serviceName}} struct {
    common dao.Common
    {{toLower $serviceName}}Biz *{{toLower $serviceName}}.{{$serviceName}}
    }

    func New(
    common dao.Common,
    {{toLower $serviceName}}Biz *{{toLower $serviceName}}.{{$serviceName}},
    ) *{{$serviceName}} {
    {{toLower $serviceName}} := &{{$serviceName}}{
    common:common,
    {{toLower $serviceName}}Biz:{{toLower $serviceName}}Biz,
    }

    return {{toLower $serviceName}}
    }

    func (s *{{$serviceName}}) BindRegistry(
    sessionRegistry commonSession.Registry) {
    }


    // Handle middleware handler
    func (s *{{$serviceName}}) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
    ctx.Next()
    }

    // RegisterHandler 注册路由
    func (s *{{$serviceName}}) RegisterHandler(router engine.Router) {
    }

{{end}}
