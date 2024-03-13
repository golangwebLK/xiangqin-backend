package candidate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"xiangqin-backend/pkg/middleware"

	xiangqin_backend "xiangqin-backend"
	"xiangqin-backend/utils"

	"gorm.io/gorm"
)

type CandidateService struct {
	DB          *gorm.DB
	MatchServer string
}

func NewCandidateService(db *gorm.DB, cfg *xiangqin_backend.Config) *CandidateService {
	return &CandidateService{
		DB:          db,
		MatchServer: cfg.MatchService.Server,
	}
}

func (candidateService *CandidateService) SavePersonalInfo(
	ctx context.Context,
	createPersonReq CreatePersonReq) error {
	msg := ctx.Value("msg").(middleware.Msg)
	UUID := uuid.NewString()
	economic, err := json.Marshal(createPersonReq.Personal.Economic)
	if err != nil {
		return err
	}
	personalInfo := PersonalInfo{
		PersonCode:                UUID,
		RealName:                  createPersonReq.Personal.RealName,
		Picture:                   createPersonReq.Personal.Picture,
		BirthYear:                 createPersonReq.Personal.BirthYear,
		Telephone:                 createPersonReq.Personal.Telephone,
		WeChat:                    createPersonReq.Personal.WeChat,
		Work:                      createPersonReq.Personal.Work,
		School:                    createPersonReq.Personal.School,
		Qualification:             createPersonReq.Personal.Qualification,
		CurrentPlace:              createPersonReq.Personal.CurrentPlace,
		AncestralHome:             createPersonReq.Personal.AncestralHome,
		Economic:                  economic,
		Hobbies:                   fmt.Sprint(createPersonReq.Personal.Hobbies),
		Height:                    createPersonReq.Personal.Height,
		Weight:                    createPersonReq.Personal.Weight,
		OriginalFamilyComposition: createPersonReq.Personal.OriginalFamilyComposition,
		ParentsSituation:          createPersonReq.Personal.ParentsSituation,
		Remarks:                   createPersonReq.Personal.Remarks,
		Gender:                    createPersonReq.Personal.Gender,
		CompanyCode:               msg.CompanyCode,
	}
	personalLike := PersonalLike{
		PersonCode:                UUID,
		BirthYear:                 createPersonReq.PersonLike.BirthYear,
		Work:                      createPersonReq.PersonLike.Work,
		Qualification:             createPersonReq.PersonLike.Qualification,
		CurrentPlace:              createPersonReq.PersonLike.CurrentPlace,
		AncestralHome:             createPersonReq.PersonLike.AncestralHome,
		Economic:                  economic,
		Hobbies:                   fmt.Sprint(createPersonReq.PersonLike.Hobbies),
		Height:                    createPersonReq.PersonLike.Height,
		Weight:                    createPersonReq.PersonLike.Weight,
		OriginalFamilyComposition: createPersonReq.PersonLike.OriginalFamilyComposition,
		ParentsSituation:          createPersonReq.PersonLike.ParentsSituation,
		CompanyCode:               msg.CompanyCode,
		Remarks:                   createPersonReq.PersonLike.Remarks,
	}
	tx := candidateService.DB.Begin()
	if err = candidateService.DB.Create(&personalInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = candidateService.DB.Create(&personalLike).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

type RequestData struct {
	Candidate  CandidateReq       `json:"candidate"`
	Candidates []CandidateReq     `json:"candidates"`
	Attributes map[string]float64 `json:"attributes"`
}

type Work struct {
	WorkId   int    `gorm:"column:work_id"`
	ParentId int    `gorm:"column:parent_id"`
	Name     string `gorm:"column:name"`
}

type DouArea struct {
	AreaId   int    `gorm:"column:area_id"`
	ParentId int    `gorm:"column:parent_id"`
	Name     string `gorm:"column:name"`
}

type PersonalInfoResult struct {
	PersonalInfo PersonalInfo `json:"personalInfo"`
	Score        float64      `json:"score"`
}

func MatchDataToShowData(
	candidateReqs []CandidateReq,
	personalInfos []PersonalInfo) (*[]PersonalInfoResult, error) {
	personalInfoResults := make([]PersonalInfoResult, 0, len(personalInfos))
	UUIDPersonalInfoMap := make(map[string]PersonalInfo, len(personalInfos))
	for _, p := range personalInfos {
		UUIDPersonalInfoMap[p.PersonCode] = p
	}
	for _, c := range candidateReqs {
		if c.Score < 60 && len(personalInfoResults) > 20 {
			break
		}
		personalInfoResult := PersonalInfoResult{
			PersonalInfo: UUIDPersonalInfoMap[c.PersonCode],
			Score:        c.Score,
		}
		personalInfoResults = append(personalInfoResults, personalInfoResult)
	}
	return &personalInfoResults, nil
}

func (candidateService *CandidateService) MatchCandidate(
	ctx context.Context,
	personalInfo PersonalInfo,
	attributes_map map[string]float64) (*[]PersonalInfoResult, error) {
	msg := ctx.Value("msg").(*middleware.Msg)
	work, area := build_work_area_tree(candidateService.DB)
	var personalInfos []PersonalInfo
	var gender string
	if personalInfo.Gender == "男" {
		gender = "女"
	} else {
		gender = "男"
	}
	result := candidateService.DB.
		Where("gender=? and company_code=?", gender, msg.CompanyCode).
		Find(&personalInfos)
	if result.Error != nil {
		return nil, errors.New("查询数据库候选人失败!")
	}
	candidate_reqs := make([]CandidateReq, 0, len(personalInfos))
	for _, p := range personalInfos {
		candidate_req := CandidateReq{
			PersonCode:    p.PersonCode,
			BirthYear:     p.BirthYear,
			Work:          get_work_or_area(work, p.Work),
			Qualification: get_qualification_index_map(p.Qualification),
			CurrentPlace:  get_work_or_area(area, p.CurrentPlace),
			AncestralHome: get_work_or_area(area, p.AncestralHome),
			Economic:      get_economic(p.Economic),
			Height:        p.Height,
			Weight:        p.Weight,
			Score:         0.0,
		}
		candidate_reqs = append(candidate_reqs, candidate_req)
	}
	candidate := CandidateReq{
		PersonCode:    uuid.NewString(),
		BirthYear:     personalInfo.BirthYear,
		Work:          get_work_or_area(work, personalInfo.Work),
		Qualification: get_qualification_index_map(personalInfo.Qualification),
		CurrentPlace:  get_work_or_area(area, personalInfo.CurrentPlace),
		AncestralHome: get_work_or_area(area, personalInfo.AncestralHome),
		Economic:      get_economic(personalInfo.Economic),
		Height:        personalInfo.Height,
		Weight:        personalInfo.Weight,
		Score:         0.0,
	}
	reqData := RequestData{
		Candidate:  candidate,
		Candidates: candidate_reqs,
		Attributes: attributes_map,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(candidateService.MatchServer+"/matching",
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("请求算法匹配服务系统错误，相应状态:" + resp.Status)
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}
	if response.Code != 200 {
		return nil, errors.New("请求匹配服务失败")
	}
	data, err := MatchDataToShowData(response.Data, personalInfos)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Response struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    []CandidateReq `json:"data"`
}

func get_economic(economic json.RawMessage) float64 {
	var e EconomicInfo
	err := json.Unmarshal(economic, &e)
	if err != nil {
		log.Println("get_economic函数有问题")
	}
	return e.Savings + e.CarMoney + e.HouseMoney
}

func get_qualification_index_map(qualification string) int {
	switch qualification {
	case "小学":
		return 1
	case "初中":
		return 2
	case "高中":
		return 3
	case "大专":
		return 4
	case "本科":
		return 5
	case "硕士":
		return 6
	case "博士":
		return 7
	}
	return 0
}

func get_work_or_area(node *utils.TreeNode, id int) []int {
	parent_id := utils.FindParentID(node, id)
	parent1_id := utils.FindParentID(node, parent_id)
	return []int{parent1_id, parent_id, id}
}

func build_work_area_tree(db *gorm.DB) (*utils.TreeNode, *utils.TreeNode) {
	var works []Work
	if db := db.
		Table("works").
		Select("work_id, parent_id, name").
		Find(&works); db.Error != nil {
		log.Println(db.Error)
	}
	var work_root *utils.TreeNode
	for _, work := range works {
		work_root = utils.Insert(work_root,
			work.WorkId, work.ParentId, work.Name)
	}
	var areas []DouArea
	if db := db.
		Table("dou_area").
		Select("area_id, parent_id, name").
		Find(&areas); db.Error != nil {
		log.Println(db.Error)
	}
	var area_root *utils.TreeNode
	for _, area := range areas {
		area_root = utils.Insert(area_root,
			area.AreaId, area.ParentId, area.Name)
	}
	return work_root, area_root
}

func (svc *CandidateService) GetPersonalInfo(ctx context.Context, page, pageSize int, name string) (*[]PersonalInfo, error) {
	msg := ctx.Value("msg").(middleware.Msg)
	var personalInfos []PersonalInfo
	query := svc.DB.Model(&PersonalInfo{}).Where("company_code=?", msg.CompanyCode)
	if name != "" {
		query = query.Where("real_name=?", "%"+name+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	offset, err := CalculateOffset(page, pageSize, total)
	if err != nil {
		return nil, err
	}
	if err = query.Offset(offset).Limit(pageSize).Find(&personalInfos).Error; err != nil {
		return nil, err
	}
	return &personalInfos, nil
}
func CalculateOffset(page, pageSize int, totalRecords int64) (int, error) {
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pageSize
	if offset > int(totalRecords) {
		return 0, errors.New("page out of range")
	}

	return offset, nil
}

func (svc *CandidateService) GetPersonalInfoByID(ctx context.Context, id int) (PersonalInfo, error) {
	msg := ctx.Value("msg").(middleware.Msg)
	var personalInfo PersonalInfo
	if err := svc.DB.Where("id=? and company_code=?", id, msg.CompanyCode).Find(&personalInfo).Error; err != nil {
		return PersonalInfo{}, err
	}
	return personalInfo, nil
}

func (svc *CandidateService) UpdatePersonalInfo(ctx context.Context, info PersonalInfo) error {
	msg := ctx.Value("msg").(middleware.Msg)
	if err := svc.DB.Where("company_code=?", msg.CompanyCode).Updates(&info).Error; err != nil {
		return err
	}
	return nil
}

func (svc *CandidateService) UpdatePersonalLike(ctx context.Context, like PersonalLike) error {
	msg := ctx.Value("msg").(middleware.Msg)
	if err := svc.DB.Where("company_code=?", msg.CompanyCode).Updates(&like).Error; err != nil {
		return err
	}
	return nil
}

func (svc *CandidateService) DeletePersonalInfo(ctx context.Context, code string) error {
	msg := ctx.Value("msg").(middleware.Msg)
	tx := svc.DB.Begin()
	if err := svc.DB.Where("company_code=? and person_code=?", msg.CompanyCode, code).Delete(&PersonalInfo{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := svc.DB.Where("company_code=? and person_code=?", msg.CompanyCode, code).Delete(&PersonalLike{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
