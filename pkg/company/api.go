package company

import (
	"encoding/json"
	"github.com/uptrace/bunrouter"
	"net/http"
	"xiangqin-backend/utils"
)

type CompanyApi struct {
	Svc *CompanyService
}

func NewUserApi(svc *CompanyService) *CompanyApi {
	return &CompanyApi{
		Svc: svc,
	}
}

func (cApi *CompanyApi) GetCompany(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}

type CreateCompanyRequestData struct {
	CompanyName      string `json:"companyName"`
	ContactPerson    string `json:"contactPerson"`
	ContactTelephone string `json:"contactTelephone"`
	CompanyTelephone string `json:"companyTelephone"`
	Address          string `json:"address"`
	IsUser           bool   `json:"isUser"`
	RootUsername     string `json:"rootUsername"`
	RootPassword     string `json:"rootPassword"`
	Remarks          string `json:"remarks"`
}

func (cApi *CompanyApi) CreateCompany(rw http.ResponseWriter, r bunrouter.Request) error {
	var reqData CreateCompanyRequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "请求参数错误",
			Data:    err,
		})
	}
	if err := cApi.Svc.CreateCompanyAndUser(reqData); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "数据保存错误",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "数据保存成功",
		Data:    nil,
	})
}
func (cApi *CompanyApi) UpdateCompany(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}
func (cApi *CompanyApi) DeleteCompany(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}
