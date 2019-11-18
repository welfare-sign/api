package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/pkg/wsgin"
)

// MerchantDisableRequest 禁用商户
type MerchantDisableRequest struct {
	wsgin.AuthRequest

	MerchantID uint64 `form:"merchant_id" json:"merchant_id"` // 商户ID
}

// MerchantDisableResponse .
type MerchantDisableResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *MerchantDisableRequest) New() wsgin.Process {
	return &MerchantDisableRequest{}
}

// Extract .
func (r *MerchantDisableRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 禁用商户
// @Summary 禁用商户
// @Description disable merchant
// @Security ApiKeyAuth
// @Tags 商户
// @Accept json
// @Produce json
// @Param args body server.MerchantDisableRequest true "参数"
// @Success 200 {object} server.MerchantDisableResponse "{"status":true}"
// @Router /merchants/disable [post]
func (r *MerchantDisableRequest) Exec(ctx context.Context) interface{} {
	resp := MerchantDisableResponse{}

	code, err := svc.DisableMerchant(ctx, r.MerchantID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
