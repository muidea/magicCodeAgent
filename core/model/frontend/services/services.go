package services

import (
	"fmt"
	"os"
	"path"

	log "github.com/cihub/seelog"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

type Services struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Services {
	return &Services{rootPath: rootPath, projectName: projectName}
}

func (s *Services) declareImport(ptr *parser.Package) (ret string) {
	moduleDeclare := `
import {genApi} from '@/utils';
`
	factoryPtr := factory.NewMemoryFactor(moduleDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal

	return
}

func (s *Services) declareObject(ptr *parser.Package) (ret string) {
	objectDeclare := `{{$ModuleVal := .ModuleVal}}{{with .ModuleVal}}
const index :{[index:string]:string} = {
{{- range .Objects}}
  filter{{formatObjectActionName $ModuleVal .}}: 'GET /{{formatObjectActionVal $ModuleVal .}}/query/',
  query{{formatObjectActionName $ModuleVal .}}: 'GET /{{formatObjectActionVal $ModuleVal .}}/query/:id',
  create{{formatObjectActionName $ModuleVal .}}: 'POST /{{formatObjectActionVal $ModuleVal .}}/create/',
  update{{formatObjectActionName $ModuleVal .}}: 'PUT /{{formatObjectActionVal $ModuleVal .}}/update/:id',
  delete{{formatObjectActionName $ModuleVal .}}: 'DELETE /{{formatObjectActionVal $ModuleVal .}}/delete/:id',
{{end}}
}{{end}}`
	factoryPtr := factory.NewMemoryFactor(objectDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Services) declareAPIFunction(ptr *parser.Package) (ret string) {
	moduleDeclare := `
const APIFunction:{[index:string]:any} = {}
Object.keys(index).forEach((key)=>{
  APIFunction[key] = genApi(index[key])
})

export default APIFunction
`
	factoryPtr := factory.NewMemoryFactor(moduleDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal

	return
}

func (s *Services) Create(ptr *parser.Package) (err error) {
	if len(ptr.Objects) == 0 {
		return
	}

	moduleName := mu.FormatModulePath(ptr)
	commonFile := fmt.Sprintf("/%s/services/index.ts", moduleName)
	commonFilePath := path.Join(s.rootPath, "pages", commonFile)
	log.Info(commonFilePath)
	func() {
		fileHandle, fileErr := os.OpenFile(commonFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if fileErr != nil {
			err = fileErr
			return
		}
		defer fileHandle.Close()

		_, err = fileHandle.WriteString(s.declareImport(ptr))
		if err != nil {
			return
		}
		_, err = fileHandle.WriteString(s.declareObject(ptr))
		if err != nil {
			return
		}
		_, err = fileHandle.WriteString(s.declareAPIFunction(ptr))
		if err != nil {
			return
		}
	}()

	return
}
