package user

import (
	"github.com/uptrace/bunrouter"
	"gorm.io/gorm"
	"log"
	xiangqin_backend "xiangqin-backend"
)

type UserRouter struct{}

// _cfg作用是候选人模块请求算法用的
func (userRouter *UserRouter) NewRouter(db *gorm.DB, _cfg *xiangqin_backend.Config, router *bunrouter.Router) *bunrouter.Router {
	log.Println("userRouter register")
	userService := NewUserService(db)
	userApi := NewUserApi(userService)

	router.POST("/api/v1/login", userApi.Login)
	router.POST("/api/v1/exit", userApi.Exit)

	router.WithGroup("/api/v1", func(g *bunrouter.Group) {
		g.GET("/menu", userApi.GetMenu)
		g.GET("/user/", userApi.GetUser)
		g.POST("/user/", userApi.CreateUser)
		g.PUT("/user/", userApi.UpdateUser)
		g.DELETE("/user/:id", userApi.DeleteUser)
	})

	return router
}
