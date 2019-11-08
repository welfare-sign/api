package server

import (
	"context"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// ExecCheckinRecordRequest .
type ExecCheckinRecordRequest struct {
	wsgin.AuthRequest

	Day uint64 `form:"day" json:"day" binding:"required"` // 签到第天数
}

// ExecCheckinRecordResponse .
type ExecCheckinRecordResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *ExecCheckinRecordRequest) New() wsgin.Process {
	return &ExecCheckinRecordRequest{}
}

// Extract .
func (r *ExecCheckinRecordRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 用户执行签到
// @Summary 用户执行签到
// @Description customer exec checkin record
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Param args body server.ExecCheckinRecordRequest true "参数"
// @Success 200 {object} server.ExecCheckinRecordResponse	"{"status":true}"
// @Router /customers/checkin_record [post]
func (r *ExecCheckinRecordRequest) Exec(ctx context.Context) interface{} {
	resp := ExecCheckinRecordResponse{}

	code, err := svc.ExecCheckinRecord(ctx, r.TokenParames.UID, r.Day)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
