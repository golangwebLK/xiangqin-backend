package candidate

import (
	"github.com/uptrace/bunrouter"
	"gorm.io/gorm"
	"log"
	xiangqin_backend "xiangqin-backend"
)

type CandidateRouter struct{}

func (candidateRouter *CandidateRouter) NewRouter(db *gorm.DB, cfg *xiangqin_backend.Config, router *bunrouter.Router) *bunrouter.Router {
	log.Println("userRouter register")
	candidateService := NewCandidateService(db, cfg)
	candidateApi := NewCandidateApi(candidateService)

	router.WithGroup("/api/v1", func(g *bunrouter.Group) {
		g.POST("/candidate/", candidateApi.CreateCandidate)
		g.POST("/match/", candidateApi.GetMatch)
	})

	return router
}
