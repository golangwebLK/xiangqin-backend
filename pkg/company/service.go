package company

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
	"xiangqin-backend/pkg/user"
	"xiangqin-backend/utils"
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

func (csvc *CompanyService) GetCompany(pageInt, pageSizeInt int, name, startTime, endTime string) (utils.PagingResp, error) {
	var companys []Company
	query := csvc.DB.Model(&Company{})
	if name != "" {
		query = query.Where("company_name=?", name)
	}
	if startTime != "" {
		startT, err := time.Parse("2006-01-02 15:04:05", startTime)
		if err != nil {
			return utils.PagingResp{}, err
		}
		query = query.Where("created_at>=?", startT)
	}
	if endTime != "" {
		endT, err := time.Parse("2006-01-02 15:04:05", endTime)
		if err != nil {
			return utils.PagingResp{}, err
		}
		query = query.Where("created_at<=?", endT)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.PagingResp{}, err
	}
	offset, err := calculateOffset(pageInt, pageSizeInt, total)
	if err != nil {
		return utils.PagingResp{}, err
	}
	if err = query.
		Offset(offset).
		Limit(pageSizeInt).
		Find(&companys).
		Error; err != nil {
		return utils.PagingResp{}, err
	}
	paging := utils.PagingResp{
		Page:     pageInt,
		PageSize: pageSizeInt,
		Total:    total,
		Data:     companys,
	}
	return paging, nil
}

func calculateOffset(page, pageSize int, totalRecords int64) (int, error) {
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pageSize
	if offset > int(totalRecords) {
		return 0, errors.New("page out of range")
	}

	return offset, nil
}
