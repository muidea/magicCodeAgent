package recorder

import (
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/framework/plugin/module"
	"github.com/muidea/magicCommon/task"

	"github.com/muidea/magicCodeAgent/internal/modules/blocks/recorder/biz"
	"github.com/muidea/magicCodeAgent/internal/modules/blocks/recorder/pkg/common"
)

func init() {
	module.Register(New())
}

func New() *Recorder {
	return &Recorder{}
}

type Recorder struct {
	bizPtr *biz.Recorder
}

func (s *Recorder) ID() string {
	return common.RecorderModule
}

func (s *Recorder) Weight() int {
	return 999
}

func (s *Recorder) Setup(eventHub event.Hub, backgroundRoutine task.BackgroundRoutine) (err *cd.Error) {

	s.bizPtr = biz.New(eventHub, backgroundRoutine)
	return
}

func (s *Recorder) Run() (err *cd.Error) {
	s.bizPtr.Run()
	return
}

func (s *Recorder) Teardown() {
}
