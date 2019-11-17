package server

import (
	"context"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// RefreshCheckinRecordRequest .
type RefreshCheckinRecordRequest struct {
	wsgin.AuthRequest
}

// RefreshCheckinRecordResponse .
type RefreshCheckinRecordResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *RefreshCheckinRecordRequest) New() wsgin.Process {
	return &RefreshCheckinRecordRequest{}
}

// Extract .
func (r *RefreshCheckinRecordRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 用户重新签到
// @Summary 用户重新签到
// @Description customer refresh checkin record
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Success 200 {object} server.RefreshCheckinRecordResponse "{"status":true}"
// @Router /customers/checkin_record/refresh [post]
func (r *RefreshCheckinRecordRequest) Exec(ctx context.Context) interface{} {
	resp := RefreshCheckinRecordResponse{}

	code, err := svc.RefreshCheckinRecord(ctx, r.TokenParames.UID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
