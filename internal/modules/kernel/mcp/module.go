package mcp

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initator"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	rc "github.com/muidea/magicToolkit/internal/initators/routeregistry/pkg/common"

	"github.com/muidea/magicToolkit/internal/modules/kernel/mcp/biz"
	"github.com/muidea/magicToolkit/internal/modules/kernel/mcp/pkg/common"
	"github.com/muidea/magicToolkit/internal/modules/kernel/mcp/service"
)

func init() {
	module.Register(New())
}

type MCP struct {
	bizPtr     *biz.MCP
	servicePtr *service.MCP
}

func New() *MCP {
	return &MCP{}
}

func (s *MCP) ID() string {
	return common.MCPModule
}

func (s *MCP) Setup(
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {
	var routeRegistryHelperVal rc.RouteRegistryHelper
	var routeRegistryHelperErr *cd.Error
	routeRegistryHelperVal, routeRegistryHelperErr = initator.GetEntity(rc.RouteRegistryInitator, routeRegistryHelperVal)
	if routeRegistryHelperErr != nil {
		err = routeRegistryHelperErr
		log.Errorf("MCP setup failed, initator.GetEntity error:%s", routeRegistryHelperErr.Error())
		return
	}

	s.bizPtr = biz.New(eventHub, backgroundRoutine)

	s.servicePtr = service.New(routeRegistryHelperVal.GetRouteRegistry(), s.bizPtr)

	return nil
}

func (s *MCP) Run() (err *cd.Error) {
	err = s.bizPtr.Initialize()
	if err != nil {
		return
	}

	s.servicePtr.RegisterRoute()
	return
}

func (s *MCP) Teardown() {
}
