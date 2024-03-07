package candidate

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"

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
	personalInfo PersonalInfo) error {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	personalInfo.PersonCode = newUUID.String()
	tx := candidateService.DB.Begin()
	if err := candidateService.DB.Create(&personalInfo).Error; err != nil {
		return err
	}
	work, area := build_work_area_tree(candidateService.DB)
	//将personalInfo生成可量化数据，储存
	candidate := &Candidate{
		PersonCode:    personalInfo.PersonCode,
		BirthYear:     personalInfo.BirthYear,
		Work:          get_work_or_area_save(work, personalInfo.Work),
		Qualification: get_qualification_index_map(personalInfo.Qualification),
		CurrentPlace:  get_work_or_area_save(area, personalInfo.CurrentPlace),
		AncestralHome: get_work_or_area_save(area, personalInfo.AncestralHome),
		Economic:      get_economic(personalInfo.Economic),
		Height:        personalInfo.Height,
		Weight:        personalInfo.Weight,
		Score:         0.0,
	}
	if err := candidateService.DB.Create(candidate).Error; err != nil {
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

func MatchDataToShowData(candidateReq []CandidateReq, db *gorm.DB, gender string) (*[]PersonalInfoResult, error) {
	UUIDCandidateMap := make(map[string]CandidateReq, len(candidateReq))
	for _, c := range candidateReq {
		UUIDCandidateMap[c.PersonCode] = c
	}
	var personalInfos []PersonalInfo
	if db = db.Where("gender=?", gender).
		Find(&personalInfos); db.Error != nil {
		return nil, errors.New("数据库查询错误!")
	}
	personalInfoResults := make([]PersonalInfoResult, 0, len(personalInfos))
	for _, p := range personalInfos {
		c := UUIDCandidateMap[p.PersonCode]
		personalInfoResult := PersonalInfoResult{
			PersonalInfo: p,
			Score:        c.Score,
		}
		personalInfoResults = append(personalInfoResults, personalInfoResult)
	}
	sort.Slice(personalInfoResults, func(i, j int) bool {
		return personalInfoResults[i].Score > personalInfoResults[j].Score
	})
	return &personalInfoResults, nil
}

func (candidateService *CandidateService) MatchCandidate(
	personalInfo PersonalInfo,
	attributes_map map[string]float64) (*[]PersonalInfoResult, error) {
	work, area := build_work_area_tree(candidateService.DB)
	var candidates []Candidate
	var gender string
	if personalInfo.Gender == "男" {
		gender = "女"
	} else {
		gender = "男"
	}
	result := candidateService.DB.
		Joins("LEFT JOIN personal_infos ON candidates.person_code = personal_infos.person_code").
		Where("personal_infos.gender=?", gender).
		Find(&candidates)
	if result.Error != nil {
		return nil, errors.New("查询数据库候选人失败!")
	}
	candidate_reqs := make([]CandidateReq, 0, 500)
	for _, c := range candidates {
		var works []int
		err := json.Unmarshal(c.Work, &works)
		if err != nil {
			log.Fatal(err)
		}
		var current_place []int
		err = json.Unmarshal(c.CurrentPlace, &current_place)
		if err != nil {
			log.Fatal(err)
		}
		var ancestral_home []int
		err = json.Unmarshal(c.AncestralHome, &ancestral_home)
		if err != nil {
			log.Fatal(err)
		}
		candidate_req := CandidateReq{
			PersonCode:    c.PersonCode,
			BirthYear:     c.BirthYear,
			Work:          works,
			Qualification: c.Qualification,
			CurrentPlace:  current_place,
			AncestralHome: ancestral_home,
			Economic:      get_economic(personalInfo.Economic),
			Height:        personalInfo.Height,
			Weight:        personalInfo.Weight,
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
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}
	if response.Code != 200 {
		return nil, errors.New("请求匹配服务失败")
	}
	data, err := MatchDataToShowData(response.Data, candidateService.DB, gender)
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
	case "专科":
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

func get_work_or_area_save(node *utils.TreeNode, id int) json.RawMessage {
	parent_id := utils.FindParentID(node, id)
	parent1_id := utils.FindParentID(node, parent_id)
	marshal, err := json.Marshal([]int{parent1_id, parent_id, id})
	if err != nil {
		log.Println(err)
	}
	return marshal
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
