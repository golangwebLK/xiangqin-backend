package candidate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

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
	Picture                   string       `json:"picture"`
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

type PersonLike struct {
	BirthYear                 int          `json:"birth_year"`
	Work                      int          `json:"work"`
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
}

type CreatePersonReq struct {
	Personal   Personal   `json:"personal"`
	PersonLike PersonLike `json:"personLike"`
}

func (cApi *CandidateApi) CreateCandidate(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	ctx := r.Request.Context()
	var createPersonReq CreatePersonReq
	if err := json.NewDecoder(r.Body).Decode(&createPersonReq); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "请求参数错误",
			Data:    err,
		})
	}
	if err := cApi.Svc.SavePersonalInfo(ctx, createPersonReq); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "保存错误",
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
	ctx := r.Request.Context()
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

	data, err := cApi.Svc.MatchCandidate(ctx, dbPersonal, req.AttributesMap)
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
	ctx := r.Request.Context()
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "请求参数错误",
			Data:    err,
		})
	}
	pageInt, _ := strconv.Atoi(queryValues.Get("page"))
	pageSizeInt, _ := strconv.Atoi(queryValues.Get("pageSize"))
	name := queryValues.Get("name")
	personalinfos, err := cApi.Svc.GetPersonalInfo(ctx, pageInt, pageSizeInt, name)
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "查询错误",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "查询成功",
		Data:    personalinfos,
	})
}

func (cApi *CandidateApi) GetPersonalInfoByID(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	ctx := r.Request.Context()
	params := r.Params()
	id, err := params.Int64("id")
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "字段错误",
			Data:    err,
		})
	}
	personalInfo, err := cApi.Svc.GetPersonalInfoByID(ctx, int(id))
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "查询失败",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "查询成功",
		Data:    personalInfo,
	})
}

func (cApi *CandidateApi) UpdatePersonalInfo(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	ctx := r.Request.Context()
	var personalInfo PersonalInfo
	if err := json.NewDecoder(r.Body).Decode(&personalInfo); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "字段错误",
			Data:    err,
		})
	}
	if err := cApi.Svc.UpdatePersonalInfo(ctx, personalInfo); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "修改失败",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "修改成功",
		Data:    nil,
	})
}

func (cApi *CandidateApi) UpdatePersonalLike(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	ctx := r.Request.Context()
	var personalLike PersonalLike
	if err := json.NewDecoder(r.Body).Decode(&personalLike); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "字段错误",
			Data:    err,
		})
	}
	if err := cApi.Svc.UpdatePersonalLike(ctx, personalLike); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "修改失败",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "修改成功",
		Data:    nil,
	})
}

func (cApi *CandidateApi) DeletePersonalInfo(
	rw http.ResponseWriter,
	r bunrouter.Request) error {
	ctx := r.Request.Context()
	params := r.Params()
	code, _ := params.Get("code")
	if err := cApi.Svc.DeletePersonalInfo(ctx, code); err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "删除失败",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "删除成功",
		Data:    nil,
	})
}

func (cApi *CandidateApi) UploadImage(rw http.ResponseWriter, r bunrouter.Request) error {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "文件大小不能超过20M",
			Data:    err,
		})
	}
	file, handler, err := r.FormFile("image")
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusBadRequest,
			Message: "获取文件失败",
			Data:    err,
		})
	}
	defer file.Close()
	tempFile, err := os.CreateTemp("uploads", "upload-*.jpg")
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "创建文件失败",
			Data:    err,
		})
	}
	defer tempFile.Close()
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return bunrouter.JSON(rw, utils.ResponseData{
			Status:  http.StatusInternalServerError,
			Message: "文件写入失败",
			Data:    err,
		})
	}
	return bunrouter.JSON(rw, utils.ResponseData{
		Status:  http.StatusOK,
		Message: "上传成功",
		Data:    handler.Filename,
	})
}
