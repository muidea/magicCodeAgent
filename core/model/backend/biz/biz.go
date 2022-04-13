package biz

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	log "github.com/cihub/seelog"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
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
	"github.com/muidea/{{.ProjectName}}/core/module/{{formatModulePath .ModuleVal}}/dao"
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

	{{formatModuleValue .ModuleVal}}Dao dao.{{formatModuleName .ModuleVal}}

	endpointName string
}

func New(
	endpointName string,
	batisClient client.Client,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) *{{formatModuleName .ModuleVal}} {

	ptr := &{{formatModuleName .ModuleVal}}{
		Base:         biz.New(common.{{formatModuleName .ModuleVal}}Module, eventHub, backgroundRoutine),
		{{formatModuleValue .ModuleVal}}Dao:   dao.New(batisClient),
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
	funcDeclare := `{{$ModuleVal := .ModuleVal}}{{$ObjectVal := .ObjectVal}}
func (s *{{formatModuleName $ModuleVal}}) Filter{{formatObjectName $ObjectVal}}(filter *bc.QueryFilter, namespace string) (ret []*model.{{formatObjectName $ObjectVal}}, total int64, err error) {
	ret, total, err = s.{{formatModuleValue $ModuleVal}}Dao.Filter{{formatObjectName $ObjectVal}}(filter, namespace)
	return
}

func (s *{{formatModuleName $ModuleVal}}) Query{{formatObjectName $ObjectVal}}(id int, namespace string) (ret *model.{{formatObjectName $ObjectVal}}, err error) {
	ret, err = s.{{formatModuleValue $ModuleVal}}Dao.Query{{formatObjectName $ObjectVal}}(id, namespace)
	return
}

func (s *{{formatModuleName $ModuleVal}}) Create{{formatObjectName $ObjectVal}}(ptr *common.{{formatObjectName $ObjectVal}}Param, creater int, namespace string) (ret *model.{{formatObjectName $ObjectVal}}, err error) {
	{{formatObjectValue $ObjectVal}}Ptr := ptr.To{{formatObjectName $ObjectVal}}(nil)
{{if .HasEntity -}}
	{{formatObjectValue $ObjectVal}}Ptr.Creater = creater
	{{formatObjectValue $ObjectVal}}Ptr.CreateTime = time.Now().UTC().Unix()
{{- end}}
	ret, err = s.{{formatModuleValue $ModuleVal}}Dao.Create{{formatObjectName $ObjectVal}}({{formatObjectValue $ObjectVal}}Ptr, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.Notify{{formatObjectActionName $ModuleVal $ObjectVal}}, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Create)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *{{formatModuleName $ModuleVal}}) Update{{formatObjectName $ObjectVal}}(id int, ptr *common.{{formatObjectName $ObjectVal}}Param, updater int, namespace string) (ret *model.{{formatObjectName $ObjectVal}}, err error) {
	current{{formatObjectName $ObjectVal}}, currentErr := s.{{formatModuleValue $ModuleVal}}Dao.Query{{formatObjectName $ObjectVal}}(id, namespace)
	if currentErr != nil {
		err = currentErr
		return
	}

	current{{formatObjectName $ObjectVal}} = ptr.To{{formatObjectName $ObjectVal}}(current{{formatObjectName $ObjectVal}})
{{if .HasEntity -}}
	current{{formatObjectName $ObjectVal}}.Creater = updater
	current{{formatObjectName $ObjectVal}}.CreateTime = time.Now().UTC().Unix()
{{- end}}
	ret, err = s.{{formatModuleValue $ModuleVal}}Dao.Update{{formatObjectName $ObjectVal}}(current{{formatObjectName $ObjectVal}}, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.Notify{{formatObjectActionName $ModuleVal $ObjectVal}}, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Update)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}

func (s *{{formatModuleName $ModuleVal}}) Delete{{formatObjectName $ObjectVal}}(id int, namespace string) (ret *model.{{formatObjectName $ObjectVal}}, err error) {
	ret, err = s.{{formatModuleValue $ModuleVal}}Dao.Delete{{formatObjectName $ObjectVal}}(id, namespace)
	if err != nil {
		return
	}

	eid := fn.FormatID(common.Notify{{formatObjectActionName $ModuleVal $ObjectVal}}, ret.ID)
	header := event.NewValues()
	header.Set("action", common.Delete)
	header.Set("namespace", namespace)
	s.BroadCast(eid, header, ret)
	return
}
`
	factoryPtr := factory.NewMemoryFactor(funcDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": modulePtr, "ObjectVal": objectPtr, "HasEntity": mu.HasEntity(objectPtr)})
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
	commonFilePath := path.Join(s.rootPath, "core/module", mu.FormatModulePath(ptr), "biz", commonFile)

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

func (s *Biz) createObject(modulePtr *parser.Package, objectPtr *parser.Object) (err error) {
	commonFile := fmt.Sprintf("%s.go", util.ToLower(objectPtr.Name))
	commonFilePath := path.Join(s.rootPath, "core/module", mu.FormatModulePath(modulePtr), "biz", commonFile)

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

		_, err = fileHandle.WriteString(s.declareObjectFunc(modulePtr, objectPtr))
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
