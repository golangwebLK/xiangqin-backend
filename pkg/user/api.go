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

func (uApi *UserApi) GetUser(rw http.ResponseWriter, r bunrouter.Request) error {
	return nil
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
