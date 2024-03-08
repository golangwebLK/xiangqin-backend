package company

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	CompanyName      string `json:"companyName" gorm:"type:varchar(255);not null"`
	ContactPerson    string `json:"contactPerson" gorm:"type:varchar(255);not null"`
	ContactTelephone string `json:"contactTelephone" gorm:"type:varchar(255);not null"`
	CompanyTelephone string `json:"companyTelephone" gorm:"type:varchar(255);not null"`
	Address          string `json:"address" gorm:"type:varchar(255);not null"`
	Code             string `json:"code" gorm:"type:varchar(255);union;not null"`
	IsUser           bool   `json:"isUser" gorm:"not null"`
	Remarks          string `json:"remarks" gorm:"type:varchar(255);"`
}
