package user

import (
	"github.com/uptrace/bunrouter"
	"gorm.io/gorm"
	"log"
	xiangqinbackend "xiangqin-backend"
	"xiangqin-backend/pkg/middleware"
	"xiangqin-backend/utils"
)

type UserRouter struct{}

// NewRouter _cfg作用是候选人模块请求算法用的
func (userRouter *UserRouter) NewRouter(db *gorm.DB, _cfg *xiangqinbackend.Config, router *bunrouter.Router, jwt *utils.JWT) *bunrouter.Router {
	log.Println("userRouter register")
	userService := NewUserService(db, jwt)
	userApi := NewUserApi(userService)

	router.POST("/api/v1/login", userApi.Login)

	router.Use(middleware.HTTPJwt(jwt)).
		WithGroup("/api/v1", func(g *bunrouter.Group) {
			g.GET("/user/", userApi.GetUser)
			g.POST("/user/", userApi.CreateUser)
			g.PUT("/user/", userApi.UpdateUser)
			g.DELETE("/user/:id", userApi.DeleteUser)
			g.POST("/exit", userApi.Exit)
		})

	return router
}
