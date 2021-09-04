package common

import (
	casCommon "github.com/muidea/magicCas/common"
	casModel "github.com/muidea/magicCas/model"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/{{.ProjectName}}/model"
)

const (
	FilterBill = "/bill/query/"
	QueryBill  = "/bill/query/:id"
	CreateBill = "/bill/create/"
	UpdateBill = "/bill/update/:id"
	DeleteBill = "/bill/delete/:id"

	FilterBillBillReport = "/bill/billreport/query/"
	QueryBillBillReport  = "/bill/billreport/query/:id"
	CreateBillBillReport = "/bill/billreport/create/"
	UpdateBillBillReport = "/bill/billreport/update/:id"
	DeleteBillBillReport = "/bill/billreport/delete/:id"

	FilterBillPayReward = "/bill/payreward/query/"
	QueryBillPayReward  = "/bill/payreward/query/:id"
	CreateBillPayReward = "/bill/payreward/create/"
	UpdateBillPayReward = "/bill/payreward/update/:id"
	DeleteBillPayReward = "/bill/payreward/delete/:id"

	FilterBillRewardPolicy = "/bill/rewardpolicy/query/"
	QueryBillRewardPolicy  = "/bill/rewardpolicy/query/:id"
	CreateBillRewardPolicy = "/bill/rewardpolicy/create/"
	UpdateBillRewardPolicy = "/bill/rewardpolicy/update/:id"
	DeleteBillRewardPolicy = "/bill/rewardpolicy/delete/:id"
)

const (
	Bill         = "bill"
	BillReport   = "billReport"
	PayReward    = "payReward"
	RewardPolicy = "rewardPolicy"
)

const BillModule = "/module/bill"

type BillView struct {
	ID         int                   `json:"id"`
	SN         string                `json:"sn"`
	Owner      *PartnerLite          `json:"owner"`
	Order      *OrderLite            `json:"order"`
	Reward     float64               `json:"reward"`
	Type       int                   `json:"type"`
	Level      int                   `json:"level"`
	Status     *model.Status         `json:"status"`
	Creater    *casCommon.EntityView `json:"creater"`
	UpdateTime int64                 `json:"updateTime"`
}

func (s *BillView) FromBill(ptr *model.Bill, entityPtr *casModel.Entity) {
	s.ID = ptr.ID
	s.SN = ptr.SN
	if ptr.Owner != nil {
		lite := &PartnerLite{}
		lite.FromPartner(ptr.Owner)
		s.Owner = lite
	}
	if ptr.Order != nil {
		lite := &OrderLite{}
		lite.FromOrder(ptr.Order)
		s.Order = lite
	}
	s.Reward = ptr.Reward
	s.Type = ptr.Type
	s.Level = ptr.Level
	s.Status = ptr.Status
	if entityPtr != nil {
		view := &casCommon.EntityView{}
		view.FromEntity(entityPtr)
		s.Creater = view
	}
	s.UpdateTime = ptr.UpdateTime
}

type BillLite struct {
	ID    int          `json:"id"`
	SN    string       `json:"sn"`
	Owner *PartnerLite `json:"owner"`
}

func (s *BillLite) FromBill(ptr *model.Bill) {
	s.ID = ptr.ID
	s.SN = ptr.SN
	if ptr.Owner != nil {
		lite := &PartnerLite{}
		lite.FromPartner(ptr.Owner)
		s.Owner = lite
	}
}

func (s *BillLite) ToBill() (ret *model.Bill) {
	ret = &model.Bill{
		ID:    s.ID,
		SN:    s.SN,
		Owner: s.Owner.ToPartner(),
	}
	return
}

type BillParam struct {
	Owner   *model.Partner `json:"owner"`
	Order   *model.Order   `json:"order"`
	Reward  float64        `json:"reward"`
	Type    int            `json:"type"`
	Level   int            `json:"level"`
	Status  *model.Status  `json:"status"`
	Creater int            `json:"creater"`
}

func (s *BillParam) ToBill() (ret *model.Bill) {
	ret = &model.Bill{
		Owner:   s.Owner,
		Order:   s.Order,
		Reward:  s.Reward,
		Type:    s.Type,
		Level:   s.Level,
		Creater: s.Creater,
	}
	if s.Status != nil {
		ret.Status = s.Status
	}

	return
}

