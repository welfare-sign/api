package server

import (
	"context"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// HelpCheckinRequest .
type HelpCheckinRequest struct {
	wsgin.AuthRequest

	CustomerID uint64 `form:"customer_id" json:"customer_id" binding:"required"` // 补签客户ID
}

// HelpCheckinResponse .
type HelpCheckinResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *HelpCheckinRequest) New() wsgin.Process {
	return &HelpCheckinRequest{}
}

// Extract .
func (r *HelpCheckinRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 为他人补签
// @Summary 为他人补签
// @Description customer help other checkin
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Param args body server.HelpCheckinRequest true "参数"
// @Success 200 {object} server.HelpCheckinResponse	"{"status":true}"
// @Router /customers/checkin_record/help [post]
func (r *HelpCheckinRequest) Exec(ctx context.Context) interface{} {
	resp := HelpCheckinResponse{}

	code, err := svc.HelpCheckinRecord(ctx, r.TokenParames.UID, r.CustomerID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
