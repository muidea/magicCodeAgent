package biz

import (
	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
)

type Base struct {
	id                string
	eventHub          event.Hub
	backgroundRoutine *task.BackgroundRoutine
}

type invokeTask struct {
	funcPtr func()
}

func (s *invokeTask) Run() {
	s.funcPtr()
}

func New(
	id string,
	eventHub event.Hub,
	backgroundRoutine *task.BackgroundRoutine) Base {
	return Base{
		id:                id,
		eventHub:          eventHub,
		backgroundRoutine: backgroundRoutine,
	}
}

func (s *Base) ID() string {
	return s.id
}

func (s *Base) PostEvent(event event.Event) {
	s.eventHub.Post(event)
}

func (s *Base) SendEvent(event event.Event) event.Result {
	return s.eventHub.Send(event)
}

func (s *Base) CallEvent(event event.Event) event.Result {
	return s.eventHub.Call(event)
}

func (s *Base) Invoke(funcPtr func()) {
	task := &invokeTask{funcPtr: funcPtr}

	s.backgroundRoutine.Post(task)
}

func (s *Base) RootDestination() string {
	return "/#"
}

func (s *Base) InnerDestination() string {
	return s.ID()
}
