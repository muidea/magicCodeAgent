package common

import (
	engine "github.com/muidea/magicEngine/http"
)

const RouteRegistryInitator = "/initators/routeregistry"

type RouteRegistryHelper interface {
	GetRouteRegistry() (ret engine.RouteRegistry)
}
