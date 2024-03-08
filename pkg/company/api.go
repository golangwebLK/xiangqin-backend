package company

import (
	"github.com/uptrace/bunrouter"
	"net/http"
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
func (cApi *CompanyApi) CreateCompany(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}
func (cApi *CompanyApi) UpdateCompany(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}
func (cApi *CompanyApi) DeleteCompany(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}
