package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// MerchantDetailRequest .
type MerchantDetailRequest struct {
	wsgin.AuthRequest

	MerchantID uint64 `json:"merchant_id" form:"merchant_id" example:"商户ID"`
}

// MerchantDetailResponse .
type MerchantDetailResponse struct {
	wsgin.BaseResponse

	Data *model.Merchant `json:"data"`
}

// New .
func (r *MerchantDetailRequest) New() wsgin.Process {
	return &MerchantDetailRequest{}
}

// Extract .
func (r *MerchantDetailRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取商户详情信息
// @Summary 获取商户详情信息
// @Description get merchant detail
// @Security ApiKeyAuth
// @Tags 商户
// @Accept json
// @Produce json
// @Param merchant_id query string false "商户ID,商户访问时可不传"
// @Success 200 {object} server.MerchantDetailResponse "{"status":true}"
// @Router /merchants/detail [get]
func (r *MerchantDetailRequest) Exec(ctx context.Context) interface{} {
	resp := MerchantDetailResponse{}

	merchantID := r.TokenParames.UID
	if r.MerchantID != 0 {
		merchantID = r.MerchantID
	}
	data, code, err := svc.GetMerchantDetailBySelfAccessToken(ctx, merchantID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