type BillResult struct {
	commonDef.Result
	Bill *BillView `json:"bill"`
}

type BillListResult struct {
	commonDef.Result
	Bill []*BillLite `json:"bill"`
}

type BillStatisticResult struct {
	commonDef.Result
	Total   int64            `json:"total"`
	Summary []*TotalizerView `json:"summary"`
	Bill    []*BillView      `json:"bill"`
}

type BillReportView struct {
	ID         int                   `json:"id"`
	SN         string                `json:"sn"`
	Owner      *PartnerLite          `json:"owner"`
	Reward     float64               `json:"reward"`
	Available  float64               `json:"available"`
	Creater    *casCommon.EntityView `json:"creater"`
	UpdateTime int64                 `json:"updateTime"`
}

func (s *BillReportView) FromBillReport(ptr *model.BillReport, entityPtr *casModel.Entity) {
	s.ID = ptr.ID
	s.SN = ptr.SN
	if ptr.Owner != nil {
		lite := &PartnerLite{}
		lite.FromPartner(ptr.Owner)
		s.Owner = lite
	}
	s.Reward = ptr.Reward
	s.Available = ptr.Available
	if entityPtr != nil {
		view := &casCommon.EntityView{}
		view.FromEntity(entityPtr)
		s.Creater = view
	}
	s.UpdateTime = ptr.UpdateTime
}

type BillReportLite struct {
	ID        int          `json:"id"`
	SN        string       `json:"sn"`
	Owner     *PartnerLite `json:"owner"`
	Reward    float64      `json:"reward"`
	Available float64      `json:"available"`
}

func (s *BillReportLite) FromBillReport(ptr *model.BillReport) {
	s.ID = ptr.ID
	s.SN = ptr.SN
	if ptr.Owner != nil {
		lite := &PartnerLite{}
		lite.FromPartner(ptr.Owner)
		s.Owner = lite
	}
	s.Reward = ptr.Reward
	s.Available = ptr.Available
}

func (s *BillReportLite) ToBillReport() (ret *model.BillReport) {
	ret = &model.BillReport{
		ID:    s.ID,
		SN:    s.SN,
		Owner: s.Owner.ToPartner(),
	}
	return
}

type BillReportParam struct {
	SN        string         `json:"sn"`
	Owner     *model.Partner `json:"owner"`
	Reward    float64        `json:"reward"`
	Available float64        `json:"available"`
}

func (s *BillReportParam) ToBillReport() (ret *model.BillReport) {
	ret = &model.BillReport{
		SN:        s.SN,
		Owner:     s.Owner,
		Reward:    s.Reward,
		Available: s.Reward,
	}

	return
}

type BillReportResult struct {
	commonDef.Result
	BillReport *BillReportView `json:"billReport"`
}

type BillReportListResult struct {
	commonDef.Result
	BillReport []*BillReportLite `json:"billReport"`
}

type BillReportStatisticResult struct {
	commonDef.Result
	Total      int64             `json:"total"`
	Summary    []*TotalizerView  `json:"summary"`
	BillReport []*BillReportView `json:"billReport"`
}

type PayRewardView struct {
	ID         int                   `json:"id"`
	SN         string                `json:"sn"`
	Owner      *PartnerLite          `json:"owner"`
	Reward     float64               `json:"reward"`
	Comment    string                `json:"comment"`
	Creater    *casCommon.EntityView `json:"creater"`
	CreateTime int64                 `json:"createTime"`
}

func (s *PayRewardView) FromPayReward(ptr *model.PayReward, entityPtr *casModel.Entity) {
	s.ID = ptr.ID
	s.SN = ptr.SN
	if ptr.Owner != nil {
		lite := &PartnerLite{}
		lite.FromPartner(ptr.Owner)
		s.Owner = lite
	}
	s.Reward = ptr.Reward
	s.Comment = ptr.Comment
	if entityPtr != nil {
		view := &casCommon.EntityView{}
		view.FromEntity(entityPtr)
		s.Creater = view
	}
	s.CreateTime = ptr.CreateTime
}

