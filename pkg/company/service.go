package company

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"xiangqin-backend/pkg/user"
)

type CompanyService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *CompanyService {
	return &CompanyService{
		DB: db,
	}
}

func (csvc *CompanyService) CreateCompanyAndUser(createCompanyRequestData CreateCompanyRequestData) error {
	//保存公司信息,
	code := uuid.NewString()
	tx := csvc.DB.Begin()
	company := Company{
		CompanyName:      createCompanyRequestData.CompanyName,
		ContactPerson:    createCompanyRequestData.ContactPerson,
		ContactTelephone: createCompanyRequestData.ContactTelephone,
		CompanyTelephone: createCompanyRequestData.CompanyTelephone,
		Address:          createCompanyRequestData.Address,
		Code:             code,
		IsUser:           createCompanyRequestData.IsUser,
		Remarks:          createCompanyRequestData.Remarks,
	}
	if err := csvc.DB.Create(&company).Error; err != nil {
		tx.Rollback()
		return err
	}
	//根据公司信息创建企业根用户
	rootUser := user.User{
		Name:        "默认姓名",
		Username:    createCompanyRequestData.RootUsername,
		Password:    createCompanyRequestData.RootPassword,
		IsUser:      true,
		Role:        "Manager",
		CompanyCode: code,
	}
	if err := csvc.DB.Create(&rootUser).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
