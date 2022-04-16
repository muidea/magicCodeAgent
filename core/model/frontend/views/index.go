package views

import (
	"os"
	"path"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

type indexer struct {
}

func (s *indexer) declareImport(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	importDeclare := `import React, {FC, useState} from 'react';
import { {{- formatObjectName .ObjectVal}}ModelState, connect, ConnectProps, Dispatch} from 'umi';
import {Button, Divider, Popconfirm} from 'antd';
import {SearchOutlined} from '@ant-design/icons';
import {PageContainer} from '@ant-design/pro-layout';
import ProTable, {ProColumns} from '@ant-design/pro-table';
import {formatTimeStamp, getAuthority, verifyAuthority, verifyOptions} from '@/utils';
import {InputFilter, OptionsFilter} from '@/components';
import {TotalizerSummary, ValueInfo} from '@/pages/components';
import {defaultPagination} from '@/data';
import * as contentDef from '@/data/content';
import * as baseDef from '@/data/base';
import Create from './create';
import Update from './update';
`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *indexer) declareInterface(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	interfaceDeclare := `
interface {{formatObjectName .ObjectVal}}Props extends ConnectProps {
  {{formatObjectValue .ObjectVal}}: {{formatObjectName .ObjectVal}}ModelState;
  dispatch: Dispatch;
}
`
	factoryPtr := factory.NewMemoryFactor(interfaceDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *indexer) declareFunc(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	funcDeclare := `
  const onFilter{{formatObjectName .ObjectVal}} = (value: any) => {
    dispatch({type: '{{formatObjectValue .ObjectVal}}/filter{{formatObjectName .ObjectVal}}', payload: {value}})
  }

  const onDelete{{formatObjectName .ObjectVal}} = (value: contentDef.{{formatObjectName .ObjectVal}}) => {
    dispatch({type: '{{formatObjectValue .ObjectVal}}/delete{{formatObjectName .ObjectVal}}', payload: {value}})
  }

  const onSubmit{{formatObjectName .ObjectVal}} = (action: string, value: any) => {
    dispatch({type: '{{formatObjectValue .ObjectVal}}/submit{{formatObjectName .ObjectVal}}', payload: {action, value}})
  }

  const checkOptionPanel = () => {
    return !verifyOptions(authorities, [baseDef.Authority.WritePrivate, baseDef.Authority.DeletePrivate])
  }

  const getCreatePanel = () => {
    if (verifyAuthority(authorities, baseDef.Authority.WritePrivate)) {
      return <Create key={'1'} title={'新建'} onSubmitHandle={onSubmit{{formatObjectName .ObjectVal}}}/>
    }
  }

  const getUpdatePanel = (record: contentDef.{{formatObjectName .ObjectVal}}) => {
    if (verifyAuthority(authorities, baseDef.Authority.WritePrivate)) {
      return <Update key={'2'} title={'更新'} initialValue={record} onSubmitHandle={onSubmit{{formatObjectName .ObjectVal}}}/>
    }
  }

  const getDeletePanel = (record: contentDef.{{formatObjectName .ObjectVal}}) => {
    if (verifyAuthority(authorities, baseDef.Authority.DeletePrivate)) {
      return <Popconfirm
        title='确认删除?'
        onConfirm={() => {
          onDelete{{formatObjectName .ObjectVal}}(record)
        }}
      >
        <Button type='text' size={'small'}>删除</Button>
      </Popconfirm>
    }
  }

  const onSubmitFilter = (dataIndex: string, filterVal: React.Key[]) => {
    currentFilter[dataIndex] = filterVal

    setCurrentFilter(currentFilter)
    onFilter{{formatObjectName .ObjectVal}}(currentFilter)
  }

  const onResetFilter = (dataIndex: string) => {
    delete currentFilter[dataIndex]

    setCurrentFilter(currentFilter)
    onFilter{{formatObjectName .ObjectVal}}(currentFilter)
  }

  const onTableChange = (pagination: any) => {
    if (pagination.current !== currentPagination.current) {
      onFilter{{formatObjectName .ObjectVal}}({...currentFilter, pageNum: pagination.current, pageSize: pagination.pageSize})
    }

    setCurrentPagination(pagination)
  }`
	factoryPtr := factory.NewMemoryFactor(funcDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *indexer) declareColumns(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	columnsDeclare := `
  const columns: ProColumns<contentDef.{{formatObjectName .ObjectVal}}>[] = [
    {
      title: '标题',
      dataIndex: 'title',
      filters: false,
      filterDropdown: (props) => (
        <InputFilter {...props}
                     dataIndex={'title'}
                     placeholder={'请输入标题'}
                     onSubmitFilter={onSubmitFilter}
                     onResetFilter={onResetFilter}/>
      ),
      filterIcon: (filtered) => (
        <SearchOutlined style={{"{{color: filtered ? '#1890ff' : undefined}}"}}/>
      ),
    },
    {
      title: '标签',
      dataIndex: 'tags',
      render: (text, record) => {
        return <ValueInfo initialValue={record.tags} valueFormatter={(val)=>{return val}}/>
      },
    },
    {
      title: '创建时间',
      dataIndex: 'createTime',
      render: (text, record) => {
        return <span>{formatTimeStamp(record.createTime)}</span>
      }
    },
    {
      title: '操作',
      dataIndex: 'option',
      hideInTable: checkOptionPanel(),
      render: (_, record) => (
        <div>
          {getUpdatePanel(record)}
          <Divider type="vertical" dashed/>
          {getDeletePanel(record)}
        </div>
      ),
    },
  ]`
	factoryPtr := factory.NewMemoryFactor(columnsDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *indexer) declareReturn(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	returnDeclare := `
  const labels = [{id: 1, label: '文章总数'}, {id: 2, label: '本周新增'}, {id: 3, label: '本月新增'}]
  return (
    <PageContainer
      header={{"{{title: '文章信息'}}"}}
      content={(<TotalizerSummary labels={labels} values={summaryList}/>)}
      tabList={[
        {
          tab: '文章列表',
          key: 'base',
        },
      ]}
      extra={[
        getCreatePanel()
      ]}
    >
      <ProTable<contentDef.{{formatObjectName .ObjectVal}}>
        cardProps={{"{{bodyStyle: {padding: 0},}}"}}
        size={'small'}
        rowKey={'id'}
        columns={columns}
        pagination={{"{{showSizeChanger: false, ...currentPagination, total: itemSize}}"}}
        search={false}
        options={false}
        dataSource={ {{- formatObjectValue .ObjectVal}}List}
        onChange={onTableChange}
      />
    </PageContainer>
  );`
	factoryPtr := factory.NewMemoryFactor(returnDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *indexer) declareObject(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	objectDeclare := `
const {{formatObjectValue .ObjectVal}}: FC<{{formatObjectName .ObjectVal}}Props> = ({ {{- formatObjectValue .ObjectVal}},dispatch, location, children}) => {
  const [currentFilter, setCurrentFilter] = useState<any>({})
  const [currentPagination, setCurrentPagination] = useState<any>({current: defaultPagination.pageNum, pageSize: defaultPagination.pageSize})

  const {summaryList, itemSize, {{formatObjectValue .ObjectVal}}List} = {{formatObjectValue .ObjectVal}}
  const authorities = getAuthority(location.pathname)
{{declareFunc .ModuleVal .ObjectVal}}
{{declareColumns .ModuleVal .ObjectVal}}
{{declareReturn .ModuleVal .ObjectVal}}
}`
	factory.RegisterFunc("declareFunc", s.declareFunc)
	factory.RegisterFunc("declareColumns", s.declareColumns)
	factory.RegisterFunc("declareReturn", s.declareReturn)
	defer factory.UnRegisterFunc("declareReturn")
	defer factory.UnRegisterFunc("declareColumns")
	defer factory.UnRegisterFunc("declareFunc")
	factoryPtr := factory.NewMemoryFactor(objectDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *indexer) declareExport(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	exportDeclare := `
export default connect(({ {{- formatObjectValue .ObjectVal -}} }: { {{- formatObjectValue .ObjectVal}}: {{formatObjectName .ObjectVal}}ModelState,}) => ({
  {{formatObjectValue .ObjectVal}},
}))({{formatObjectValue .ObjectVal}});`
	factoryPtr := factory.NewMemoryFactor(exportDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *indexer) index(ptr *parser.Package, objPtr *parser.Object, commonPath string) (err error) {
	commonFile := path.Join(commonPath, "index.tsx")
	log.Info(commonFile)

	fileHandle, fileErr := os.OpenFile(commonFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if fileErr != nil {
		err = fileErr
		return
	}
	defer fileHandle.Close()

	_, err = fileHandle.WriteString(s.declareImport(ptr, objPtr))
	if err != nil {
		return
	}
	_, err = fileHandle.WriteString(s.declareInterface(ptr, objPtr))
	if err != nil {
		return
	}
	_, err = fileHandle.WriteString(s.declareObject(ptr, objPtr))
	if err != nil {
		return
	}
	_, err = fileHandle.WriteString(s.declareExport(ptr, objPtr))
	if err != nil {
		return
	}
	return
}
