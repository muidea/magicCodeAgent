{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}

package {{$packageName}}

import (
log "github.com/cihub/seelog"
batisClient "github.com/muidea/magicBatis/client"
"github.com/muidea/magicCommon/event"
"github.com/muidea/magicCommon/task"
engine "github.com/muidea/magicEngine"
"github.com/muidea/{{$projectName}}/core/biz"
"github.com/muidea/{{$projectName}}/core/handler"
"github.com/muidea/{{$projectName}}/core/service"
{{range $sv := $serviceList}}
"github.com/muidea/{{$projectName}}/core/service/{{toLower $sv.Name}}"{{end}}
)


// New 新建Core
func New(clnt batisClient.Client) (*Core, error) {
    err := model.InitializeModel(clnt)
    if err != nil {
    log.Errorf("initializeModel failed, err:%s", err.Error())
    return nil, err
    }

    eventHub := event.NewHub()
    backgroundRoutine := task.NewBackgroundRoutine()

    commonHandler := handler.NewCommonHandler(clnt)

    {{range $sv := $serviceList}}
        {{toLower $sv.Name}}Handler := handler.New{{$sv.Name}}(clnt)
    {{end}}

    {{range $sv := $serviceList}}
        {{toLower $sv.Name}}Biz := biz.New{{$sv.Name}}({{toLower $sv.Name}}Handler, eventHub, backgroundRoutine)
    {{end}}

    {{range $sv := $serviceList}}
        {{toLower $sv.Name}}Service := service.New{{$sv.Name}}(commonHandler, {{toLower $sv.Name}}Biz)
    {{end}}

    core := &Core{ {{range $sv := $serviceList}}
            {{toLower $sv.Name}}Service: {{toLower $sv.Name}}Service, {{end}}
        eventHub:          eventHub,
        backgroundRoutine: backgroundRoutine,
    }

    return core, nil
}

// Core Core对象
type Core struct {
    {{range $sv := $serviceList}}
        {{toLower $sv.Name}}Service  *{{toLower $sv.Name}}.{{$sv.Name}} {{end}}
    eventHub          event.Hub
    backgroundRoutine *task.BackgroundRoutine
}

// Startup 启动
func (s *Core) Startup(router engine.Router) {
    {{range $sv := $serviceList}}
    s.{{toLower $sv.Name}}Service.RegisterHandler(router)
    {{end}}
}

// Teardown 销毁
func (s *Core) Teardown() {
}
