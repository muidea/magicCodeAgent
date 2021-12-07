package model

const (
	TotalizeRealtime = 1
	TotalizeWeek     = 2
	TotalizeMonth    = 3
)

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

func (s *Totalizer) Reset() {
	s.ID = 0
	s.Value = 0
	s.Catalog = TotalizeCurrent
}

func NewTotalizer(owner string, typeVal int, namespace string) *Totalizer {
	if typeVal < TotalizeRealtime || typeVal > TotalizeMonth {
		return nil
	}

	return &Totalizer{Owner: owner, Type: typeVal, Catalog: TotalizeCurrent, Namespace: namespace}
}
