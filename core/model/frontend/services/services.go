package services

import (
	"fmt"
	"os"
	"path"

	log "github.com/cihub/seelog"

	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/core/util"
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

func (s *Services) declareImport(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	moduleDeclare := `import {genApi} from '@/utils';
`
	factoryPtr := factory.NewMemoryFactor(moduleDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal

	return
}

func (s *Services) declareObject(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	objectDeclare := `
const api :{[index:string]:string} = {
  filter{{formatObjectName .ObjectVal}}: 'GET /{{formatObjectActionValue .ModuleVal .ObjectVal}}/query/',
  query{{formatObjectName .ObjectVal}}: 'GET /{{formatObjectActionValue .ModuleVal .ObjectVal}}/query/:id',
  create{{formatObjectName .ObjectVal}}: 'POST /{{formatObjectActionValue .ModuleVal .ObjectVal}}/create/',
  update{{formatObjectName .ObjectVal}}: 'PUT /{{formatObjectActionValue .ModuleVal .ObjectVal}}/update/:id',
  delete{{formatObjectName .ObjectVal}}: 'DELETE /{{formatObjectActionValue .ModuleVal .ObjectVal}}/delete/:id',
}
`
	factoryPtr := factory.NewMemoryFactor(objectDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Services) declareAPIFunction(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	moduleDeclare := `
const APIFunction:{[index:string]:any} = {}
Object.keys(api).forEach((key)=>{
  APIFunction[key] = genApi(api[key])
})

export default APIFunction
`
	factoryPtr := factory.NewMemoryFactor(moduleDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
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

	for _, val := range ptr.Objects {
		commonPath := fmt.Sprintf("%s/pages/%s/services", s.rootPath, mu.FormatObjectActionValue(ptr, val))
		if !util.Exists(commonPath) {
			err = os.MkdirAll(commonPath, os.ModePerm)
			if err != nil {
				return
			}
		}
		commonFile := path.Join(commonPath, "index.ts")
		log.Info(commonFile)
		func() {
			fileHandle, fileErr := os.OpenFile(commonFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if fileErr != nil {
				err = fileErr
				return
			}
			defer fileHandle.Close()

			_, err = fileHandle.WriteString(s.declareImport(ptr, val))
			if err != nil {
				return
			}
			_, err = fileHandle.WriteString(s.declareObject(ptr, val))
			if err != nil {
				return
			}
			_, err = fileHandle.WriteString(s.declareAPIFunction(ptr, val))
			if err != nil {
				return
			}
		}()
	}

	return
}
