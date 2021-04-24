package client

import (
	"github.com/muidea/magicToolkit/core/util"
	"github.com/muidea/magicToolkit/model"
)

func init() {
	util.RegisterFunc("clientServiceFuncDeclare", ServiceFuncDeclare)
	util.RegisterFunc("clientServiceFuncDefine", ServiceFuncDefine)
}

func ServiceFuncDeclare(servicePtr *model.Module) string {
	serviceVal := map[string]interface{}{"Service": servicePtr}
	factoryPtr := util.NewValueFactory("serviceFuncDeclare.tpl", serviceVal)
	valStr, valErr := factoryPtr.Create("client")
	if valErr == nil {
		return valStr
	}

	return ""
}

func ServiceFuncDefine(servicePtr *model.Module) string {
	serviceVal := map[string]interface{}{"Service": servicePtr}
	factoryPtr := util.NewValueFactory("serviceFuncDefine.tpl", serviceVal)
	valStr, valErr := factoryPtr.Create("client")
	if valErr == nil {
		return valStr
	}

	return ""
}
