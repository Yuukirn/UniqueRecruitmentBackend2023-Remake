package models

import (
	"UniqueRecruitmentBackend/global"
	"UniqueRecruitmentBackend/internal/request"
	"errors"
	"time"
)

type ApplicationEntity struct {
	Common
	Grade     string `gorm:"not null"` //constants.Grade
	Institute string `gorm:"not null"`
	Major     string `gorm:"not null"`
	Rank      string `gorm:"not null"` //constants.Rank
	Group     string `gorm:"not null"` //constants.Group
	Intro     string `gorm:"not null"`
	IsQuick   bool   `gorm:"column:isQuick;not null"`
	Referrer  string

	Resume string

	Abandoned                 bool               `gorm:"not null; default false" `
	Rejected                  bool               `gorm:"not null; default false"`
	Step                      string             `gorm:"not null"`                                                      //constants.Step
	CandidateID               string             `gorm:"column:candidateId;uniqueIndex:UQ_CandidateID_RecruitmentID"`   //manytoone
	RecruitmentID             string             `gorm:"column:recruitmentId;uniqueIndex:UQ_CandidateID_RecruitmentID"` //manytoone
	InterviewAllocationsGroup time.Time          `gorm:"column:interviewAllocationsGroup;"`
	InterviewAllocationsTeam  time.Time          `gorm:"column:interviewAllocationsTeam;"`
	InterviewSelections       []*InterviewEntity `gorm:"many2many:interview_selections; constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"` //manytomany
	Comments                  []CommentEntity    `gorm:"foreignKey:ApplicationID; references:Uid; constraint:OnDelete:CASCADE;"`        //onetomany
}

func (a ApplicationEntity) TableName() string {
	return "applications"
}
func CreateAndSaveApplication(req *request.CreateApplicationRequest, filename string) (*ApplicationEntity, error) {
	db := global.GetDB()
	row := db.Model(&ApplicationEntity{}).Where("recruitmentId = ?", req.RecruitmentID).Find(&ApplicationEntity{}).RowsAffected

	//check now user's recruitment application >1
	if row != 0 {
		return nil, errors.New("A candidate can only apply once at the same recruitment")
	}
	a := ApplicationEntity{
		Grade:         req.Grade,
		Institute:     req.Institute,
		Major:         req.Major,
		Rank:          req.Rank,
		Group:         req.Group,
		Intro:         req.Intro,
		RecruitmentID: req.RecruitmentID,
		Referrer:      req.Referrer,
		IsQuick:       req.IsQuick,
		Resume:        filename,
		// TODO(wwb)
		// Add step status
		Step: "",
	}
	err := db.Model(&ApplicationEntity{}).Create(&a).Error
	return &a, err
}

/*

 */
