package views

import (
	"os"
	"path"

	log "github.com/cihub/seelog"

	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

type updater struct {
}

func (s *updater) declareImport(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	importDeclare := `import React from 'react';
import {Button} from 'antd';
import ProForm, {DrawerForm, ProFormSelect, ProFormText, ProFormTextArea} from '@ant-design/pro-form';
import * as contentDef from "@/data/content";
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

func (s *updater) declareInterface(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	interfaceDeclare := `
interface updateProps {
  title: string;
  initialValue: contentDef.{{formatObjectName .ObjectVal}};
  onSubmitHandle: (action: string, value: any) => void;
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

func (s *updater) declareObject(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	objectDeclare := `
const Update: React.FC<updateProps> = (props) => {
  const {title, initialValue, onSubmitHandle} = props;

  return (
    <DrawerForm
      title={title}
      width="500px"
      drawerProps={ {destroyOnClose: true,} }
      trigger={<Button type="text" size={"small"}>{title}</Button>}
      onVisibleChange={(visible)=>{
        if (!visible){
          return
        }
      }}
      onFinish={async (values) => {
        onSubmitHandle("update", {
          ...values,
          id: initialValue.id,
        })
        return true;
      }}
    >
      <ProForm.Group>
        <ProFormText
          name="title"
          label="标题"
          width={"lg"}
          placeholder="请输入标题"
          initialValue={initialValue.title}
          rules={[{required: true, message: '请输入标题'}]}
        />
      </ProForm.Group>
      <ProForm.Group>
        <ProFormTextArea
          name="content"
          label="内容"
          width={"lg"}
          placeholder="请输入内容"
          initialValue={initialValue.content}
          rules={[{required: true, message: '请输入内容'}]}
        />
      </ProForm.Group>
      <ProForm.Group>
        <ProFormSelect
          mode="tags"
          name="tags"
          label="标签"
          width={"lg"}
          placeholder="请输入标签"
          initialValue={initialValue?.tags?.length > 0? initialValue.tags:undefined}
          valueEnum={ {} }
        />
      </ProForm.Group>
    </DrawerForm>
  );
};
`
	factoryPtr := factory.NewMemoryFactor(objectDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *updater) declareExport(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	exportDeclare := `
export default Update;`
	factoryPtr := factory.NewMemoryFactor(exportDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *updater) update(ptr *parser.Package, objPtr *parser.Object, commonPath string) (err error) {
	commonFile := path.Join(commonPath, "update.tsx")
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
