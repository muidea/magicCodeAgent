package biz

import (
	"bytes"
	"fmt"
	log "github.com/cihub/seelog"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
	"os"
	"os/exec"
	"path"
)

type Biz struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Biz {
	return &Biz{rootPath: rootPath, projectName: projectName}
}

func (s *Biz) declarePackage() (ret string) {
	packageDeclare := `
package biz
`
	byteBuff := bytes.Buffer{}
	byteBuff.WriteString(packageDeclare)
	ret = byteBuff.String()
	return
}

func (s *Biz) declareBizImport(ptr *parser.Package) (ret string) {
	importDeclare := `
import (
	"github.com/muidea/magicBatis/client"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/{{.ProjectName}}/common"
	"github.com/muidea/{{.ProjectName}}/core/base/biz"
	"github.com/muidea/{{.ProjectName}}/core/module/{{formatModuleVal .ModuleVal}}/dao"
)
`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Biz) declareObjectImport(ptr *parser.Package) (ret string) {
	importDeclare := `
import (
	"time"

	bc "github.com/muidea/magicBatis/common"

	"github.com/muidea/magicCommon/event"
	fn "github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/{{.ProjectName}}/common"
	"github.com/muidea/{{.ProjectName}}/model"
)
`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Biz) declareDefine(ptr *parser.Package) (ret string) {
	defineDeclare := `
type {{formatModuleName .ModuleVal}} struct {
	biz.Base

	{{formatModuleVal .ModuleVal}}Dao dao.{{formatModuleName .ModuleVal}}

	endpointName string
}

func New(
	endpointName string,
	batisClient client.Client,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) *{{formatModuleName .ModuleVal}} {

	ptr := &{{formatModuleName .ModuleVal}}{
		Base:         biz.New(common.{{formatModuleName .ModuleVal}}Module, eventHub, backgroundRoutine),
		{{formatModuleVal .ModuleVal}}Dao:   dao.New(batisClient),
		endpointName: endpointName,
	}

	return ptr
}

func (s *{{formatModuleName .ModuleVal}}) Notify(event event.Event, result event.Result) {
	return
}`
	factoryPtr := factory.NewMemoryFactor(defineDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Biz) declareObjectFunc(modulePtr *parser.Package, objectPtr *parser.Object) (ret string) {
	funcDeclare := `{{$ModuleVal := .ModuleVal}}
{{range .Items}}
func (s *{{formatModuleName $ModuleVal}}) Filter{{formatObjectName .}}(filter *bc.QueryFilter, namespace string) (ret []*model.{{formatObjectName .}}, total int64, err error) {
	ret, total, err = s.{{formatModuleVal $ModuleVal}}Dao.Filter{{formatObjectName .}}(filter, namespace)
	return
}

func (s *{{formatModuleName $ModuleVal}}) Query{{formatObjectName .}}(id int, namespace string) (ret *model.{{formatObjectName .}}, err error) {
	ret, err = s.{{formatModuleVal $ModuleVal}}Dao.Query{{formatObjectName .}}(id, namespace)
	return
}

func (s *{{formatModuleName $ModuleVal}}) Create{{formatObjectName .}}(ptr *common.{{formatObjectName .}}Param, creater int, namespace string) (ret *model.{{formatObjectName .}}, err error) {
	{{formatObjectVal .}}Ptr := ptr.To{{formatObjectName .}}(nil)
	{{formatObjectVal .}}Ptr.Creater = creater
	{{formatObjectVal .}}Ptr.UpdateTime = time.Now().UTC().Unix()
	ret, err = s.{{formatModuleVal $ModuleVal}}Dao.Create{{formatObjectName .}}({{formatObjectVal .}}Ptr, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.Notify{{formatObjectActionName $ModuleVal .}}, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Create)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *{{formatModuleName $ModuleVal}}) Update{{formatObjectName .}}(id int, ptr *common.{{formatObjectName .}}Param, updater int, namespace string) (ret *model.{{formatObjectName .}}, err error) {
	current{{formatObjectName .}}, currentErr := s.{{formatModuleVal $ModuleVal}}Dao.Query{{formatObjectName .}}(id, namespace)
	if currentErr != nil {
		return
	}

	current{{formatObjectName .}} = ptr.To{{formatObjectName .}}(current{{formatObjectName .}})
	current{{formatObjectName .}}.Creater = updater
	current{{formatObjectName .}}.UpdateTime = time.Now().UTC().Unix()
	ret, err = s.{{formatModuleVal $ModuleVal}}Dao.Update{{formatObjectName .}}(current{{formatObjectName .}}, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.Notify{{formatObjectActionName $ModuleVal .}}, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Update)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *{{formatModuleName $ModuleVal}}) Delete{{formatObjectName .}}(id int, namespace string) (ret *model.{{formatObjectName .}}, err error) {
	ret, err = s.{{formatModuleVal $ModuleVal}}Dao.Delete{{formatObjectName .}}(id, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.Notify{{formatObjectActionName $ModuleVal .}}, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Delete)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}
{{end}}`
	factoryPtr := factory.NewMemoryFactor(funcDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": modulePtr, "Items": []*parser.Object{objectPtr}})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Biz) Create(ptr *parser.Package) (err error) {
	if len(ptr.Objects) == 0 {
		return
	}

	commonFile := fmt.Sprintf("%s.go", "biz")
	commonFilePath := path.Join(s.rootPath, "core/module", util.ToLower(mu.FormatModuleName(ptr)), "biz", commonFile)

	log.Info(commonFilePath)
	func() {
		fileHandle, fileErr := os.OpenFile(commonFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if fileErr != nil {
			err = fileErr
			return
		}
		defer fileHandle.Close()

		_, err = fileHandle.WriteString(s.declarePackage())
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareBizImport(ptr))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareDefine(ptr))
		if err != nil {
			return
		}
	}()

	if err != nil {
		return
	}

	cmd := exec.Command("gofmt", "-w", commonFilePath)
	err = cmd.Run()

	for _, val := range ptr.Objects {
		err = s.createObject(ptr, val)
		if err != nil {
			return
		}
	}

	return
}

func (s *Biz) createObject(modulePtr *parser.Package, ptr *parser.Object) (err error) {
	commonFile := fmt.Sprintf("%s.go", util.ToLower(ptr.Name))
	commonFilePath := path.Join(s.rootPath, "core/module", util.ToLower(mu.FormatModuleName(modulePtr)), "biz", commonFile)

	log.Info(commonFilePath)
	func() {
		fileHandle, fileErr := os.OpenFile(commonFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if fileErr != nil {
			err = fileErr
			return
		}
		defer fileHandle.Close()

		_, err = fileHandle.WriteString(s.declarePackage())
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareObjectImport(modulePtr))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareObjectFunc(modulePtr, ptr))
		if err != nil {
			return
		}
	}()

	if err != nil {
		return
	}

	cmd := exec.Command("gofmt", "-w", commonFilePath)
	err = cmd.Run()
	return
}
