package models

import (
	"fmt"
	"github.com/muidea/magicToolkit/core/util"
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

func (s *Models) declareImport(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	importDeclare := `import {Effect, Reducer, Subscription} from 'umi';
import {showErrorMessage} from "@/components";
import {defaultPagination} from "@/data";
import api from "@/pages/{{formatObjectActionValue .ModuleVal .ObjectVal}}/services";
`
	factoryPtr := factory.NewMemoryFactor(importDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareInterface(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	interfaceDeclare := `
export interface {{formatObjectName .ObjectVal}}ModelState {
  summaryList: [];
  itemSize: number;
  {{formatObjectValue .ObjectVal}}List: any[];
  current{{formatObjectName .ObjectVal}}?: any;
}

export interface {{formatObjectName .ObjectVal}}ModelType {
  namespace: '{{formatObjectName .ObjectVal}}';
  state: {{formatObjectName .ObjectVal}}ModelState;
  effects: {
    filter{{formatObjectName .ObjectVal}}: Effect;
    query{{formatObjectName .ObjectVal}}: Effect;
    delete{{formatObjectName .ObjectVal}}: Effect;
    submit{{formatObjectName .ObjectVal}}: Effect;
    create{{formatObjectName .ObjectVal}}: Effect;
    update{{formatObjectName .ObjectVal}}: Effect;
  };
  reducers: {
    save: Reducer<{{formatObjectName .ObjectVal}}ModelState>;
    clean: Reducer<{{formatObjectName .ObjectVal}}ModelState>;
  };
  subscriptions: { setup: Subscription };
}
`
	factoryPtr := factory.NewMemoryFactor(interfaceDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareSubscriptions(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	subscriptionsDeclare := `
  subscriptions: {
    setup({dispatch, history}) {
      return history.listen((location) => {
		if (location.pathname === '/{{formatObjectActionName .ModuleVal .ObjectVal}}') {
		  dispatch({
			type: 'filter{{formatObjectName .ObjectVal}}',
			payload: {
			  value: {...defaultPagination},
			}
		  })
          return
		}
      });
    }
  },`
	factoryPtr := factory.NewMemoryFactor(subscriptionsDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareEffects(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	effectsDeclare := `
  effects: {
    * filter{{formatObjectName .ObjectVal}}({payload}, {call, put}) {
      const data = yield call(api.filter{{formatObjectName .ObjectVal}})
      const {errorCode, reason} = data
      if (errorCode === 0) {
        const { {{formatObjectValue .ObjectVal}} } = data
        yield put({
          type: 'save',
          payload: {
            {{formatObjectValue .ObjectVal}}List: {{formatObjectValue .ObjectVal}},
          },
        });
        return
      }
      showErrorMessage(reason);
    },
    * query{{formatObjectName .ObjectVal}}({payload}, {call, put}) {
      const {value} = payload
      const data = yield call(api.query{{formatObjectName .ObjectVal}}, {id: value.id})
      const {errorCode, reason} = data
      if (errorCode === 0) {
        const { {{formatObjectValue .ObjectVal}} } = data
        yield put({
          type: 'save',
          payload: {
            current{{formatObjectName .ObjectVal}}: {{formatObjectValue .ObjectVal}},
          },
        });
        return
      }
      showErrorMessage(reason);
    },
    * submit{{formatObjectName .ObjectVal}}({payload}, {put}) {
      const {action, value} = payload
      if (action === 'create') {
        yield put({
          type: 'create{{formatObjectName .ObjectVal}}',
          payload: {
            value,
          }
        })
        return
      }
      if (action === 'update') {
        yield put({
          type: 'update{{formatObjectName .ObjectVal}}',
          payload: {
            value,
          }
        })
        return
      }
    },
    * create{{formatObjectName .ObjectVal}}({payload}, {call, put}) {
      const {value} = payload
      const data = yield call(api.create{{formatObjectName .ObjectVal}}, {...value})
      const {errorCode, reason} = data
      if (errorCode === 0) {
        yield put({
          type: 'filter{{formatObjectName .ObjectVal}}',
          payload: {
            value: {...defaultPagination},
          }
        });
        return
      }
      showErrorMessage(reason);
    },
    * update{{formatObjectName .ObjectVal}}({payload}, {call, put}) {
      const {value} = payload
      const data = yield call(api.update{{formatObjectName .ObjectVal}}, {...value})
      const {errorCode, reason} = data
      if (errorCode === 0) {
        yield put({
          type: 'filter{{formatObjectName .ObjectVal}}',
          payload: {
            value: {...defaultPagination},
          }
        });
        return
      }
      showErrorMessage(reason);
    },
    * delete{{formatObjectName .ObjectVal}}({payload}, {call, put}) {
      const {value} = payload
      const data = yield call(api.delete{{formatObjectName .ObjectVal}}, {id: value.id})
      const {errorCode, reason} = data
      if (errorCode === 0) {
        yield put({
          type: 'filter{{formatObjectName .ObjectVal}}',
          payload: {
            value: {...defaultPagination},
          }
        });
        return
      }
      showErrorMessage(reason);
    },
  },`
	factoryPtr := factory.NewMemoryFactor(effectsDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareReducers(ptr *parser.Package, objPtr *parser.Object) (ret string) {
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
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareObject(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	objectDeclare := `
const {{formatObjectName .ObjectVal}}Model: {{formatObjectName .ObjectVal}}ModelType = {
  namespace: '{{formatObjectName .ObjectVal}}',

  state: {
    summaryList: [],
    itemSize: 0,
  	{{formatObjectValue .ObjectVal}}List: [],
  },
{{declareSubscriptions .ModuleVal .ObjectVal}}
{{declareEffects .ModuleVal .ObjectVal}}
{{declareReducers .ModuleVal .ObjectVal}}
};`

	factory.RegisterFunc("declareSubscriptions", s.declareSubscriptions)
	factory.RegisterFunc("declareEffects", s.declareEffects)
	factory.RegisterFunc("declareReducers", s.declareReducers)
	defer factory.UnRegisterFunc("declareReducers")
	defer factory.UnRegisterFunc("declareEffects")
	defer factory.UnRegisterFunc("declareSubscriptions")

	factoryPtr := factory.NewMemoryFactor(objectDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ModuleVal": ptr, "ObjectVal": objPtr})
	if contentErr != nil {
		log.Errorf("create content failed, err:%s", contentErr.Error())
		return
	}

	ret = contentVal
	return
}

func (s *Models) declareExport(ptr *parser.Package, objPtr *parser.Object) (ret string) {
	exportDeclare := `
export default {{formatObjectName .ObjectVal}}Model;
`
	factoryPtr := factory.NewMemoryFactor(exportDeclare)
	contentVal, contentErr := factoryPtr.Create(map[string]interface{}{"ProjectName": s.projectName, "ModuleVal": ptr, "ObjectVal": objPtr})
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

	for _, val := range ptr.Objects {
		commonPath := fmt.Sprintf("%s/pages/%s/models", s.rootPath, mu.FormatObjectActionValue(ptr, val))
		if !util.Exists(commonPath) {
			err = os.MkdirAll(commonPath, os.ModePerm)
			if err != nil {
				return
			}
		}
		commonFile := path.Join(commonPath, "index.ts")
		log.Info(commonFile)
		func() {
			fileHandle, fileErr := os.OpenFile(commonFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if fileErr != nil {
				err = fileErr
				return
			}
			defer fileHandle.Close()

			_, err = fileHandle.WriteString(s.declareImport(ptr, val))
			if err != nil {
				return
			}
			_, err = fileHandle.WriteString(s.declareInterface(ptr, val))
			if err != nil {
				return
			}
			_, err = fileHandle.WriteString(s.declareObject(ptr, val))
			if err != nil {
				return
			}
			_, err = fileHandle.WriteString(s.declareExport(ptr, val))
			if err != nil {
				return
			}
		}()
	}

	return
}
