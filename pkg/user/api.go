package user

import (
	"encoding/json"
	"github.com/uptrace/bunrouter"
	"net/http"
	"time"
	"xiangqin-backend/utils"
)

type UserApi struct {
	Svc *UserService
}

func NewUserApi(svc *UserService) *UserApi {
	return &UserApi{
		Svc: svc,
	}
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (uApi *UserApi) Login(rw http.ResponseWriter, r bunrouter.Request) error {
	var loginReq LoginReq
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "字段序列化错误",
			Data:    err,
		})
	}
	user, err := uApi.Svc.ComparePassword(loginReq)
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "密码验证错误",
			Data:    err,
		})
	}
	exp := time.Now().AddDate(0, 0, 1)
	companyCodeAndID := uApi.Svc.StrConcatenation(user.ID, user.CompanyCode)
	tokenStr, err := uApi.Svc.SignByID(companyCodeAndID, exp)
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "token生成错误",
			Data:    err,
		})
	}
	http.SetCookie(rw, &http.Cookie{
		Name:    "xq-session",
		Value:   tokenStr,
		Path:    "/",
		Expires: exp,
	})
	contents, err := uApi.Svc.GetContent(user)
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "目录生成失败",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "登陆成功",
		Data:    contents,
	})
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
