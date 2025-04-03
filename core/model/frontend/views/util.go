package views

import (
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

func getColumns(ptr *parser.Object) (ret string) {
	return
}

func getStringColumns(ptr *parser.Item) (ret string) {
	importDeclare :=
		`    {
      title: '标题',
      dataIndex: '{{format}}',
      filters: false,
      filterDropdown: (props) => (
        <InputFilter {...props}
                     dataIndex={'title'}
                     placeholder={'请输入标题'}
                     onSubmitFilter={onSubmitFilter}
                     onResetFilter={onResetFilter}/>
      ),
      filterIcon: (filtered) => (
        <SearchOutlined style={{color: filtered ? '#1890ff' : undefined}}/>
      ),
    },`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ItemVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}
