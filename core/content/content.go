package content

import (
	"encoding/json"
	"github.com/muidea/magicToolkit/core/content/common"
	"github.com/muidea/magicToolkit/core/content/service"
	"github.com/muidea/magicToolkit/model"
	"io/ioutil"
	"os"
	"path"
)

// Content core define
type Content struct {
	rootPath string
}

// New new core
func New(rootPath string) *Content {
	return &Content{rootPath: rootPath}
}

// Create create
func (s *Content) Create(projectName string) (err error) {
	definePath := path.Join(s.rootPath, "doc", "define.json")
	defineFile, defineErr := os.Open(definePath)
	if defineErr != nil {
		err = defineErr
		return
	}
	defer defineFile.Close()

	byteContent, byteErr := ioutil.ReadAll(defineFile)
	if byteErr != nil {
		err = byteErr
		return
	}

	type Define struct {
		Module []*model.Module `json:"module"`
	}

	contentDefine := &Define{}
	err = json.Unmarshal(byteContent, contentDefine)
	if err != nil {
		return
	}

	for _, val := range contentDefine.Module {
		common := common.New(projectName)
		err = common.Create(val)
		if err != nil {
			return
		}

		service := service.New(projectName)
		err = service.Create(val)
		if err != nil {
			return
		}
	}

	return
}
