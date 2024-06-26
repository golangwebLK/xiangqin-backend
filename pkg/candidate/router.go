package candidate

import (
	"github.com/uptrace/bunrouter"
	"gorm.io/gorm"
	"log"
	xiangqin_backend "xiangqin-backend"
	"xiangqin-backend/pkg/middleware"
	"xiangqin-backend/utils"
)

type CandidateRouter struct{}

func (candidateRouter *CandidateRouter) NewRouter(db *gorm.DB, cfg *xiangqin_backend.Config, router *bunrouter.Router, jwt *utils.JWT) *bunrouter.Router {
	log.Println("candidateRouter register")
	candidateService := NewCandidateService(db, cfg)
	candidateApi := NewCandidateApi(candidateService, cfg)

	router.Use(middleware.HTTPJwt(jwt)).
		WithGroup("/api/v1", func(g *bunrouter.Group) {
			g.POST("/candidate/", candidateApi.CreateCandidate)
			g.POST("/match/", candidateApi.GetMatch)
			g.GET("/candidate/", candidateApi.GetPersonalInfo)
			g.GET("/candidate/:id", candidateApi.GetPersonalInfoByID)
			g.PUT("/candidate/", candidateApi.UpdatePersonalInfo)
			g.PUT("/candidate_like/", candidateApi.UpdatePersonalLike)
			g.DELETE("/candidate/:code", candidateApi.DeletePersonalInfo)
			g.POST("/uploadImage", candidateApi.UploadImage)
			g.GET("/speech/:name", candidateApi.Speech())
		})

	return router
}
