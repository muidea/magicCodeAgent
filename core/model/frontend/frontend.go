package frontend

import (
	log "github.com/cihub/seelog"
	"github.com/muidea/magicToolkit/core/model/frontend/data"
	"github.com/muidea/magicToolkit/core/model/frontend/models"
	"github.com/muidea/magicToolkit/core/model/frontend/services"
	"github.com/muidea/magicToolkit/core/model/frontend/views"
	"github.com/muidea/magicToolkit/parser"
)

type Frontend struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Frontend {
	return &Frontend{rootPath: rootPath, projectName: projectName}
}

func (s *Frontend) createModule(ptr *parser.Package) (err error) {
	return
}

func (s *Frontend) Create(packagePtr *parser.Package) (err error) {
	moduleErr := s.createModule(packagePtr)
	if moduleErr != nil {
		log.Errorf("create module failed, name:%s, err:%s", packagePtr.Name, moduleErr.Error())
		err = moduleErr
		return
	}

	dataPtr := data.New(s.rootPath, s.projectName)
	dataErr := dataPtr.Create(packagePtr)
	if dataErr != nil {
		log.Errorf("create data file failed, name:%s, err:%s", packagePtr.Name, dataErr.Error())
		return
	}

	modelsPtr := models.New(s.rootPath, s.projectName)
	modelsErr := modelsPtr.Create(packagePtr)
	if modelsErr != nil {
		log.Errorf("create models file failed, name:%s, err:%s", packagePtr.Name, modelsErr.Error())
		return
	}

	servicesPtr := services.New(s.rootPath, s.projectName)
	servicesErr := servicesPtr.Create(packagePtr)
	if servicesErr != nil {
		log.Errorf("create servies file failed, name:%s, err:%s", packagePtr.Name, servicesErr.Error())
		return
	}

	viewsPtr := views.New(s.rootPath, s.projectName)
	viewsErr := viewsPtr.Create(packagePtr)
	if viewsErr != nil {
		log.Errorf("create view file failed, name:%s, err:%s", packagePtr.Name, viewsErr.Error())
		return
	}
	return
}
