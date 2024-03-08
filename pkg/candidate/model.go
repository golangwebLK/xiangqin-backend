package candidate

import (
	"encoding/json"

	"gorm.io/gorm"
)

// 上线前重点检查，有些字段要有gorm:union;not null
type PersonalInfo struct {
	gorm.Model
	PersonCode                string          `gorm:"type:varchar(255);unique;not null" json:"person_code"`
	RealName                  string          `gorm:"type:varchar(255);not null" json:"real_name"`
	Gender                    string          `gorm:"type:varchar(255);not null" json:"gender"`
	BirthYear                 int             `gorm:"not null" json:"birth_year"`
	Telephone                 string          `gorm:"type:varchar(255);not null" json:"telephone"`
	WeChat                    string          `gorm:"type:varchar(255);not null" json:"wechat"`
	Work                      int             `gorm:"not null" json:"work"`
	School                    string          `gorm:"type:varchar(255);not null" json:"school"`
	Qualification             string          `gorm:"type:varchar(255);not null" json:"qualification"`
	CurrentPlace              int             `gorm:"not null" json:"current_place"`
	AncestralHome             int             `gorm:"not null" json:"ancestral_home"`
	Economic                  json.RawMessage `gorm:"type:json;not null" json:"economic"`
	Hobbies                   string          `gorm:"type:varchar(255);not null" json:"hobbies"`
	Height                    float64         `gorm:"not null" json:"height"`
	Weight                    float64         `gorm:"not null" json:"weight"`
	OriginalFamilyComposition string          `gorm:"type:varchar(255);not null" json:"original_family_composition"`
	ParentsSituation          string          `gorm:"type:varchar(255);not null" json:"parents_situation"`
	CompanyCode               string          `json:"companyCode" gorm:"type:varchar(255);"`
	Remarks                   string          `gorm:"type:varchar(255);not null" json:"remarks"`
}

type EconomicInfo struct {
	Savings    float64 `json:"savings"`
	House      string  `json:"house"`
	HouseMoney float64 `json:"house_money"`
	Car        string  `json:"car"`
	CarMoney   float64 `json:"car_money"`
}

type CandidateReq struct {
	PersonCode    string  `json:"person_code"`
	BirthYear     int     `json:"birth_year"`     // 实际年龄
	Work          []int   `json:"work"`           // 按照包含关系，填入编号
	Qualification int     `json:"qualification"`  // 学历编号1-7，
	CurrentPlace  []int   `json:"current_place"`  // 按照包含关系，填入编号
	AncestralHome []int   `json:"ancestral_home"` // 按照包含关系，填入编号
	Economic      float64 `json:"economic"`       // 实际财富
	Height        float64 `json:"height"`         // 实际身高
	Weight        float64 `json:"weight"`         // 实际体重
	Score         float64 `json:"score"`
}

type PersonalLike struct {
	gorm.Model
	PersonCode                string          `gorm:"type:varchar(255);unique;not null" json:"person_code"`
	BirthYear                 int             `gorm:"not null" json:"birth_year"`
	Work                      int             `gorm:"not null" json:"work"`
	Qualification             string          `gorm:"type:varchar(255);not null" json:"qualification"`
	CurrentPlace              int             `gorm:"not null" json:"current_place"`
	AncestralHome             int             `gorm:"not null" json:"ancestral_home"`
	Economic                  json.RawMessage `gorm:"type:json;not null" json:"economic"`
	Hobbies                   string          `gorm:"type:varchar(255);not null" json:"hobbies"`
	Height                    float64         `gorm:"not null" json:"height"`
	Weight                    float64         `gorm:"not null" json:"weight"`
	OriginalFamilyComposition string          `gorm:"type:varchar(255);not null" json:"original_family_composition"`
	ParentsSituation          string          `gorm:"type:varchar(255);not null" json:"parents_situation"`
	Remarks                   string          `gorm:"type:varchar(255);not null" json:"remarks"`
}
