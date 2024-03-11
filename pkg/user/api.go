package user

import (
	"github.com/uptrace/bunrouter"
	"net/http"
)

type UserApi struct {
	Svc *UserService
}

func NewUserApi(svc *UserService) *UserApi {
	return &UserApi{
		Svc: svc,
	}
}

func (uApi *UserApi) Login(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}
func (uApi *UserApi) Exit(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}

func (uApi *UserApi) GetMenu(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}

func (uApi *UserApi) GetUser(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}

type RequestUser struct {
	Name        string `json:"name"`
	Birth       string `json:"birth"`
	Telephone   string `json:"telephone"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	IsUser      bool   `json:"isUser"`
	Role        string `json:"role"`
	CompanyCode string `json:"companyCode"`
	Remarks     string `json:"remarks"`
}

func (uApi *UserApi) CreateUser(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}
func (uApi *UserApi) UpdateUser(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}
func (uApi *UserApi) DeleteUser(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
}
