package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// RegisterStatRequest .
type RegisterStatRequest struct {
	wsgin.MustAuthRequest

	BeginDate string `form:"begin_date" json:"begin_date" binding:"required"` // 开始日期
	EndDate   string `form:"end_date" json:"end_date" binding:"required"`     // 结束日期
}

// RegisterStatResponse .
type RegisterStatResponse struct {
	wsgin.BaseResponse

	Data []*model.RegisterStat `json:"data"`
}

// New .
func (r *RegisterStatRequest) New() wsgin.Process {
	return &RegisterStatRequest{}
}

// Extract .
func (r *RegisterStatRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 统计用户注册数
// @Summary 统计用户注册数
// @Description stat customer register num
// @Tags 统计
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param begin_date query string true "开始日期"
// @Param end_date query string true "结束日期"
// @Success 200 {object} server.RegisterStatResponse "{"status":true}"
// @Router /stat/register [get]
func (r *RegisterStatRequest) Exec(ctx context.Context) interface{} {
	resp := RegisterStatResponse{}

	data, code, err := svc.GetRegisterStat(ctx, r.BeginDate, r.EndDate)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
