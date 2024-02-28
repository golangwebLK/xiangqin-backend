package pkg

import (
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"gorm.io/gorm"
	"xiangqin-backend/pkg/user"
)

type Router interface {
	NewRouter(db *gorm.DB, router *bunrouter.Router) *bunrouter.Router
}

func NewRouter(db *gorm.DB) *bunrouter.Router {
	router := bunrouter.New(
		bunrouter.Use(
			reqlog.NewMiddleware(),
		))
	routers := []Router{
		&user.UserRouter{},
	}
	for _, r := range routers {
		r.NewRouter(db, router)
	}
	return router
}
