{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}

package {{$packageName}}

import (
log "github.com/cihub/seelog"
batisClient "github.com/muidea/magicBatis/client"
casToolkit "github.com/muidea/magicCas/toolkit"
"github.com/muidea/magicCommon/event"
"github.com/muidea/magicCommon/session"
"github.com/muidea/magicCommon/task"
engine "github.com/muidea/magicEngine"

"github.com/muidea/{{$projectName}}/config"
"github.com/muidea/{{$projectName}}/model"

"github.com/muidea/{{$projectName}}/core/biz"
"github.com/muidea/{{$projectName}}/core/biz/cas"
"github.com/muidea/{{$projectName}}/core/dao"
"github.com/muidea/{{$projectName}}/core/service"
"github.com/muidea/{{$projectName}}/core/service/base"

authorityBiz "github.com/muidea/{{$projectName}}/core/biz/authority"
"github.com/muidea/{{$projectName}}/core/service/authority"

{{range $sv := $serviceList}}
"github.com/muidea/{{$projectName}}/core/service/{{toLower $sv.Name}}"{{end}}
)


// New 新建Core
func New(clnt batisClient.Client,endpointName string) (ret *Core, err error) {
    err = model.InitializeModel(clnt)
    if err != nil {
        log.Errorf("initializeModel failed, err:%s", err.Error())
        return
    }

    eventHub := event.NewHub()
    backgroundRoutine := task.NewBackgroundRoutine()

    common := dao.NewCommon(clnt, config.CasService())

    {{range $sv := $serviceList}}
        {{toLower $sv.Name}} := dao.New{{$sv.Name}}(clnt)
    {{- end }}

    {{range $sv := $serviceList}}
        {{toLower $sv.Name}}Biz := biz.New{{$sv.Name}}({{toLower $sv.Name}}, eventHub, backgroundRoutine)
    {{- end }}

    {{range $sv := $serviceList}}
        {{toLower $sv.Name}}Service := service.New{{$sv.Name}}(common, {{toLower $sv.Name}}Biz)
    {{- end }}

    casBiz := cas.NewCas(config.CasService(), eventHub, backgroundRoutine)
    baseService := base.New(common, casBiz, config.FileService())

    authorityBiz := authorityBiz.New(config.CasService(), eventHub, backgroundRoutine)
    authorityService := authority.New(common, authorityBiz)

    core := &Core{ {{range $sv := $serviceList}}
            {{toLower $sv.Name}}Service: {{toLower $sv.Name}}Service, {{end}}
        baseService:       baseService,
        authorityService:  authorityService,
        eventHub:          eventHub,
        backgroundRoutine: backgroundRoutine,
    }

    ret = core
    return
}

// Core Core对象
type Core struct {
    {{range $sv := $serviceList}}
        {{toLower $sv.Name}}Service  *{{toLower $sv.Name}}.{{$sv.Name}}
    {{- end}}

    baseService       *base.Base
    authorityService *authority.Authority
    eventHub          event.Hub
    backgroundRoutine *task.BackgroundRoutine
}

// Startup 启动
func (s *Core) Startup(router engine.Router) {
    sessionRegistry := session.CreateRegistry(nil)
    casRouteRegistry := casToolkit.NewCasRegistry(s.baseService, router)
    roleRouteRegistry := casToolkit.NewRoleRegistry(s.baseService, router)

    s.baseService.BindRegistry(sessionRegistry, casRouteRegistry, roleRouteRegistry)
    s.baseService.RegisterHandler(router)

    s.authorityService.BindRegistry(sessionRegistry, roleRouteRegistry)
    s.authorityService.RegisterHandler(router)

    {{range $sv := $serviceList}}
        s.{{toLower $sv.Name}}Service.BindRegistry(sessionRegistry,roleRouteRegistry)
        s.{{toLower $sv.Name}}Service.RegisterHandler(router)
    {{end}}
}

// Teardown 销毁
func (s *Core) Teardown() {
}
