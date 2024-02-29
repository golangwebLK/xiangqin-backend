package candidate

import (
	"github.com/uptrace/bunrouter"
	"gorm.io/gorm"
	"log"
)

type CandidateRouter struct{}

func (candidateRouter *CandidateRouter) NewRouter(db *gorm.DB, router *bunrouter.Router) *bunrouter.Router {
	log.Println("userRouter register")
	candidateService := NewCandidateService(db)
	candidateApi := NewCandidateApi(candidateService)

	router.WithGroup("/api/v1", func(g *bunrouter.Group) {
		g.POST("/candidate/", candidateApi.CreateCandidate)
	})

	return router
}
