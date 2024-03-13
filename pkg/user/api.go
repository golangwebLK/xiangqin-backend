package user

import (
	"encoding/json"
	"github.com/uptrace/bunrouter"
	"net/http"
	"net/url"
	"strconv"
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
	http.SetCookie(rw, &http.Cookie{
		Name:    "xq-session",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
	})
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "退出成功",
		Data:    nil,
	})
}

func (uApi *UserApi) GetUser(rw http.ResponseWriter, r bunrouter.Request) error {
	ctx := r.Request.Context()
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "字段错误",
			Data:    err,
		})
	}
	pageInt, _ := strconv.Atoi(queryValues.Get("page"))
	pageSizeInt, _ := strconv.Atoi(queryValues.Get("pageSize"))
	name := queryValues.Get("name")
	users, err := uApi.Svc.GetUser(ctx, pageInt, pageSizeInt, name)
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "查询错误",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "查询成功",
		Data:    users,
	})
}

type RequestUser struct {
	Name      string `json:"name"`
	Birth     string `json:"birth"`
	Telephone string `json:"telephone"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	IsUser    bool   `json:"isUser"`
	Role      string `json:"role"`
	Remarks   string `json:"remarks"`
}

func (uApi *UserApi) CreateUser(rw http.ResponseWriter, r bunrouter.Request) error {
	ctx := r.Request.Context()
	var rUser RequestUser
	if err := json.NewDecoder(r.Body).Decode(&rUser); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "请求字段错误",
			Data:    err,
		})
	}
	if err := uApi.Svc.CreateUser(ctx, rUser); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "用户保存失败",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "保存成功",
		Data:    nil,
	})
}
func (uApi *UserApi) UpdateUser(rw http.ResponseWriter, r bunrouter.Request) error {
	ctx := r.Request.Context()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "请求字段错误",
			Data:    err,
		})
	}
	if err := uApi.Svc.UpdateUser(ctx, user); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "更新失败",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "更新成功",
		Data:    nil,
	})
}
func (uApi *UserApi) DeleteUser(rw http.ResponseWriter, r bunrouter.Request) error {
	ctx := r.Request.Context()
	params := r.Params()
	id, _ := params.Int64("id")
	if err := uApi.Svc.DeleteUser(ctx, int(id)); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "删除失败",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "删除成功",
		Data:    nil,
	})
}

func (uApi *UserApi) UpdatePassword(rw http.ResponseWriter, r bunrouter.Request) error {
	ctx := r.Request.Context()
	var p LoginReq
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "字段序列化失败",
			Data:    err,
		})
	}
	if err := uApi.Svc.UpdatePassword(ctx, p); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "密码修改失败",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "密码修改成功",
		Data:    nil,
	})
}
