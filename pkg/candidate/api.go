package candidate

import (
	"encoding/json"
	"github.com/uptrace/bunrouter"
	"net/http"
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

type Personal struct {
	RealName                  string       `json:"real_name"`
	BirthYear                 int          `json:"birth_year"`
	Telephone                 string       `json:"telephone"`
	WeChat                    string       `json:"we_chat"`
	Work                      int          `json:"work"`
	School                    string       `json:"school"`
	Qualification             string       `json:"qualification"`
	CurrentPlace              int          `json:"current_place"`
	AncestralHome             int          `json:"ancestral_home"`
	Economic                  EconomicInfo `json:"economic"`
	Hobbies                   string       `json:"hobbies"`
	Height                    int          `json:"height"`
	Weight                    int          `json:"weight"`
	OriginalFamilyComposition string       `json:"original_family_composition"`
	ParentsSituation          string       `json:"parents_situation"`
	Remarks                   string       `json:"remarks"`
	Gender                    string       `json:"gender"`
}

func (cApi *CandidateApi) CreateCandidate(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	var personal Personal
	if err := json.NewDecoder(r.Body).Decode(&personal); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "请求参数错误",
			Data:    err,
		})
	}
	dbPersonal := PersonalInfo{
		RealName:                  personal.RealName,
		BirthYear:                 personal.BirthYear,
		Telephone:                 personal.Telephone,
		WeChat:                    personal.WeChat,
		Work:                      personal.Work,
		School:                    personal.School,
		Qualification:             personal.Qualification,
		CurrentPlace:              personal.CurrentPlace,
		AncestralHome:             personal.AncestralHome,
		Economic:                  personal.Economic,
		Hobbies:                   personal.Hobbies,
		Height:                    personal.Height,
		Weight:                    personal.Weight,
		OriginalFamilyComposition: personal.OriginalFamilyComposition,
		ParentsSituation:          personal.ParentsSituation,
		Remarks:                   personal.Remarks,
		Gender:                    personal.Gender,
	}
	err := cApi.Svc.SavePersonalInfo(dbPersonal)
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "数据库保存错误",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "保存成功",
		Data:    nil,
	})
}
