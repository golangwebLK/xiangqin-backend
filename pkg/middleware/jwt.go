package middleware

import (
	"context"
	"github.com/uptrace/bunrouter"
	"net/http"
	"strconv"
	"strings"
	"xiangqin-backend/utils"
)

type Msg struct {
	companyCode string
	ID          int
}

func HTTPJwt(jwt *utils.JWT) bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			cookie, err := req.Request.Cookie("xq-session")
			if err != nil {
				return bunrouter.JSON(w, utils.ResponseData{
					Status:  http.StatusBadRequest,
					Message: "没有找到cookie",
					Data:    err,
				})
			}
			claim, err := jwt.Verify(cookie.Value)
			if err != nil {
				return bunrouter.JSON(w, utils.ResponseData{
					Status:  http.StatusBadRequest,
					Message: "cookie错误",
					Data:    err,
				})
			}
			ids := strings.Split(claim.Subject, "@")
			id, err := strconv.Atoi(ids[1])
			if err != nil {
				return err
			}
			msg := Msg{
				companyCode: ids[0],
				ID:          id,
			}
			ctx := context.WithValue(context.TODO(), "msg", msg)
			req.Request = req.Request.WithContext(ctx)
			return next(w, req)
		}
	}
}
