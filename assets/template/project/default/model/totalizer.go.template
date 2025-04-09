package model

import "time"

const (
	TotalizeCurrent = 1
	TotalizeHistory = 2
)

/*
Totalizer
Owner 统计名称
Type 类型(统计周期)
TimeStamp 统计时间
Value 统计结果
*/
type Totalizer struct {
	ID        int     `json:"id" orm:"id key auto"`
	Owner     string  `json:"owner" orm:"owner"`
	Type      int     `json:"type" orm:"type"`
	TimeStamp int64   `json:"timeStamp" orm:"timeStamp"`
	Value     float64 `json:"value" orm:"value"`
	Catalog   int     `json:"catalog" orm:"catalog"`
	Namespace string  `json:"namespace" orm:"namespace"`
}

func (s *Totalizer) Reset(timeStamp time.Time) {
	s.ID = 0
	s.Value = 0
	s.TimeStamp = timeStamp.UTC().Unix()
	s.Catalog = TotalizeCurrent
}

func (s *Totalizer) DuplicateHistory() *Totalizer {
	return &Totalizer{
		Owner:     s.Owner,
		Type:      s.Type,
		TimeStamp: s.TimeStamp,
		Value:     s.Value,
		Catalog:   TotalizeHistory,
		Namespace: s.Namespace,
	}
}

func NewTotalizer(owner string, typeVal int, namespace string) *Totalizer {
	return &Totalizer{Owner: owner, Type: typeVal, Catalog: TotalizeCurrent, Namespace: namespace}
}
