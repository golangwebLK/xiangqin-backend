package candidate

import "gorm.io/gorm"

type PersonalInfo struct {
	gorm.Model
	RealName                  string   `gorm:"type:varchar(255);not null" json:"real_name"`
	Gender                    string   `gorm:"type:varchar(255);not null" json:"gender"`
	BirthYear                 int      `gorm:"not null" json:"birth_year"`
	Telephone                 int      `gorm:"unique;not null" json:"telephone"`
	WeChat                    string   `gorm:"type:varchar(255);unique;not null" json:"wechat"`
	Work                      int      `gorm:"not null" json:"work"`
	School                    string   `gorm:"type:varchar(255);not null" json:"school"`
	Qualification             string   `gorm:"type:varchar(255);not null" json:"qualification"`
	CurrentPlace              int      `gorm:"not null" json:"current_place"`
	AncestralHome             int      `gorm:"not null" json:"ancestral_home"`
	Economic                  []string `gorm:"type:json;not null" json:"economic"`
	Hobbies                   string   `gorm:"type:varchar(255);not null" json:"hobbies"`
	Height                    int      `gorm:"not null" json:"height"`
	Weight                    int      `gorm:"not null" json:"weight"`
	OriginalFamilyComposition []string `gorm:"type:json;not null" json:"original_family_composition"`
	ParentsSituation          []string `gorm:"type:json;not null" json:"parents_situation"`
	Remarks                   string   `gorm:"type:varchar(255);not null" json:"remarks"`
}