type PayRewardLite struct {
	ID    int          `json:"id"`
	SN    string       `json:"sn"`
	Owner *PartnerLite `json:"owner"`
}

func (s *PayRewardLite) FromPayReward(ptr *model.PayReward) {
	s.ID = ptr.ID
	s.SN = ptr.SN
	if ptr.Owner != nil {
		lite := &PartnerLite{}
		lite.FromPartner(ptr.Owner)
		s.Owner = lite
	}
}

func (s *PayRewardLite) ToPayReward() (ret *model.PayReward) {
	ret = &model.PayReward{
		ID:    s.ID,
		SN:    s.SN,
		Owner: s.Owner.ToPartner(),
	}
	return
}

type PayRewardParam struct {
	Owner   *PartnerLite `json:"owner" validate:"required"`
	Reward  float64      `json:"reward" validate:"required"`
	Comment string       `json:"comment"`
}

func (s *PayRewardParam) ToPayReward() (ret *model.PayReward) {
	ret = &model.PayReward{
		Owner:   s.Owner.ToPartner(),
		Reward:  s.Reward,
		Comment: s.Comment,
	}
	return
}

type PayRewardResult struct {
	commonDef.Result
	PayReward *PayRewardView `json:"payReward"`
}

type PayRewardListResult struct {
	commonDef.Result
	PayReward []*PayRewardLite `json:"payReward"`
}

type PayRewardStatisticResult struct {
	commonDef.Result
	Total     int64            `json:"total"`
	Summary   []*TotalizerView `json:"summary"`
	PayReward []*PayRewardView `json:"payReward"`
}

type RewardPolicyView struct {
	ID          int                   `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	ValueItem   []model.ValueItem     `json:"item"`
	ValueScope  model.ValueScope      `json:"scope"`
	Status      *model.Status         `json:"status"`
	Creater     *casCommon.EntityView `json:"creater"`
	UpdateTime  int64                 `json:"updateTime"`
}

func (s *RewardPolicyView) FromRewardPolicy(ptr *model.RewardPolicy, entityPtr *casModel.Entity) {
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	s.ValueItem = ptr.ValueItem
	s.ValueScope = ptr.ValueScope
	s.Status = ptr.Status
	if entityPtr != nil {
		view := &casCommon.EntityView{}
		view.FromEntity(entityPtr)
		s.Creater = view
	}
	s.UpdateTime = ptr.UpdateTime
}

type RewardPolicyLite struct {
	ID     int           `json:"id"`
	Name   string        `json:"name"`
	Status *model.Status `json:"status"`
}

func (s *RewardPolicyLite) FromRewardPolicy(ptr *model.RewardPolicy) {
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Status = ptr.Status
}

func (s *RewardPolicyLite) ToRewardPolicy() (ret *model.RewardPolicy) {
	ret = &model.RewardPolicy{
		ID:     s.ID,
		Name:   s.Name,
		Status: s.Status,
	}
	return
}

type RewardPolicyParam struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	ValueItem   []model.ValueItem `json:"item"`
	ValueScope  model.ValueScope  `json:"scope"`
	Status      *model.Status     `json:"status"`
}

func (s *RewardPolicyParam) ToRewardPolicy() (ret *model.RewardPolicy) {
	ret = &model.RewardPolicy{
		Name:        s.Name,
		Description: s.Description,
		ValueItem:   s.ValueItem,
		ValueScope:  s.ValueScope,
	}
	if s.Status != nil {
		ret.Status = s.Status
	}
	return
}

type RewardPolicyResult struct {
	commonDef.Result
	RewardPolicy *RewardPolicyView `json:"rewardPolicy"`
}

type RewardPolicyListResult struct {
	commonDef.Result
	RewardPolicy []*RewardPolicyLite `json:"rewardPolicy"`
}

type RewardPolicyStatisticResult struct {
	commonDef.Result
	Total        int64               `json:"total"`
	Summary      []*TotalizerView    `json:"summary"`
	RewardPolicy []*RewardPolicyView `json:"rewardPolicy"`
	Status       []*model.Status     `json:"status"`
}
