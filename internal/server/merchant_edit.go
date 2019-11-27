package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// MerchantEditRequest 编辑商户
type MerchantEditRequest struct {
	wsgin.MustAuthRequest

	Merchant   *model.MerchantVO `json:"merchant" binding:"required,dive"` // 商户信息
	MerchantID uint64            `json:"merchant_id"`                      // 商户ID
}

// MerchantEditResponse .
type MerchantEditResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *MerchantEditRequest) New() wsgin.Process {
	return &MerchantEditRequest{}
}

// Extract .
func (r *MerchantEditRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 编辑商户
// @Summary 编辑商户
// @Description edit merchant
// @Security ApiKeyAuth
// @Tags 商户
// @Accept json
// @Produce json
// @Param args body server.MerchantEditRequest true "参数"
// @Success 200 {object} server.MerchantEditResponse "{"status":true}"
// @Router /merchants [put]
func (r *MerchantEditRequest) Exec(ctx context.Context) interface{} {
	resp := MerchantEditResponse{}

	code, err := svc.EditMerchant(ctx, r.MerchantID, r.Merchant)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
