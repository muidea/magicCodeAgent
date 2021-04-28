{{$packageName := .PackageName}}
{{$projectName := .ProjectName}}
{{$serviceList := .ServiceList}}

package {{$packageName}}

type timerCheckTask struct {
    preTimeStamp *time.Time
    eventHub     event.Hub
}

func (s *timerCheckTask) Run() {
    log.Info("timerCheckTask task running..")

    current := time.Now()
    if s.preTimeStamp == nil {
    s.preTimeStamp = &current
    }

    preTime := *s.preTimeStamp
    s.preTimeStamp = &current

    log.Infof("timer check....., date:%s", current.Format("2006-01-02 15:04:05"))
    notifyTimer := &ed.TimerNotify{PreTime: preTime, CurTime: current}
    ev := event.NewEvent(ed.NotifyTimer, "/", "/#", notifyTimer)
    s.eventHub.Post(ev)
}

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
    sessionRegistry := session.CreateRegistry(nil)

    {{range $sv := $serviceList}}
    s.{{toLower $sv.Name}}Service.RegisterHandler(router)
    {{end}}

    s.backgroundRoutine.Timer(&timerCheckTask{eventHub: s.eventHub}, 24*time.Hour, 2*time.Hour)
}

// Teardown 销毁
func (s *Core) Teardown() {
}
