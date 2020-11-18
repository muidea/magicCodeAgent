package service

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

type Define struct {
	projectName string
}

func New(projectName string) *Define {
	return &Define{projectName: projectName}
}

func (s *Define) getTemplate() (ret string) {
	ret = `
	{{$projectName := .ProjectName}}
	{{$packageName := .PackageName }}
	{{$module := .Module}}
	{{$moduleName := $module.Name}}
	{{$actionList := $module.Action}}
	{{$childList := $module.Child}}
	package {{$packageName}}
	
	import (
		"fmt"
		"net/http"
	
		commonCommon "github.com/muidea/magicCommon/common"
		commonSession "github.com/muidea/magicCommon/session"
	
		"github.com/muidea/{{$projectName}}/common"
		"github.com/muidea/{{$projectName}}/core/handler"
		"github.com/muidea/{{$projectName}}/core/service/{{$packageName}}/biz"
	
		casModel "github.com/muidea/magicCas/model"
		casToolkit "github.com/muidea/magicCas/toolkit"
	
		engine "github.com/muidea/magicEngine"
	)
	
	// {{$moduleName}} 路由信息
	type {{$moduleName}} struct {
		sessionRegistry      commonSession.Registry
		privateRouteRegistry casToolkit.PrivateRegistry
		casRouteRegistry     casToolkit.CasRegistry
		
		commonHandler handler.CommonHandler
	}
	
	// New create route
	func New(
		commonHandler handler.CommonHandler,
	) *{{$moduleName}} {
		route := &{{$moduleName}}{
			commonHandler: commonHandler,
		}
	
		return route
	}
	
	func (s *{{$moduleName}}) BindRegistry(
		sessionRegistry commonSession.Registry,
		casRouteRegistry casToolkit.CasRegistry,
		privateRouteRegistry casToolkit.PrivateRegistry) {
	
		s.sessionRegistry = sessionRegistry
		s.casRouteRegistry = casRouteRegistry
		s.privateRouteRegistry = privateRouteRegistry
	}
	
	// Handle middleware handler
	func (s *{{$moduleName}}) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
		curSession := s.sessionRegistry.GetSession(res, req)
	
		sessionInfo := curSession.GetSessionInfo()
		sessionInfo.Scope = commonCommon.ShareSession
	
		values := req.URL.Query()
		values = sessionInfo.Encode(values)
		req.URL.RawQuery = values.Encode()
	
		ctx.Next()
	}
	
	// RegisterHandler 注册路由
	func (s *{{$moduleName}}) RegisterHandler(router engine.Router) {

	}
	
	func (s *{{$moduleName}}) getCurrentEntity(res http.ResponseWriter, req *http.Request) (ret *casModel.Entity, err error) {
		curSession := s.sessionRegistry.GetSession(res, req)
		authVal, ok := curSession.GetOption(commonCommon.AuthAccount)
		if !ok {
			err = fmt.Errorf("无效权限,未通过验证")
			return
		}
		ret = authVal.(*casModel.Entity)
		return
	}

	func (s *Bill) getCurrentNamespace(res http.ResponseWriter, req *http.Request) (ret string, err error) {
		return
	}
	`

	return
}

func (s *Define) Create(module *model.Module) (err error) {
	if len(module.Action) == 0 && len(module.Child) == 0 {
		return
	}

	filePath := fmt.Sprintf("./core/service/%s", strings.ToLower(module.Name))
	existFlag, existErr := util.Exists(filePath)
	if existErr != nil {
		err = existErr
		return
	}
	if !existFlag {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return
		}
	}

	filePath = fmt.Sprintf("%s/%s.go", filePath, "service")
	existsFlag, existsErr := util.Exists(filePath)
	if existsErr != nil {
		err = fmt.Errorf("unexpect, fileName:%s", filePath)
		return
	}
	if !existsFlag {
		func() {
			fileHandle, fileErr := os.Create(filePath)
			if fileErr != nil {
				err = fileErr
				return
			}
			defer fileHandle.Close()

			templateHandle, templateErr := template.New(module.Name).Parse(s.getTemplate())
			if templateErr != nil {
				err = templateErr
				return
			}

			data := map[string]interface{}{"ProjectName": s.projectName, "PackageName": module.LowerName(), "Module": module}
			err = templateHandle.Execute(fileHandle, data)
		}()
		if err != nil {
			return
		}

		cmd := exec.Command("gofmt", "-w", filePath)
		err = cmd.Run()
		if err != nil {
			return
		}
	}

	instance := &Instance{projectName: s.projectName, packageName: module.LowerName()}
	err = instance.Create(module)
	if err != nil {
		return
	}
	for _, val := range module.Child {
		instance := &Instance{projectName: s.projectName, packageName: module.LowerName()}
		err = instance.Create(val)
		if err != nil {
			return
		}
	}

	return
}
