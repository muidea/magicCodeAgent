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

type Service struct {
	projectName string
}

func New(projectName string) *Service {
	return &Service{projectName: projectName}
}

func (s *Service) getTemplate() (ret string) {
	ret = `
{{$projectName := .ProjectName}}
{{$serviceName := .ServiceName}}
{{$packageName := .PackageName }}
{{$module := .Module}}
{{$moduleName := $module.Name}}
{{$actionList := $module.Action}}
{{$childList := $module.Child}}
{{$lowName := $module.Lower $moduleName}}
package {{$packageName}}

import (
	"fmt"
	"net/http"

	commonCommon "github.com/muidea/magicCommon/common"
	commonSession "github.com/muidea/magicCommon/session"

	casModel "github.com/muidea/magicCas/model"
	casToolkit "github.com/muidea/magicCas/toolkit"

	engine "github.com/muidea/magicEngine"

	"github.com/muidea/{{$projectName}}/common"
	"github.com/muidea/{{$projectName}}/core/handler"
	"github.com/muidea/{{$projectName}}/core/biz/{{$lowName}}"
)

// {{$serviceName}} 路由信息
type {{$serviceName}} struct {
	sessionRegistry      commonSession.Registry
	privateRouteRegistry casToolkit.PrivateRegistry
	casRouteRegistry     casToolkit.CasRegistry

	{{$lowName}}Biz *{{$lowName}}.{{$serviceName}}
	commonHandler handler.CommonHandler
}

// New create route
func New(
	commonHandler handler.CommonHandler,
	{{$lowName}}Biz *{{$lowName}}.{{$serviceName}},
) *{{$serviceName}} {
	route := &{{$serviceName}}{
		commonHandler: commonHandler,
		{{$lowName}}Biz: {{$lowName}}Biz,
	}

	return route
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
{{range $k, $v := $actionList}}
{{$url := $module.UrlRoute $serviceName $v }}
{{$url}}
{{end}}
{{range $mk, $mv := $childList}}
{{$subModuleName := $mv.Name}}
{{$subModuleAction := $mv.Action}}
{{range $ak, $av := $subModuleAction}}
{{$url := $mv.UrlRoute $serviceName $av }}
{{$url}}
{{end}}
{{end}}
}

func (s *{{$serviceName}}) getCurrentEntity(res http.ResponseWriter, req *http.Request) (ret *casModel.Entity, err error) {
	curSession := s.sessionRegistry.GetSession(res, req)
	authVal, ok := curSession.GetOption(commonCommon.AuthAccount)
	if !ok {
		err = fmt.Errorf("无效权限,未通过验证")
		return
	}
	ret = authVal.(*casModel.Entity)
	return
}

func (s *{{$serviceName}}) queryEntity(id int, sessionInfo *commonCommon.SessionInfo, namespace string) (ret *casModel.Entity) {
	entityPtr, entityErr := s.commonHandler.QueryEntity(id, sessionInfo, namespace)
	if entityErr != nil {
		ret = casModel.UnknownEntity()
		return
	}

	ret = entityPtr
	return
}

func (s *{{$serviceName}}) getCurrentNamespace(res http.ResponseWriter, req *http.Request) (ret string) {
	namespace := req.Header.Get(casCommon.NamespaceID)
	if namespace != "" {
		ret = namespace
		return
	}

	items := strings.Split(req.Host, ":")
	if nil != cn.ParseIP(items[0]) {
		ret = "default"
		return
	}

	items = strings.Split(items[0], ".")
	if len(items) <= 1 {
		ret = "default"
		return
	}

	ret = items[0]
	return
}
`
	return
}

func (s *Service) Create(module *model.Module, forceFlag bool) (err error) {
	if len(module.Action) == 0 && len(module.Child) == 0 {
		return
	}

	filePath := fmt.Sprintf("./core/service/%s", strings.ToLower(module.Name))
	existFlag := util.Exists(filePath)
	if !existFlag {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return
		}
	}

	filePath = fmt.Sprintf("%s/%s.go", filePath, "service")
	existsFlag := util.Exists(filePath)
	if !existsFlag || forceFlag {
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

			data := map[string]interface{}{"ProjectName": s.projectName, "ServiceName": module.Name, "PackageName": module.Lower(module.Name), "Module": module}
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

	if !module.IsPackage() {
		instance := &Instance{projectName: s.projectName, serviceName: module.Name, packageName: module.Lower(module.Name)}
		err = instance.Create(module, forceFlag)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	for _, val := range module.Child {
		if val.IsPackage() {
			continue
		}

		instance := &Instance{projectName: s.projectName, serviceName: module.Name, packageName: module.Lower(module.Name)}
		err = instance.Create(val, forceFlag)
		if err != nil {
			return
		}
	}

	return
}
