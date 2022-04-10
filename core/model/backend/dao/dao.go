package dao

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	log "github.com/cihub/seelog"

	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

type Dao struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Dao {
	return &Dao{rootPath: rootPath, projectName: projectName}
}

func (s *Dao) declarePackage() (ret string) {
	packageDeclare := `
package dao
`
	byteBuff := bytes.Buffer{}
	byteBuff.WriteString(packageDeclare)
	ret = byteBuff.String()
	return
}

func (s *Dao) declareImport() (ret string) {
	importDeclare := `
import (
	"fmt"

	"github.com/muidea/magicBatis/client"
	bc "github.com/muidea/magicBatis/common"
	"github.com/muidea/{{.ProjectName}}/model"
)
`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Dao) declareInterface(ptr *parser.Package) (ret string) {
	interfaceDeclare := `
type {{formatModuleName .ModuleVal}} interface { {{with .ModuleVal}}
{{range .Objects}}
	Filter{{formatObjectName .}}(filter *bc.QueryFilter, namespace string) (ret []*model.{{formatObjectName .}}, total int64, err error)
	Query{{formatObjectName .}}(id int, namespace string) (ret *model.{{formatObjectName .}}, err error)
	Create{{formatObjectName .}}(ptr *model.{{formatObjectName .}}, namespace string) (ret *model.{{formatObjectName .}}, err error)
	Update{{formatObjectName .}}(ptr *model.{{formatObjectName .}}, namespace string) (ret *model.{{formatObjectName .}}, err error)
	Delete{{formatObjectName .}}(id int, namespace string) (ret *model.{{formatObjectName .}}, err error)
{{end}}{{end}}
}`
	factoryPtr := factory.NewMemoryFactor(interfaceDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Dao) declareDefine(ptr *parser.Package) (ret string) {
	defineDeclare := `
func New(clnt client.Client) {{formatModuleName .ModuleVal}} {
	return &{{formatModuleVal .ModuleVal}}{batisClient: clnt}
}

type {{formatModuleVal .ModuleVal}} struct {
	batisClient client.Client
}`
	type Item struct {
		NameVal string
	}

	factoryPtr := factory.NewMemoryFactor(defineDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Dao) declareFunc(ptr *parser.Package) (ret string) {
	funcDeclare := `{{$ModuleVal := .ModuleVal}}{{with .ModuleVal}}
{{range .Objects}}
func (s *{{formatModuleVal $ModuleVal}}) Filter{{formatObjectName .}}(filter *bc.QueryFilter, namespace string) (ret []*model.{{formatObjectName .}}, total int64, err error) {
	ret = []*model.{{formatObjectName .}}{}
	if filter == nil {
		filter = bc.NewFilter()
	}
	filter.Equal("Namespace", namespace)
	filter.ValueMask(&{{defineObject .}})
	total, err = s.batisClient.BatchQueryEntity(&ret, filter)
	return
}

func (s *{{formatModuleVal $ModuleVal}}) Query{{formatObjectName .}}(id int, namespace string) (ret *model.{{formatObjectName .}}, err error) {
	if id <= 0 {
		err = fmt.Errorf("illegal id value, id:%d", id)
		return
	}
	ptr := &{{defineObject .}}
	ptr.ID = id
	ptr.Namespace = namespace
	err = s.batisClient.QueryEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *{{formatModuleVal $ModuleVal}}) Create{{formatObjectName .}}(ptr *model.{{formatObjectName .}}, namespace string) (ret *model.{{formatObjectName .}}, err error) {
	if ptr == nil {
		err = fmt.Errorf("illegal ptr value")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.InsertEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *{{formatModuleVal $ModuleVal}}) Update{{formatObjectName .}}(ptr *model.{{formatObjectName .}}, namespace string) (ret *model.{{formatObjectName .}}, err error) {
	if ptr == nil || ptr.ID <= 0 {
		err = fmt.Errorf("illegal update ptr")
		return
	}
	ptr.Namespace = namespace
	err = s.batisClient.UpdateEntity(ptr)
	if err != nil {
		return
	}
	ret = ptr
	return
}

func (s *{{formatModuleVal $ModuleVal}}) Delete{{formatObjectName .}}(id int, namespace string) (ret *model.{{formatObjectName .}}, err error) {
	valPtr, valErr := s.Query{{formatObjectName .}}(id, namespace)
	if valErr != nil {
		err = valErr
		return
	}
	err = s.batisClient.DeleteEntity(valPtr)
	if err != nil {
		return
	}
	ret = valPtr
	return
}
{{end}}{{end}}`

	factoryPtr := factory.NewMemoryFactor(funcDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Dao) Create(ptr *parser.Package) (err error) {
	if len(ptr.Objects) == 0 {
		return
	}

	commonFile := fmt.Sprintf("%s.go", "dao")
	commonFilePath := path.Join(s.rootPath, "core/module", mu.FormatModulePath(ptr), "dao", commonFile)

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

		_, err = fileHandle.WriteString(s.declareImport())
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareInterface(ptr))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareDefine(ptr))
		if err != nil {
			return
		}

		_, err = fileHandle.WriteString(s.declareFunc(ptr))
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
