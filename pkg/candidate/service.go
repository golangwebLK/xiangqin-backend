package candidate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
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
		Work:          get_work_or_area(work, personalInfo.Work),
		Qualification: get_qualification_index_map(personalInfo.Qualification),
		CurrentPlace:  get_work_or_area(area, personalInfo.CurrentPlace),
		AncestralHome: get_work_or_area(area, personalInfo.AncestralHome),
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
	Candidate  Candidate          `json:"candidate"`
	Candidates []Candidate        `json:"candidates"`
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
	attributes_map map[string]float64) error {
	work, area := build_work_area_tree(candidateService.DB)
	var candidates []Candidate
	result := candidateService.DB.Limit(200).
		Find(&candidates)
	if result.Error != nil {
		panic("failed to query database")
	}
	candidate := Candidate{
		PersonCode:    personalInfo.PersonCode,
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
	fmt.Println(candidate)
	reqData := RequestData{
		Candidate:  candidate,
		Candidates: candidates,
		Attributes: attributes_map,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST",
		candidateService.MatchServer+"/matching",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(resp.StatusCode)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var i string
	fmt.Println(body)
	err = json.Unmarshal(body, &i)
	if err != nil {
		log.Fatal("JSON unmarshalling failed: ", err)
		return err
	}
	fmt.Println(i)
	return nil
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func get_economic(economic interface{}) float64 {
	e, ok := economic.(EconomicInfo)
	if !ok {
		log.Println("economicStr 不是字符串类型")
		return 0
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

func get_work_or_area(node *utils.TreeNode, id int) json.RawMessage {
	parent_id := utils.FindParentID(node, id)
	parent1_id := utils.FindParentID(node, parent_id)
	jsonData, err := json.Marshal([]int{parent1_id, parent_id, id})
	if err != nil {
		log.Println("JSON marshaling failed:", err)
	}

	return jsonData
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
