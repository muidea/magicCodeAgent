package content

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

type Core struct {
	projectName string
}

func (s *Core) getTemplate() (ret string) {
	ret = `
{{$projectName := .ProjectName}}
{{$moduleList := .Define.Module}}
package core

import (
	"log"

	engine "github.com/muidea/magicEngine"

	batisClient "github.com/muidea/magicBatis/client"
	casToolkit "github.com/muidea/magicCas/toolkit"
	"github.com/muidea/magicCommon/session"

	"github.com/muidea/{{$projectName}}/config"
	"github.com/muidea/{{$projectName}}/core/handler"
	"github.com/muidea/{{$projectName}}/core/service/base"
	"github.com/muidea/{{$projectName}}/core/service/system"

	"github.com/muidea/{{$projectName}}/model"
	{{range $k, $v := $moduleList}}
	{{$lowName := $v.Lower $v.Name}}
	"github.com/muidea/{{$projectName}}/core/service/{{$lowName}}"
	{{end}}
)

// New 新建Core
func New() (*Core, error) {
	var err error
	clnt := batisClient.NewClient(config.BatisService(), config.EndpointName())
	defer func() {
		if err != nil {
			clnt.Release()
		}
	}()

	err = model.InitializeModel(clnt)
	if err != nil {
		log.Printf("initializeModel failed, err:%s", err.Error())
		return nil, err
	}

	commonHandler := handler.NewCommonHandler(clnt, config.CasService())
	{{range $k, $v := $moduleList}}
	{{$lowName := $v.Lower $v.Name}}
	{{$lowName}}Handler := handler.New{{$v.Name}}Handler(clnt)
	{{end}}

	baseService := base.New(
		commonHandler,
		config.CasService(),
		config.FileService(),
	)

	{{range $k, $v := $moduleList}}
	{{$lowName := $v.Lower $v.Name}}
	{{$lowName}}Service := {{$lowName}}.New(commonHandler,{{$lowName}}Handler)
	{{end}}

	systemService := system.New(commonHandler)

	core := &Core{
		baseService:      baseService,
		{{range $k, $v := $moduleList}}
		{{$lowName := $v.Lower $v.Name}}
		{{$lowName}}Service: {{$lowName}}Service,
		{{end}}
		systemService:    systemService,
	}

	return core, nil
}

// Core Core对象
type Core struct {
	baseService      *base.Base
	{{range $k, $v := $moduleList}}
	{{$lowName := $v.Lower $v.Name}}
	{{$lowName}}Service *{{$lowName}}.{{$v.Name}}
	{{end}}
	systemService    *system.System
}

// Startup 启动
func (s *Core) Startup(router engine.Router) {
	sessionRegistry := session.CreateRegistry(nil)
	casRouteRegistry := casToolkit.NewCasRegistry(s.baseService, router)
	privateRouteRegistry := casToolkit.NewPrivateRegistry(s.baseService, router)

	s.baseService.BindRegistry(sessionRegistry, casRouteRegistry, privateRouteRegistry)
	{{range $k, $v := $moduleList}}
	{{$lowName := $v.Lower $v.Name}}
	s.{{$lowName}}Service.BindRegistry(sessionRegistry, casRouteRegistry, privateRouteRegistry)
	{{end}}
	s.systemService.BindRegistry(sessionRegistry, casRouteRegistry, privateRouteRegistry)

	s.baseService.RegisterHandler(router)
	{{range $k, $v := $moduleList}}
	{{$lowName := $v.Lower $v.Name}}
	s.{{$lowName}}Service.RegisterHandler(router)
	{{end}}
	s.systemService.RegisterHandler(router)

	casRouteRegistry.RegisterHandler(router)
	privateRouteRegistry.RegisterHandler(router)
}

// Teardown 销毁
func (s *Core) Teardown() {
}
`

	return
}

func (s *Core) Create(define *Define) (err error) {
	filePath := fmt.Sprintf("./core/%s.go", "core")
	func() {
		fileHandle, fileErr := os.Create(filePath)
		if fileErr != nil {
			err = fileErr
			return
		}
		defer fileHandle.Close()

		templateHandle, templateErr := template.New("core").Parse(s.getTemplate())
		if templateErr != nil {
			err = templateErr
			return
		}

		data := map[string]interface{}{"ProjectName": s.projectName, "Define": define}
		err = templateHandle.Execute(fileHandle, data)
	}()
	if err != nil {
		return
	}

	cmd := exec.Command("gofmt", "-w", filePath)
	err = cmd.Run()
	return
}
