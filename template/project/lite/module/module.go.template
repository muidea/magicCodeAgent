package module

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/session"
	"github.com/muidea/magicCommon/task"
	engine "github.com/muidea/magicEngine"
)

type Module interface {
	ID() string
	BindEventHub(eventHub event.Hub)
	BindBackgroundRoutine(backgroundRoutine *task.BackgroundRoutine)
	BindRegistry(sessionRegistry session.Registry)
	Startup(endpointName string, router engine.Router)
	Shutdown()
}

var moduleList []Module

func Register(module Module) {
	moduleList = append(moduleList, module)
}

func GetList() []Module {
	return moduleList
}
