package candidate

import (
	"gorm.io/gorm"
	xiangqin_backend "xiangqin-backend"
)

type CandidateService struct {
	DB          *gorm.DB
	MatchServer string
}

func NewCandidateService(db *gorm.DB, cfg *xiangqin_backend.Config) *CandidateService {
	return &CandidateService{
		DB:          db,
		MatchServer: cfg.MatchService.Server,
	}
}

func (candidateService *CandidateService) SavePersonalInfo(
	personalInfo PersonalInfo) error {
	if err := candidateService.DB.Create(&personalInfo).Error; err != nil {
		return err
	}
	return nil
}

type Candidate struct {
	BirthYear     int8    `json:"birth_year"`     // 实际年龄
	Work          []int8  `json:"work"`           // 按照包含关系，填入编号
	Qualification int8    `json:"qualification"`  // 学历编号1-7，
	CurrentPlace  []int8  `json:"current_place"`  // 按照包含关系，填入编号
	AncestralHome []int8  `json:"ancestral_home"` // 按照包含关系，填入编号
	Economic      float64 `json:"economic"`       // 实际财富
	Height        float64 `json:"height"`         // 实际身高
	Weight        float64 `json:"weight"`         // 实际体重
	Score         float64 `json:"score"`
}

type RequestData struct {
	Candidate  Candidate          `json:"candidate"`
	Candidates []Candidate        `json:"candidates"`
	Attributes map[string]float64 `json:"attributes"`
}

func (candidateService *CandidateService) MatchCandidate(
	personalInfo PersonalInfo,
	attributes_map map[string]float64) error {

	return nil
}
