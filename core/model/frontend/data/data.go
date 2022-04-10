package data

import (
	"fmt"
	"os"
	"path"

	log "github.com/cihub/seelog"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

type Data struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Data {
	return &Data{rootPath: rootPath, projectName: projectName}
}

func (s *Data) declareObject(ptr *parser.Package) (ret string) {
	objectDeclare := `{{$ModuleVal := .ModuleVal}}{{with .ModuleVal}}{{range .Objects}}
{{declareObjectView .}}

export const formatter{{.Name}} = (val: any) => {
  return '-'
}{{end}}{{end}}`

	factory.RegisterFunc("declareObjectView", DeclareObjectView)
	defer factory.UnRegisterFunc("declareObjectView")

	factoryPtr := factory.NewMemoryFactor(objectDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Data) Create(ptr *parser.Package) (err error) {
	if len(ptr.Objects) == 0 {
		return
	}

	moduleName := mu.FormatModulePath(ptr)
	commonFile := fmt.Sprintf("%s.ts", moduleName)
	commonFilePath := path.Join(s.rootPath, "data", commonFile)

	log.Info(commonFilePath)
	func() {
		fileHandle, fileErr := os.OpenFile(commonFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if fileErr != nil {
			err = fileErr
			return
		}
		defer fileHandle.Close()

		_, err = fileHandle.WriteString(s.declareObject(ptr))
		if err != nil {
			return
		}
	}()

	return
}
