package models

const (
	FLOW_STATUS_WAIT = iota
	FLOW_STATUS_FINISH
)

type Flow struct {
	Id        string       `json:"id"`
	Title     string       `json:"title"`
	Message   *MessageInfo `json:"message"`
	Status    int          `json:"status"`
	Parent    *Flow        `json:"parent"`
	Childrens []*Flow      `json:"childrens"`
}

func (s *Flow) GetChildrens() []*Flow {
	return s.Childrens
}

func (s *Flow) GetParent() *Flow {
	return s.Parent
}

func (s *Flow) IsFinish() bool {
	return s.Status == FLOW_STATUS_FINISH
}
