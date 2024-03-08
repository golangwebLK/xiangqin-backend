package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Birth       string `json:"birth" gorm:"type:varchar(255);"`
	Telephone   string `json:"telephone" gorm:"type:varchar(255);"`
	Username    string `json:"username" gorm:"type:varchar(255);not null"`
	Password    string `json:"password" gorm:"type:varchar(255);not null"`
	IsUser      bool   `json:"isUser" gorm:"not null"`
	CompanyCode string `json:"companyCode" gorm:"type:varchar(255);not null"`
	Remarks     string `json:"remarks" gorm:"type:varchar(255);"`
}

type Permission struct {
	gorm.Model
	UserID    int `json:"userID" gorm:"not null"`
	ContentID int `json:"contentID" gorm:"not null"`
}

type Content struct {
	gorm.Model
	Name       string `json:"name" gorm:"type:varchar(255);not null"`
	Logo       string `json:"logo" gorm:"type:varchar(255);not null"`
	Code       string `json:"code" gorm:"type:varchar(255);not null"`
	ParentCode string `json:"parentCode" gorm:"type:varchar(255);"`
	Target     string `json:"target" gorm:"type:varchar(255);"`
}
