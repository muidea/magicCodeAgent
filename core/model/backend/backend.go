package backend

import (
	log "github.com/cihub/seelog"

	"github.com/muidea/magicToolkit/core/model/backend/biz"
	"github.com/muidea/magicToolkit/core/model/backend/common"
	"github.com/muidea/magicToolkit/core/model/backend/dao"
	"github.com/muidea/magicToolkit/core/model/backend/service"
	bu "github.com/muidea/magicToolkit/core/model/backend/util"
	"github.com/muidea/magicToolkit/core/module"
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/model"
	"github.com/muidea/magicToolkit/parser"
)

type Backend struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Backend {
	return &Backend{rootPath: rootPath, projectName: projectName}
}

func (s *Backend) createModule(ptr *parser.Package) (ret *model.Module, err error) {
	modulePtr := module.New(s.rootPath)
	ret, err = modulePtr.Create(s.projectName, util.ToTitle(ptr.FileName))
	return
}

func (s *Backend) Create(packagePtr *parser.Package) (err error) {
	factory.RegisterFunc("defineItem", bu.DefineItem)
	factory.RegisterFunc("defineObject", bu.DefineObject)
	factory.RegisterFunc("defineObjectLite", bu.DefineObjectLite)
	defer factory.UnRegisterFunc("defineObjectLite")
	defer factory.UnRegisterFunc("defineObject")
	defer factory.UnRegisterFunc("defineItem")

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
		log.Errorf("create service file failed, name:%s, err:%s", packagePtr.Name, serviceErr.Error())
		return
	}
	return
}
