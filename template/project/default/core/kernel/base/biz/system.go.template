package biz

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	log "github.com/cihub/seelog"

	cc "github.com/muidea/magicCas/common"
)

const moduleDefine = `
{
    "module": [
        {
            "name": "Partner",
            "icon": "icon-partner",
            "title": "会员管理",
			"content": "会员信息管理，支持会员信息的查询，新增，修改，删除操作，不同权限组的账号可以执行不同的操作。支持会员信息推荐，每个新会员可以由一名会员推荐。不允许同一个手机号被多个会员使用",
            "child": [
                {
                    "name": "Partner",
                    "icon": "icon-subaccount",
                    "title": "会员信息",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/",
                        "update/:id",
                        "delete/:id"
                    ]
                }
            ]
        },
        {
            "name": "Order",
            "icon": "icon-cart",
            "title": "订单管理",
			"content": "订单信息管理，支持订单信息的查询，新增操作。",
            "child": [
                {
                    "name": "Order",
                    "icon": "icon-order",
                    "title": "订单信息",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/"
                    ]
                }
            ]
        },
        {
            "name": "Bill",
            "icon": "icon-card",
            "title": "积分管理",
			"content": "积分信息管理，支持积分清单、积分报告、积分兑换、积分策略的查询操作，支持积分兑换，积分策略新建操作，支持积分策略的更新和删除操作。不同权限组的账号可以执行不同的操作。",
            "child": [
                {
                    "name": "Bill",
                    "icon": "icon-bill",
                    "title": "积分清单",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/"
                    ]
                },
                {
                    "name": "BillReport",
                    "icon": "icon-bulletin",
                    "title": "积分报表",
                    "action": [
                        "query/",
                        "query/:id"
                    ]
                },
                {
                    "name": "PayReward",
                    "icon": "icon-pay_collect",
                    "title": "积分兑换",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/"
                    ]
                },
                {
                    "name": "RewardPolicy",
                    "icon": "icon-balance",
                    "title": "积分策略",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/",
                        "update/:id",
                        "delete/:id"
                    ]
                }
            ]
        },
        {
            "name": "Product",
            "icon": "icon-catalog",
            "title": "产品管理",
			"content": "产品信息管理，支持产品信息的查询，新增，修改，删除操作，不同权限组的账号可以执行不同的操作。",
            "child": [
                {
                    "name": "Product",
                    "icon": "icon-gift",
                    "title": "产品信息",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/",
                        "update/:id",
						"delete/:id"
                    ]
                }
            ]
        },
        {
            "name": "Store",
            "icon": "icon-huititle",
            "title": "店铺管理",
			"content": "店铺信息管理，支持店铺信息的查询，新增，修改，删除操作，支持产品入库，商品出库，商品信息查看管理，不同权限组的账号可以执行不同的操作。",
            "child": [
                {
                    "name": "StockIn",
                    "icon": "icon-downland",
                    "title": "入库管理",
                    "action": [
                        "query/",
                        "query/:id",
						"update/:id",
                        "create/"
                    ]
                },
                {
                    "name": "StockOut",
                    "icon": "icon-upload",
                    "title": "出库管理",
                    "action": [
                        "query/",
                        "query/:id",
						"update/:id",
                        "create/"
                    ]
                },
                {
                    "name": "Goods",
                    "icon": "icon-shop",
                    "title": "商品管理",
                    "action": [
                        "query/",
                        "query/:id",
                        "update/:id"
                    ]
                },
                {
                    "name": "Store",
                    "icon": "icon-compass",
                    "title": "店铺信息",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/",
                        "update/:id",
                        "delete/:id"
                    ]
                }
            ]
        },
        {
            "name": "Warehouse",
            "icon": "icon-warehouse",
            "title": "仓库管理",
			"content": "仓库信息管理，支持仓库信息的查询，新增，修改，删除操作，支持货架信息查询，新增，修改，删除操作，不同权限组的账号可以执行不同的操作。",
            "child": [
                {
                    "name": "Shelf",
                    "icon": "icon-barchart",
                    "title": "货架管理",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/",
                        "update/:id",
                        "delete/:id"
                    ]
                },
                {
                    "name": "Warehouse",
                    "icon": "icon-home",
                    "title": "仓库信息",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/",
                        "update/:id",
                        "delete/:id"
                    ]
                }
            ]
        }
    ]
}
`
const systemDefine = `
{
    "module": [
        {
            "name": "Authority",
            "icon": "icon-people",
            "title": "账号管理",
			"content": "账号信息管理，支持账号信息的查询，新增，修改，删除操作，不同权限组的账号可以执行不同的操作。",
            "child": [
                {
                    "name": "Account",
                    "icon": "icon-me",
                    "title": "账号信息",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/",
                        "update/:id",
                        "delete/:id"
                    ]
                },
                {
                    "name": "Endpoint",
                    "icon": "icon-endpoint",
                    "title": "接入终端",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/",
                        "update/:id",
                        "delete/:id"
                    ]
                },
                {
                    "name": "Role",
                    "icon": "icon-private",
                    "title": "权限信息",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/",
                        "update/:id",
                        "delete/:id"
                    ]
                },
                {
                    "name": "Namespace",
                    "icon": "icon-namespace",
                    "title": "租户信息",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/",
                        "update/:id",
                        "delete/:id"
                    ]
                }
            ]
        },
        {
            "name": "Image",
            "icon": "icon-folder",
            "title": "图片管理",
			"content": "图片信息管理，支持对系统里上传的图片信息查询，更新操作，不同权限组的账号可以执行不同的操作。",
            "child": [
				{
                    "name": "Image",
                    "icon": "icon-image",
                    "title": "图片信息",
                    "action": [
                        "query/",
                        "query/:id",
                        "create/",
                        "update/:id",
                        "delete/:id"
                    ]
                }
            ]
        },
        {
            "name": "Setting",
            "icon": "icon-setting",
            "title": "系统设置",
			"content": "系统信息管理，支持系统信息的查询，更新操作，不同权限组的账号可以执行不同的操作。",
            "child": [
                {
                    "name": "Profile",
                    "icon": "icon-personal",
                    "title": "个人信息",
                    "action": [
                        "query/"
                    ]
                },
				{
                    "name": "Setting",
                    "icon": "icon-info",
                    "title": "系统信息",
                    "action": [
                        "query/"
                    ]
                }
            ]
        }
    ]
}
`

type ModuleInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *ModuleInfo) FromModule(modulePtr *Module) {
	s.Name = modulePtr.Title
	s.Description = modulePtr.Content
}

type Privilege struct {
	Read   bool `json:"read"`
	Write  bool `json:"write"`
	Delete bool `json:"delete"`
}

type Route struct {
	Name      string     `json:"name"`
	Icon      string     `json:"icon"`
	Path      string     `json:"path"`
	Privilege *Privilege `json:"privilege"`
	Routes    []*Route   `json:"routes"`
}

func (s *Route) FromModule(prefix string, modulePtr *Module) {
	s.Name = modulePtr.Title
	s.Icon = modulePtr.Icon
	if prefix != "" {
		if prefix != modulePtr.Name {
			s.Path = strings.ToLower(path.Join("/", prefix, modulePtr.Name) + "/")
		} else {
			s.Path = strings.ToLower(path.Join("/", modulePtr.Name) + "/")
		}
	}
}

type Module struct {
	Name    string    `json:"name"`
	Icon    string    `json:"icon"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Action  []string  `json:"action"`
	Child   []*Module `json:"child"`
}

type Define struct {
	Module []*Module `json:"module"`
}

type System struct {
	moduleDefine *Define
	systemDefine *Define
}

func (s *System) loadModule() (err error) {
	moduleDefine, moduleErr := s.loadDefine(moduleDefine)
	if moduleErr != nil {
		err = moduleErr
		return
	}
	s.moduleDefine = moduleDefine

	systemDefine, systemErr := s.loadDefine(systemDefine)
	if systemErr != nil {
		err = systemErr
		return
	}
	s.systemDefine = systemDefine
	return
}

func (s *System) loadDefine(strDefine string) (ret *Define, err error) {
	dataBuffer := bytes.NewBufferString(strDefine)
	dataContent, dataErr := ioutil.ReadAll(dataBuffer)
	if dataErr != nil {
		err = dataErr
		log.Errorf("ioutil.ReadAll failed,err:%s", err.Error())
		return
	}

	define := &Define{}
	err = json.Unmarshal(dataContent, define)
	if err != nil {
		return
	}
	ret = define
	return
}

func (s *System) SystemInfo(ptr *cc.RoleView, isSuper bool) (route []*Route, content []*ModuleInfo) {
	err := s.loadModule()
	if err != nil {
		log.Errorf("load module failed, err:%s", err.Error())
		return
	}

	//if ptr == nil && isSuper {
	//	moduleRoutes := getRoutes(s.moduleDefine)
	//	ret = append(ret, moduleRoutes...)
	//	systemRoutes := getRoutes(s.systemDefine)
	//	ret = append(ret, systemRoutes...)
	//	return
	//}

	route = []*Route{}
	content = []*ModuleInfo{}
	if !isSuper {
		for _, val := range s.moduleDefine.Module {
			rt := &Route{}
			rt.FromModule("", val)
			for _, sv := range val.Child {
				st := &Route{}
				st.FromModule(val.Name, sv)
				privilegePtr := s.checkPrivilege(st.Path, ptr)
				if privilegePtr != nil {
					st.Privilege = privilegePtr
					rt.Routes = append(rt.Routes, st)
				}
			}
			if len(rt.Routes) > 0 {
				route = append(route, rt)

				cnt := &ModuleInfo{}
				cnt.FromModule(val)

				content = append(content, cnt)
			}
		}
	}

	for _, val := range s.systemDefine.Module {
		rt := &Route{}
		rt.FromModule("", val)
		for _, sv := range val.Child {
			st := &Route{}
			st.FromModule(val.Name, sv)

			if s.IsSpecialRoute(st.Path) && !isSuper {
				continue
			}

			if s.IsImageRoute(st.Path) && isSuper {
				continue
			}

			privilegePtr := s.checkPrivilege(st.Path, ptr)
			if privilegePtr != nil {
				st.Privilege = privilegePtr
				rt.Routes = append(rt.Routes, st)
			}
		}

		if len(rt.Routes) > 0 {
			route = append(route, rt)

			cnt := &ModuleInfo{}
			cnt.FromModule(val)

			content = append(content, cnt)
		}
	}

	return
}

var patternReg, _ = regexp.Compile("(/[a-z]+)+/(query|create|update|delete)")

func (s *System) checkPrivilege(prefix string, ptr *cc.RoleView) (ret *Privilege) {
	if ptr == nil {
		return
	}

	if ptr.IsSuper() {
		ret = &Privilege{Read: true, Write: true, Delete: true}
		return
	}

	privates := map[int]int{}
	for _, val := range ptr.Privilege {
		subStr := patternReg.FindString(val.Path)
		if strings.Index(subStr, prefix) == 0 {
			privates[val.Value.Value] = val.Value.Value
		}
	}

	if len(privates) == 0 {
		return
	}

	ret = &Privilege{}
	for k, _ := range privates {
		switch k {
		case cc.ReadPermission:
			ret.Read = true
		case cc.WritePermission:
			ret.Write = true
		case cc.DeletePermission:
			ret.Delete = true
		}
	}

	return
}

var namespaceRoute = "/authority/namespace"
var fullNamespaceRoute = "/api/v1/authority/namespace"

func (s *System) IsSpecialRoute(prefix string) bool {
	if 0 == strings.Index(prefix, namespaceRoute) {
		return true
	}

	if 0 == strings.Index(prefix, fullNamespaceRoute) {
		return true
	}

	return false
}

var imageRoute = "/image"
var fullImageRoute = "/api/v1/image"

func (s *System) IsImageRoute(prefix string) bool {
	if 0 == strings.Index(prefix, imageRoute) {
		return true
	}

	if 0 == strings.Index(prefix, fullImageRoute) {
		return true
	}

	return false
}

var authorityRoute = "/authority/"
var fullAuthorityRoute = "/api/v1/authority/"
var settingRoute = "/setting/"
var fullSettingRoute = "/api/v1/setting/"

func (s *System) IsSettingRoute(prefix string) bool {
	if 0 == strings.Index(prefix, authorityRoute) {
		return true
	}
	if 0 == strings.Index(prefix, fullAuthorityRoute) {
		return true
	}

	if 0 == strings.Index(prefix, settingRoute) {
		return true
	}
	if 0 == strings.Index(prefix, fullSettingRoute) {
		return true
	}

	return false
}
