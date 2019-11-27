package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// MerchantAddRequest .
type MerchantAddRequest struct {
	wsgin.MustAuthRequest

	Merchant *model.MerchantVO `json:"merchant" binding:"required,dive"`
}

// MerchantAddResponse .
type MerchantAddResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *MerchantAddRequest) New() wsgin.Process {
	return &MerchantAddRequest{}
}

// Extract .
func (r *MerchantAddRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 新增商户
// @Summary 新增商户
// @Description create new merchant
// @Security ApiKeyAuth
// @Tags 商户
// @Accept json
// @Produce json
// @Param args body server.MerchantAddRequest true "参数"
// @Success 200 {object} server.MerchantAddResponse "{"status":true}"
// @Router /merchants [post]
func (r *MerchantAddRequest) Exec(ctx context.Context) interface{} {
	resp := MerchantAddResponse{}

	code, err := svc.AddMerchant(ctx, r.Merchant)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
