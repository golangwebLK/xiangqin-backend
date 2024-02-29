package candidate

import (
	"encoding/json"
	"github.com/uptrace/bunrouter"
	"net/http"
	"strings"
	"xiangqin-backend/utils"
)

type CandidateApi struct {
	Svc *CandidateService
}

func NewCandidateApi(svc *CandidateService) *CandidateApi {
	return &CandidateApi{
		Svc: svc,
	}
}

type Candidate struct {
	RealName                  string   `json:"realName"`
	BirthYear                 int      `json:"birthYear"`
	Telephone                 int      `json:"telephone"`
	WeChat                    string   `json:"weChat"`
	Work                      int      `json:"work"`
	School                    string   `json:"school"`
	Qualification             string   `json:"qualification"`
	CurrentPlace              int      `json:"currentPlace"`
	AncestralHome             int      `json:"ancestralHome"`
	Economic                  []string `json:"economic"`
	Hobbies                   string   `json:"hobbies"`
	Height                    int      `json:"height"`
	Weight                    int      `json:"weight"`
	OriginalFamilyComposition string   `json:"originalFamilyComposition"`
	ParentsSituation          string   `json:"parentsSituation"`
	Remarks                   string   `json:"remarks"`
	Gender                    string   `json:"gender"`
}

func (cApi *CandidateApi) CreateCandidate(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	var candidate Candidate
	if err := json.NewDecoder(r.Body).Decode(&candidate); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "请求参数错误",
			Data:    err,
		})
	}
	dbCandidate := &PersonalInfo{
		RealName:                  candidate.RealName,
		BirthYear:                 candidate.BirthYear,
		Telephone:                 candidate.Telephone,
		WeChat:                    candidate.WeChat,
		Work:                      candidate.Work,
		School:                    candidate.School,
		Qualification:             candidate.Qualification,
		CurrentPlace:              candidate.CurrentPlace,
		AncestralHome:             candidate.AncestralHome,
		Economic:                  candidate.Economic,
		Hobbies:                   candidate.Hobbies,
		Height:                    candidate.Height,
		Weight:                    candidate.Weight,
		OriginalFamilyComposition: strings.Split(candidate.OriginalFamilyComposition, ","),
		ParentsSituation:          strings.Split(candidate.ParentsSituation, ","),
		Remarks:                   candidate.Remarks,
		Gender:                    candidate.Gender,
	}

	return nil
}
