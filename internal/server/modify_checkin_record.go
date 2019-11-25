package server

import (
	"context"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// ModifyCheckinRecordRequest .
type ModifyCheckinRecordRequest struct {
	wsgin.AuthRequest

	CheckinRecordID uint64 `json:"checkin_record_id"` // 签到记录ID
	Status          string `json:"status"`            // 用户状态：U，A
}

// ModifyCheckinRecordResponse .
type ModifyCheckinRecordResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *ModifyCheckinRecordRequest) New() wsgin.Process {
	return &ModifyCheckinRecordRequest{}
}

// Extract .
func (r *ModifyCheckinRecordRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 后台用户更改用户签到记录(临时)
// @Summary 后台用户更改用户签到记录(临时)
// @Description users modify customer checkin record
// @Security ApiKeyAuth
// @Tags 后台用户
// @Accept json
// @Produce json
// @Param args body server.ModifyCheckinRecordRequest true "参数"
// @Success 200 {object} server.ModifyCheckinRecordResponse	"{"status":true}"
// @Router /users/checkin_record_list/modify [post]
func (r *ModifyCheckinRecordRequest) Exec(ctx context.Context) interface{} {
	resp := ModifyCheckinRecordResponse{}

	code, err := svc.ModifyCustomerCheckinRecord(ctx, r.CheckinRecordID, r.Status)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
