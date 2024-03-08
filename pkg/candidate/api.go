package candidate

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/uptrace/bunrouter"

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
	Hobbies                   []string     `json:"hobbies"`
	Height                    float64      `json:"height"`
	Weight                    float64      `json:"weight"`
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
	economic, err := json.Marshal(personal.Economic)
	if err != nil {
		return err
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
		Economic:                  economic,
		Hobbies:                   fmt.Sprint(personal.Hobbies),
		Height:                    personal.Height,
		Weight:                    personal.Weight,
		OriginalFamilyComposition: personal.OriginalFamilyComposition,
		ParentsSituation:          personal.ParentsSituation,
		Remarks:                   personal.Remarks,
		Gender:                    personal.Gender,
	}
	err = cApi.Svc.SavePersonalInfo(dbPersonal)
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

type RequestMatch struct {
	Personal      Personal           `json:"personal"`
	AttributesMap map[string]float64 `json:"attributes_map"`
}

func (cApi *CandidateApi) GetMatch(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	var req RequestMatch
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "请求参数错误",
			Data:    err,
		})
	}
	economic, err := json.Marshal(req.Personal.Economic)
	if err != nil {
		return err
	}
	dbPersonal := PersonalInfo{
		RealName:                  req.Personal.RealName,
		BirthYear:                 req.Personal.BirthYear,
		Telephone:                 req.Personal.Telephone,
		WeChat:                    req.Personal.WeChat,
		Work:                      req.Personal.Work,
		School:                    req.Personal.School,
		Qualification:             req.Personal.Qualification,
		CurrentPlace:              req.Personal.CurrentPlace,
		AncestralHome:             req.Personal.AncestralHome,
		Economic:                  economic,
		Hobbies:                   fmt.Sprint(req.Personal.Hobbies),
		Height:                    req.Personal.Height,
		Weight:                    req.Personal.Weight,
		OriginalFamilyComposition: req.Personal.OriginalFamilyComposition,
		ParentsSituation:          req.Personal.ParentsSituation,
		Remarks:                   req.Personal.Remarks,
		Gender:                    req.Personal.Gender,
	}

	data, err := cApi.Svc.MatchCandidate(dbPersonal, req.AttributesMap)
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "匹配失败",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "匹配成功",
		Data:    data,
	})
}

func (cApi *CandidateApi) GetPersonalInfo(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	return nil
}

func (cApi *CandidateApi) GetPersonalInfoByID(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	params := r.Params()
	id, _ := params.Int64("id")
	fmt.Println(id)
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "匹配成功",
		Data:    nil,
	})
}

func (cApi *CandidateApi) UpdatePersonalInfo(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	return nil
}

func (cApi *CandidateApi) DeletePersonalInfo(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	return nil
}

func (cApi *CandidateApi) UpdatePersonalInfoAndScore(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	return nil
}
