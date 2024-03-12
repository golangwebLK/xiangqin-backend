package pkg

import (
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"gorm.io/gorm"
	"os"
	xiangqin_backend "xiangqin-backend"
	"xiangqin-backend/pkg/candidate"
	"xiangqin-backend/pkg/company"
	"xiangqin-backend/pkg/user"
	"xiangqin-backend/utils"
)

type Router interface {
	NewRouter(db *gorm.DB, cfg *xiangqin_backend.Config, router *bunrouter.Router, jwt *utils.JWT) *bunrouter.Router
}

func NewRouter(db *gorm.DB, cfg *xiangqin_backend.Config) *bunrouter.Router {
	keyBytes, err := os.ReadFile("private_key.pem")
	if err != nil {
		panic(err)
	}
	jwt, err := utils.NewJWTFromKeyBytes(keyBytes)
	if err != nil {
		panic(err)
	}
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
		r.NewRouter(db, cfg, router, jwt)
	}
	return router
}
