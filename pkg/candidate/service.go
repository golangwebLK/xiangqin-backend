package candidate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
	if err := candidateService.DB.Create(&personalInfo).Error; err != nil {
		return err
	}
	return nil
}

type Candidate struct {
	BirthYear     int     `json:"birth_year"`     // 实际年龄
	Work          []int   `json:"work"`           // 按照包含关系，填入编号
	Qualification int     `json:"qualification"`  // 学历编号1-7，
	CurrentPlace  []int   `json:"current_place"`  // 按照包含关系，填入编号
	AncestralHome []int   `json:"ancestral_home"` // 按照包含关系，填入编号
	Economic      float64 `json:"economic"`       // 实际财富
	Height        float64 `json:"height"`         // 实际身高
	Weight        float64 `json:"weight"`         // 实际体重
	Score         float64 `json:"score"`
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

type Person struct {
	BirthYear     int         `gorm:"column:birth_year"`
	Work          int         `gorm:"column:work"`
	Qualification string      `gorm:"column:qualification"`
	CurrentPlace  int         `gorm:"column:current_place"`
	AncestralHome int         `gorm:"column:ancestral_home"`
	Economic      interface{} `gorm:"column:economic"`
	Height        int         `gorm:"column:height"`
	Weight        int         `gorm:"column:weight"`
}

func (candidateService *CandidateService) MatchCandidate(
	personalInfo PersonalInfo,
	attributes_map map[string]float64) error {
	work, area := build_work_area_tree(candidateService.DB)
	var persons []Person
	result := candidateService.DB.Table("public.personal_infos").
		Select("birth_year, work, qualification, current_place, ancestral_home, economic, height, weight").
		Find(&persons)
	if result.Error != nil {
		panic("failed to query database")
	}
	candidates := make([]Candidate, len(persons), len(persons))
	for index, person := range persons {
		candidates[index].BirthYear = person.BirthYear
		candidates[index].Work = get_work_or_area(work, person.Work)
		candidates[index].Qualification = get_qualification_index_map(person.Qualification)
		candidates[index].CurrentPlace = get_work_or_area(area, person.CurrentPlace)
		candidates[index].AncestralHome = get_work_or_area(area, person.AncestralHome)
		candidates[index].Economic = get_economic(person.Economic)
		candidates[index].Height = float64(person.Height)
		candidates[index].Weight = float64(person.Weight)
	}
	var candidate Candidate
	candidate.BirthYear = personalInfo.BirthYear
	candidate.Work = get_work_or_area(work, personalInfo.Work)
	candidate.Qualification = get_qualification_index_map(personalInfo.Qualification)
	candidate.CurrentPlace = get_work_or_area(area, personalInfo.CurrentPlace)
	candidate.AncestralHome = get_work_or_area(area, personalInfo.AncestralHome)
	candidate.Economic = get_economic(PersonalInfo.Economic)
	candidate.Height = float64(personalInfo.Height)
	candidate.Weight = float64(personalInfo.Weight)
	reqData := RequestData{
		Candidate:  candidate,
		Candidates: candidates,
		Attributes: attributes_map,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		log.Println("序列化时出错:", err)
	}
	resp, err := http.Post(candidateService.MatchServer+"/matching", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("发起请求时出错:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应时出错:", err)
	}
	fmt.Println(body)
	return nil
}

func get_economic(economicStr interface{}) float64 {
	str, ok := economicStr.(string)
	if !ok {
		log.Println("economicStr 不是字符串类型")
		return 0
	}
	data := []byte(str)
	var economic EconomicInfo
	err := json.Unmarshal([]byte(data), &economic)
	if err != nil {
		log.Println(err)
		return 0
	}
	i, err := strconv.Atoi(economic.Savings)
	if err != nil {
		log.Println(err)
		return 0
	}
	return float64(i)
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
