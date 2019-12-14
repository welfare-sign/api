package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// CheckinStatRequest .
type CheckinStatRequest struct {
	wsgin.MustAuthRequest

	BeginDate string `form:"begin_date" json:"begin_date" binding:"required"` // 开始日期
	EndDate   string `form:"end_date" json:"end_date" binding:"required"`     // 结束日期
}

// CheckinStatResponse .
type CheckinStatResponse struct {
	wsgin.BaseResponse

	Data []*model.CheckinStat `json:"data"`
}

// New .
func (r *CheckinStatRequest) New() wsgin.Process {
	return &CheckinStatRequest{}
}

// Extract .
func (r *CheckinStatRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 统计用户执行签到数目
// @Summary 统计用户执行签到数目
// @Description stat customer exec checkin num
// @Tags 统计
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param begin_date query string true "开始日期"
// @Param end_date query string true "结束日期"
// @Success 200 {object} server.CheckinStatResponse "{"status":true}"
// @Router /stat/checkin [get]
func (r *CheckinStatRequest) Exec(ctx context.Context) interface{} {
	resp := CheckinStatResponse{}

	data, code, err := svc.GetCheckinStat(ctx, r.BeginDate, r.EndDate)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
