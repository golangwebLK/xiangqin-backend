package company

import "gorm.io/gorm"

type CompanyService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *CompanyService {
	return &CompanyService{
		DB: db,
	}
}
