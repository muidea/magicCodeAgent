/**

Package -> Service -> Model -> Action

InterfaceFunc

DeclareFunc

ImplementFunc

*/

package content

import (
	"encoding/json"
	"github.com/muidea/magicOrm/provider/remote"
	"github.com/muidea/magicToolkit/core/content/common"
	"github.com/muidea/magicToolkit/core/content/core"
	"github.com/muidea/magicToolkit/core/content/core/biz"
	"github.com/muidea/magicToolkit/core/content/core/service"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
	"io/ioutil"
	"os"
	"path"
)

// Content core define
type Content struct {
	projectPath string
	name2model  map[string]*remote.Object
}

// New new core
func New(projectPath string, modelList []*remote.Object) *Content {
	name2model := map[string]*remote.Object{}
	for _, val := range modelList {
		name := util.GetModelName(val)
		name2model[name] = val
	}

	return &Content{projectPath: projectPath, name2model: name2model}
}

// Create create
func (s *Content) Create(projectName string, forceFlag bool) (err error) {
	definePath := path.Join(s.projectPath, "doc", "define.json")
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

	contentDefine := &model.Define{}
	err = json.Unmarshal(byteContent, contentDefine)
	if err != nil {
		return
	}

	for _, val := range contentDefine.Module {
		common := common.New(contentDefine.ApiVersion, projectName)
		err = common.Create(val, forceFlag)
		if err != nil {
			return
		}

		biz := biz.New(projectName)
		err = biz.Create(val, forceFlag)
		if err != nil {
			return
		}

		service := service.New(projectName)
		err = service.Create(val, forceFlag)
		if err != nil {
			return
		}
	}

	core := core.New(projectName)
	err = core.Create(contentDefine, forceFlag)
	if err != nil {
		return
	}

	return
}
