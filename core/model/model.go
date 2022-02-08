package model

import (
	"io/ioutil"
	"path"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicToolkit/core/model/biz"
	"github.com/muidea/magicToolkit/core/model/common"
	"github.com/muidea/magicToolkit/core/model/dao"
	"github.com/muidea/magicToolkit/core/model/service"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/core/module"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/model"
	"github.com/muidea/magicToolkit/parser"
)

type Model struct {
	rootPath    string
	projectName string
}

func New(rootPath string) *Model {
	return &Model{rootPath: rootPath}
}

func (s *Model) Create(projectName string) (err error) {
	factory.RegisterFunc("isSimple", mu.IsSimple)
	factory.RegisterFunc("defineItem", mu.DefineItem)
	factory.RegisterFunc("defineObject", mu.DefineObject)
	defer factory.UnRegisterFunc("defineObject")
	defer factory.UnRegisterFunc("defineItem")
	defer factory.UnRegisterFunc("isSimple")

	s.projectName = projectName

	modelPath := path.Join(s.rootPath, "model")
	existFlag := util.Exists(modelPath)
	if !existFlag {
		return
	}

	fdList, fdErr := ioutil.ReadDir(modelPath)
	if fdErr != nil {
		err = fdErr
		return
	}

	for _, fd := range fdList {
		fileName := fd.Name()
		fullPath := path.Join(modelPath, fileName)
		if !util.Exists(fullPath) || fileName == "log.go" || fileName == "totalizer.go" || fileName == "helper.go" {
			continue
		}

		packagePtr, packageErr := parser.Parse(fullPath)
		if packageErr != nil {
			log.Errorf("parse failed, file:%s, err:%s", fullPath, packageErr)
		}

		_, moduleErr := s.createModule(packagePtr)
		if moduleErr != nil {
			log.Errorf("create module failed, name:%s, err:%s", packagePtr.Name, moduleErr.Error())
			return
		}

		commonPtr := common.New(s.rootPath, s.projectName)
		commonErr := commonPtr.Create(packagePtr)
		if commonErr != nil {
			log.Errorf("create common file failed, name:%s, err:%s", packagePtr.Name, commonErr.Error())
			return
		}

		daoPtr := dao.New(s.rootPath, s.projectName)
		daoErr := daoPtr.Create(packagePtr)
		if daoErr != nil {
			log.Errorf("create dao file failed, name:%s, err:%s", packagePtr.Name, daoErr.Error())
			return
		}

		bizPtr := biz.New(s.rootPath, s.projectName)
		bizErr := bizPtr.Create(packagePtr)
		if bizErr != nil {
			log.Errorf("create biz file failed, name:%s, err:%s", packagePtr.Name, bizErr.Error())
			return
		}

		servicePtr := service.New(s.rootPath, s.projectName)
		serviceErr := servicePtr.Create(packagePtr)
		if serviceErr != nil {
			log.Errorf("create service file failed, name:%s, err:%s", packagePtr.Name, bizErr.Error())
			return
		}
	}

	return
}

func (s *Model) createModule(ptr *parser.Package) (ret *model.Module, err error) {
	modulePtr := module.New(s.rootPath)
	ret, err = modulePtr.Create(s.projectName, util.ToTitle(ptr.FileName))
	return
}
