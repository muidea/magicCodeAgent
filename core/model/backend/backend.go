package backend

import (
	log "github.com/cihub/seelog"

	"github.com/muidea/magicToolkit/core/model/backend/biz"
	"github.com/muidea/magicToolkit/core/model/backend/common"
	"github.com/muidea/magicToolkit/core/model/backend/dao"
	"github.com/muidea/magicToolkit/core/model/backend/service"
	mu "github.com/muidea/magicToolkit/core/model/util"
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
	factory.RegisterFunc("isComplex", mu.IsComplex)
	factory.RegisterFunc("defineItem", mu.DefineItem)
	factory.RegisterFunc("defineObject", mu.DefineObject)
	factory.RegisterFunc("defineObjectLite", mu.DefineObjectLite)
	factory.RegisterFunc("formatModuleName", mu.FormatModuleName)
	factory.RegisterFunc("formatObjectName", mu.FormatObjectName)
	factory.RegisterFunc("formatObjectActionName", mu.FormatObjectActionName)
	factory.RegisterFunc("formatModuleVal", mu.FormatModuleVal)
	factory.RegisterFunc("formatObjectVal", mu.FormatObjectVal)
	factory.RegisterFunc("formatObjectActionVal", mu.FormatObjectActionVal)
	factory.RegisterFunc("formatModulePath", mu.FormatModulePath)
	factory.RegisterFunc("formatObjectJson", mu.FormatObjectJson)

	defer factory.UnRegisterFunc("formatObjectJson")
	defer factory.UnRegisterFunc("formatModulePath")
	defer factory.UnRegisterFunc("formatObjectActionVal")
	defer factory.UnRegisterFunc("formatObjectVal")
	defer factory.UnRegisterFunc("formatModuleVal")
	defer factory.UnRegisterFunc("formatObjectActionName")
	defer factory.UnRegisterFunc("formatObjectName")
	defer factory.UnRegisterFunc("formatModuleName")
	defer factory.UnRegisterFunc("defineObjectLite")
	defer factory.UnRegisterFunc("defineObject")
	defer factory.UnRegisterFunc("defineItem")
	defer factory.UnRegisterFunc("isComplex")

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
