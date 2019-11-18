package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/pkg/wsgin"
)

// MerchantDelRequest 删除商户
type MerchantDelRequest struct {
	wsgin.AuthRequest

	MerchantID uint64 `form:"merchant_id" json:"merchant_id"` // 商户ID
}

// MerchantDelResponse .
type MerchantDelResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *MerchantDelRequest) New() wsgin.Process {
	return &MerchantDelRequest{}
}

// Extract .
func (r *MerchantDelRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 删除商户
// @Summary 删除商户
// @Description delete merchant
// @Security ApiKeyAuth
// @Tags 商户
// @Accept json
// @Produce json
// @Param args body server.MerchantDelRequest true "参数"
// @Success 200 {object} server.MerchantDelResponse "{"status":true}"
// @Router /merchants [delete]
func (r *MerchantDelRequest) Exec(ctx context.Context) interface{} {
	resp := MerchantDelResponse{}

	code, err := svc.DeleteMerchant(ctx, r.MerchantID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
