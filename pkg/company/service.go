package company

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

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
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(createCompanyRequestData.RootPassword), 12)
	rootUser := user.User{
		Name:        "默认姓名",
		Username:    createCompanyRequestData.RootUsername,
		Password:    string(passwordHash),
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
		query = query.Where("company_name like ?", "%"+name+"%")
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
	offset, err := CalculateOffset(pageInt, pageSizeInt, total)
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

	codes := make([]string, 0, len(companys))
	for _, company := range companys {
		codes = append(codes, company.Code)
	}
	var users []user.User
	if err = csvc.DB.Where("company_code in (?)", codes).Find(&users).Error; err != nil {
		return utils.PagingResp{}, err
	}
	var companyAndUsers []CompanyAndUser
	codeUsersMap := make(map[string][]user.User, len(users)/10)
	for _, u := range users {
		if _, exist := codeUsersMap[u.CompanyCode]; exist {
			codeUsersMap[u.CompanyCode] = append(codeUsersMap[u.CompanyCode], u)
		} else {
			codeUsersMap[u.CompanyCode] = []user.User{}
			codeUsersMap[u.CompanyCode] = append(codeUsersMap[u.CompanyCode], u)
		}
	}
	for _, company := range companys {
		companyAndUser := CompanyAndUser{
			Company: company,
		}
		if _, exist := codeUsersMap[company.Code]; exist {
			companyAndUser.User = codeUsersMap[company.Code]
		}
		companyAndUsers = append(companyAndUsers, companyAndUser)
	}
	paging := utils.PagingResp{
		Page:     pageInt,
		PageSize: pageSizeInt,
		Total:    total,
		Data:     companyAndUsers,
	}
	return paging, nil
}

type CompanyAndUser struct {
	Company Company     `json:"company"`
	User    []user.User `json:"users"`
}

func CalculateOffset(page, pageSize int, totalRecords int64) (int, error) {
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pageSize
	if offset > int(totalRecords) {
		return 0, errors.New("page out of range")
	}

	return offset, nil
}

func (csvc *CompanyService) UpdateCompany(company Company) error {
	if err := csvc.DB.Updates(company).Where("id=?", company.ID).Error; err != nil {
		return err
	}
	return nil
}

func (csvc *CompanyService) DeleteCompany(code string) error {
	tx := csvc.DB.Begin()
	if err := csvc.DB.Where("code=?", code).Delete(&Company{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := csvc.DB.Where("company_code=?", code).Delete(&user.User{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
