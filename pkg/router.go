package pkg

import (
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"gorm.io/gorm"
	xiangqin_backend "xiangqin-backend"
	"xiangqin-backend/pkg/candidate"
	"xiangqin-backend/pkg/company"
	"xiangqin-backend/pkg/user"
)

type Router interface {
	NewRouter(db *gorm.DB, cfg *xiangqin_backend.Config, router *bunrouter.Router) *bunrouter.Router
}

func NewRouter(db *gorm.DB, cfg *xiangqin_backend.Config) *bunrouter.Router {
	router := bunrouter.New(
		bunrouter.Use(
			reqlog.NewMiddleware(),
		))
	routers := []Router{
		&user.UserRouter{},
		&candidate.CandidateRouter{},
		&company.CompanyRouter{},
	}
	for _, r := range routers {
		r.NewRouter(db, cfg, router)
	}
	return router
}
