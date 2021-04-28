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
    {{toLower $serviceName}}Biz *{{toLower $serviceName}}.{{$serviceName}}
    }

    func New(
    commonHandler handler.CommonHandler,
    {{toLower $serviceName}}Biz *{{toLower $serviceName}}.{{$serviceName}},
    ) *{{$serviceName}} {
    {{toLower $serviceName}} := &{{$serviceName}}{
    commonHandler:commonHandler,
    {{toLower $serviceName}}Biz:{{toLower $serviceName}}Biz,
    }

    return {{toLower $serviceName}}
    }

    func (s *{{$serviceName}}) BindRegistry(
    sessionRegistry commonSession.Registry,
    casRouteRegistry casToolkit.CasRegistry,
    privateRouteRegistry casToolkit.PrivateRegistry) {

    s.sessionRegistry = sessionRegistry
    s.casRouteRegistry = casRouteRegistry
    s.privateRouteRegistry = privateRouteRegistry
    }


    // Handle middleware handler
    func (s *{{$serviceName}}) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
    curSession := s.sessionRegistry.GetSession(res, req)

    sessionInfo := curSession.GetSessionInfo()
    sessionInfo.Scope = commonCommon.ShareSession

    values := req.URL.Query()
    values = sessionInfo.Encode(values)
    req.URL.RawQuery = values.Encode()

    ctx.Next()
    }

    // RegisterHandler 注册路由
    func (s *{{$serviceName}}) RegisterHandler(router engine.Router) {
    }

{{end}}
