package pprof

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/foundation/log"
	"github.com/muidea/magicCommon/framework/plugin/initator"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/config"
	"github.com/muidea/magicCodeAgent/internal/initators/pprof/pkg/common"
)

func init() {
	initator.Register(New())
}

type pprof struct{}

func New() *pprof {
	return &pprof{}
}

func (s *pprof) ID() string {
	return common.PProfInitator
}

func (s *pprof) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) {
	addr := ":1" + config.BindPort()
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Criticalf("funcRetErr=http.ListenAndServe||err=%s", err.Error())
		}
	}()

	log.Infof("pprof workspace start on %s", addr)
}

func (s *pprof) Run() {
	// nothing
}

func (s *pprof) Teardown() {

}
