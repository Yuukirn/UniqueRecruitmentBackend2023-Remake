package controllers

import (
	"UniqueRecruitmentBackend/internal/common"
	error2 "UniqueRecruitmentBackend/internal/error"
	"UniqueRecruitmentBackend/internal/models"
	"UniqueRecruitmentBackend/internal/request"
	"UniqueRecruitmentBackend/internal/utils"
	"github.com/gin-gonic/gin"
	"time"
)

// Post recruitment/
// admin role

// CreateRecruitment create new recruitment
func CreateRecruitment(c *gin.Context) {
	var req request.CreateRecruitmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, error2.RequestBodyError.WithDetail(err.Error()))
		return
	}
	if time.Now().After(req.Beginning) || req.Beginning.After(req.Deadline) || req.Deadline.After(req.End) {
		common.Error(c, error2.RequestBodyError.WithDetail("time set up wrong"))
		return
	}
	recruitmentId, err := models.CreateRecruitment(&req)
	if err != nil {
		common.Error(c, error2.SaveDatabaseError.WithData("recruitment"))
		return
	}
	common.Success(c, "Success create recruitment", map[string]interface{}{
		"rid": recruitmentId,
	})
}

// Post recruitment/:rid/schedule
// admin role

// UpdateRecruitment update recruitment details
func UpdateRecruitment(c *gin.Context) {
	recruitmentId := c.Param("rid")
	var req request.UpdateRecruitmentRequest
	if err := c.ShouldBindJSON(&req); err != nil || recruitmentId == "" {
		common.Error(c, error2.RequestBodyError.WithDetail(err.Error()))
		return
	}
	if err := models.UpdateRecruitment(recruitmentId, &req); err != nil {
		common.Error(c, error2.UpdateDatabaseError.WithData("recruitment").WithDetail(err.Error()))
		return
	}
	common.Success(c, "Success update recruitment", map[string]interface{}{
		"rid": recruitmentId,
	})
}

// Get recruitment/:rid

// GetRecruitmentById get recruitment details by id
func GetRecruitmentById(c *gin.Context) {
	recruitmentId := c.Param("rid")
	if recruitmentId == "" {
		common.Error(c, error2.RequestBodyError.WithDetail("lost http query params [rid]"))
		return
	}
	resp, err := models.GetRecruitmentById(recruitmentId)
	if err != nil {
		common.Error(c, error2.GetDatabaseError.WithData("recruitment").WithDetail(err.Error()))
		return
	}
	common.Success(c, "Success get one recruitment", resp)
}

// Get recruitment/

// GetAllRecruitment get all recruitment details
func GetAllRecruitment(c *gin.Context) {
	// TODO(wwb)
	// compare member joinin time and recruitment time
	resp, err := models.GetAllRecruitment()
	if err != nil {
		common.Error(c, error2.GetDatabaseError.WithData("recruitment").WithDetail(err.Error()))
		return
	}
	common.Success(c, "Success get all recruitment", resp)
}

func compareTime(a string, b string) bool {
	ta := utils.TimeParse(a)
	tb := utils.TimeParse(b)
	return ta.After(tb)
}
