package company

import (
	"github.com/uptrace/bunrouter"
	"gorm.io/gorm"
	"log"
	xiangqin_backend "xiangqin-backend"
)

type CompanyRouter struct{}

// _cfg作用是候选人模块请求算法用的
func (userRouter *CompanyRouter) NewRouter(db *gorm.DB, _cfg *xiangqin_backend.Config, router *bunrouter.Router) *bunrouter.Router {
	log.Println("userRouter register")
	companyService := NewUserService(db)
	companyApi := NewUserApi(companyService)

	router.WithGroup("/api/v1", func(g *bunrouter.Group) {
		g.GET("/company/", companyApi.GetCompany)
		g.POST("/company/", companyApi.CreateCompany)
		g.PUT("/company/", companyApi.UpdateCompany)
		g.DELETE("/company/:code", companyApi.DeleteCompany)
	})

	return router
}
