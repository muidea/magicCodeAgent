package models

import (
	"fmt"
	"os"
	"path"

	log "github.com/cihub/seelog"
	mu "github.com/muidea/magicToolkit/core/model/util"
	"github.com/muidea/magicToolkit/factory"
	"github.com/muidea/magicToolkit/parser"
)

type Models struct {
	rootPath    string
	projectName string
}

func New(rootPath, projectName string) *Models {
	return &Models{rootPath: rootPath, projectName: projectName}
}

func (s *Models) declareImport(ptr *parser.Package) (ret string) {
	importDeclare := `
import {Effect, Reducer, Subscription} from 'umi';
import {showErrorMessage} from "@/components";
import {defaultPagination} from "@/data";
import api from "@/pages/{{formatModulePath .ModuleVal}}/services";
`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareInterface(ptr *parser.Package) (ret string) {
	interfaceDeclare := `{{$ModuleVal := .ModuleVal}}
export interface {{formatModuleName $ModuleVal}}ModelState {
  summaryList: [];
  itemSize: number;
{{- with .ModuleVal -}}{{range .Objects}}
  {{formatObjectVal .}}List: any[];
  current{{formatObjectName .}}?: any;
{{end}}{{end}}
}

export interface {{formatModuleName $ModuleVal}}ModelType {
  namespace: '{{formatModuleName $ModuleVal}}';
  state: {{formatModuleName $ModuleVal}}ModelState;
  effects: {
{{- with .ModuleVal -}}{{range .Objects}}
    filter{{formatObjectName .}}: Effect;
    delete{{formatObjectName .}}: Effect;
    submit{{formatObjectName .}}: Effect;
    create{{formatObjectName .}}: Effect;
    update{{formatObjectName .}}: Effect;
{{end}}{{end}}
  };
  reducers: {
    save: Reducer<{{formatModuleName $ModuleVal}}ModelState>;
    clean: Reducer<{{formatModuleName $ModuleVal}}ModelState>;
  };
  subscriptions: { setup: Subscription };
}
`
	factoryPtr := factory.NewMemoryFactor(interfaceDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareSubscriptions(ptr *parser.Package) (ret string) {
	subscriptionsDeclare := `{{$ModuleVal := .ModuleVal}}
  subscriptions: {
    setup({dispatch, history}) {
      return history.listen((location) => {
		{{- with .ModuleVal -}}{{range .Objects}}
		if (location.pathname === '/{{formatObjectActionName $ModuleVal .}}') {
		  dispatch({
			type: 'filter{{formatObjectName .}}',
			payload: {
			  value: {...defaultPagination},
			}
		  })
          return
		}
		{{end}}{{end}}
      });
    }
  },`
	factoryPtr := factory.NewMemoryFactor(subscriptionsDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareEffects(ptr *parser.Package) (ret string) {
	effectsDeclare := `{{$ModuleVal := .ModuleVal}}
  effects: {
	{{- with .ModuleVal -}}{{range .Objects}}
    * filter{{formatObjectName .}}({payload}, {call, put}) {
      const data = yield call(api.filter{{formatObjectName .}})
      const {errorCode, reason} = data
      if (errorCode === 0) {
        const { {{formatObjectValue .}} } = data
        yield put({
          type: 'save',
          payload: {
            {{formatObjectValue .}}List: {{formatObjectValue .}},
          },
        });
        return
      }
      showErrorMessage(reason);
    },
    * query{{formatObjectName .}}({payload}, {call, put}) {
      const {value} = payload
      const data = yield call(api.query{{formatObjectName .}}, {id: value.id})
      const {errorCode, reason} = data
      if (errorCode === 0) {
        const { {{formatObjectValue .}} } = data
        yield put({
          type: 'save',
          payload: {
            current{{formatObjectName .}}: {{formatObjectValue .}},
          },
        });
        return
      }
      showErrorMessage(reason);
    },
    * submit{{formatObjectName .}}({payload}, {put}) {
      const {action, value} = payload
      if (action === 'create') {
        yield put({
          type: 'create{{formatObjectName .}}',
          payload: {
            value,
          }
        })
        return
      }
      if (action === 'update') {
        yield put({
          type: 'update{{formatObjectName .}}',
          payload: {
            value,
          }
        })
        return
      }
    },
    * create{{formatObjectName .}}({payload}, {call, put}) {
      const {value} = payload
      const data = yield call(api.create{{formatObjectName .}}, {...value})
      const {errorCode, reason} = data
      if (errorCode === 0) {
        yield put({
          type: 'filter{{formatObjectName .}}',
          payload: {
            value: {...defaultPagination},
          }
        });
        return
      }
      showErrorMessage(reason);
    },
    * update{{formatObjectName .}}({payload}, {call, put}) {
      const {value} = payload
      const data = yield call(api.update{{formatObjectName .}}, {...value})
      const {errorCode, reason} = data
      if (errorCode === 0) {
        yield put({
          type: 'filter{{formatObjectName .}}',
          payload: {
            value: {...defaultPagination},
          }
        });
        return
      }
      showErrorMessage(reason);
    },
    * delete{{formatObjectName .}}({payload}, {call, put}) {
      const {value} = payload
      const data = yield call(api.delete{{formatObjectName .}}, {id: value.id})
      const {errorCode, reason} = data
      if (errorCode === 0) {
        yield put({
          type: 'filter{{formatObjectName .}}',
          payload: {
            value: {...defaultPagination},
          }
        });
        return
      }
      showErrorMessage(reason);
    },
	{{end}}{{end}}
  },`
	factoryPtr := factory.NewMemoryFactor(effectsDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareReducers(ptr *parser.Package) (ret string) {
	reducersDeclare := `
  reducers: {
    save(state, action) {
      return {
        ...state,
        ...action.payload,
      };
    },

    clean(state, action) {
      return {
        ...state,
        ...action.payload,
      };
    },
  },
`
	factoryPtr := factory.NewMemoryFactor(reducersDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareObject(ptr *parser.Package) (ret string) {
	objectDeclare := `{{$ModuleVal := .ModuleVal}}
const {{formatModuleName $ModuleVal}}Model: {{formatModuleName $ModuleVal}}ModelType = {
  namespace: '{{formatModuleName $ModuleVal}}',

  state: {
    summaryList: [],
    itemSize: 0,
{{- with .ModuleVal -}}{{range .Objects}}
  {{toLower .Name}}List: [],
{{end}}{{end}}
  },
{{declareSubscriptions $ModuleVal}}
{{declareEffects $ModuleVal}}
{{declareReducers $ModuleVal}}
};`

	factory.RegisterFunc("declareSubscriptions", s.declareSubscriptions)
	factory.RegisterFunc("declareEffects", s.declareEffects)
	factory.RegisterFunc("declareReducers", s.declareReducers)
	defer factory.UnRegisterFunc("declareReducers")
	defer factory.UnRegisterFunc("declareEffects")
	defer factory.UnRegisterFunc("declareSubscriptions")

	factoryPtr := factory.NewMemoryFactor(objectDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareExport(ptr *parser.Package) (ret string) {
	exportDeclare := `{{$ModuleVal := .ModuleVal}}
export default {{formatModuleName $ModuleVal}}Model;
`
	factoryPtr := factory.NewMemoryFactor(exportDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) Create(ptr *parser.Package) (err error) {
	if len(ptr.Objects) == 0 {
		return
	}

	moduleName := mu.FormatModulePath(ptr)
	commonFile := fmt.Sprintf("/%s/models/index.ts", moduleName)
	commonFilePath := path.Join(s.rootPath, "pages", commonFile)
	log.Info(commonFilePath)
	func() {
		fileHandle, fileErr := os.OpenFile(commonFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if fileErr != nil {
			err = fileErr
			return
		}
		defer fileHandle.Close()

		_, err = fileHandle.WriteString(s.declareImport(ptr))
		if err != nil {
			return
		}
		_, err = fileHandle.WriteString(s.declareInterface(ptr))
		if err != nil {
			return
		}
		_, err = fileHandle.WriteString(s.declareObject(ptr))
		if err != nil {
			return
		}
		_, err = fileHandle.WriteString(s.declareExport(ptr))
		if err != nil {
			return
		}
	}()

	return
}
