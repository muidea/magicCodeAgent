package models

const (
	FLOW_STATUS_WAIT = iota
	FLOW_STATUS_FINISH
)

type MessagePair struct {
	Send *MessageInfo `json:"send"`
	Recv *MessageInfo `json:"recv"`
}

type Flow struct {
	Id        string         `json:"id"`
	Title     string         `json:"title"`
	Message   []*MessagePair `json:"message"`
	Status    int            `json:"status"`
	Parent    *Flow          `json:"parent"`
	Childrens []*Flow        `json:"childrens"`
}

func (s *Flow) AppendSend(m *MessageInfo) {
	if s.Message == nil {
		s.Message = make([]*MessagePair, 0)
	}
	if len(s.Message) > 0 {
		if s.Message[len(s.Message)-1].Recv == nil {
			s.Message[len(s.Message)-1].Send = m
			return
		}
	}

	s.Message = append(s.Message, &MessagePair{Send: m})
}

func (s *Flow) AppendRecv(m *MessageInfo) {
	if len(s.Message) == 0 {
		return
	}
	s.Message[len(s.Message)-1].Recv = m
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
