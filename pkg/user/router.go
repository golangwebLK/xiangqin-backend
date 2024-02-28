package user

import (
	"github.com/uptrace/bunrouter"
	"gorm.io/gorm"
	"log"
)

type UserRouter struct{}

func (userRouter *UserRouter) NewRouter(db *gorm.DB, router *bunrouter.Router) *bunrouter.Router {
	log.Println("userRouter register")
	userService := NewUserService(db)
	userApi := NewUserApi(userService)

	router.GET("/user/", userApi.GetUser)

	return router
}
