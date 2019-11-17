package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/pkg/wsgin"
)

// WXPayRequest .
type WXPayRequest struct {
	wsgin.AuthRequest
}

// WXPayResponse .
type WXPayResponse struct {
	wsgin.BaseResponse

	Data string `json:"data"`
}

// New .
func (r *WXPayRequest) New() wsgin.Process {
	return &WXPayRequest{}
}

// Extract .
func (r *WXPayRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 用户支付
// @Summary 用户支付
// @Description customer pay
// @Security ApiKeyAuth
// @Tags 微信
// @Accept json
// @Produce json
// @Success 200 {object} server.WXPayResponse	"{"status":true}"
// @Router /wx/pay [post]
func (r *WXPayRequest) Exec(ctx context.Context) interface{} {
	resp := WXPayResponse{}

	param, code, err := svc.WXPay(ctx, r.TokenParames.UID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = param
	return resp
}
