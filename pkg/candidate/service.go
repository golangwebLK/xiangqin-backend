package candidate

import "gorm.io/gorm"

type CandidateService struct {
	DB *gorm.DB
}

func NewCandidateService(db *gorm.DB) *CandidateService {
	return &CandidateService{
		DB: db,
	}
}

func (candidateService *CandidateService) SavePersonalInfo(personalInfo PersonalInfo) error {
	if err := candidateService.DB.Create(personalInfo).Error; err != nil {
		return err
	}
	return nil
}
