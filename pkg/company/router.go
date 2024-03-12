package company

import (
	"github.com/uptrace/bunrouter"
	"gorm.io/gorm"
	"log"
	xiangqin_backend "xiangqin-backend"
	"xiangqin-backend/pkg/middleware"
	"xiangqin-backend/utils"
)

type CompanyRouter struct{}

// NewRouter _cfg作用是候选人模块请求算法用的
func (userRouter *CompanyRouter) NewRouter(db *gorm.DB, _cfg *xiangqin_backend.Config, router *bunrouter.Router, jwt *utils.JWT) *bunrouter.Router {
	log.Println("userRouter register")
	companyService := NewUserService(db)
	companyApi := NewUserApi(companyService)

	router.Use(middleware.HTTPJwt(jwt)).
		WithGroup("/api/v1", func(g *bunrouter.Group) {
			g.GET("/company/", companyApi.GetCompany)
			g.POST("/company/", companyApi.CreateCompany)
			g.PUT("/company/", companyApi.UpdateCompany)
			g.DELETE("/company/:code", companyApi.DeleteCompany)
		})

	return router
}
