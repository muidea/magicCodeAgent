package views

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

type Views struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Views {
	return &Views{rootPath: rootPath, projectName: projectName}
}

func (s *Views) declareObject(ptr *parser.Package) (ret string) {
	objectDeclare := `{{$ModuleVal := .ModuleVal}}{{with .ModuleVal}}
{{range .Objects}}
export interface {{.Name}} {
}

export const formatter{{.Name}} = (val: any) => {
  return '-'
}
{{end}}{{end}}`
	factoryPtr := factory.NewMemoryFactor(objectDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Views) Create(ptr *parser.Package) (err error) {
	if len(ptr.Objects) == 0 {
		return
	}

	for _, val := range ptr.Objects {
		commonPath := fmt.Sprintf("%s/pages/%s", s.rootPath, mu.FormatObjectActionValue(ptr, val))
		if !util.Exists(commonPath) {
			err = os.MkdirAll(commonPath, os.ModePerm)
			if err != nil {
				return
			}
		}
		commonFile := path.Join(commonPath, "index.tsx")
		log.Info(commonFile)

		func() {
			fileHandle, fileErr := os.OpenFile(commonFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
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
	}
	return
}
