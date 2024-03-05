package candidate

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

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

func (candidateService *CandidateService) MatchCandidate(
	personalInfo PersonalInfo,
	attributes_map map[string]float64) *[]CandidateReq {
	work, area := build_work_area_tree(candidateService.DB)
	var candidates []Candidate
	result := candidateService.DB.Table("candidates").
		Limit(5000).
		Find(&candidates)
	if result.Error != nil {
		panic("failed to query database")
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
	return &response.Data
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
