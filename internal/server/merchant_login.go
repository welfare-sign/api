package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// MerchantLoginRequest .
type MerchantLoginRequest struct {
	wsgin.BaseRequest

	ContactPhone string `json:"contact_phone" form:"contact_phone" binding:"required,mobile" example:"手机号"`
	Code         string `json:"code" form:"code" binding:"required" example:"验证码"`
}

// MerchantLoginResponse .
type MerchantLoginResponse struct {
	wsgin.BaseResponse

	Data string `json:"data"` // access token
}

// New .
func (r *MerchantLoginRequest) New() wsgin.Process {
	return &MerchantLoginRequest{}
}

// Extract .
func (r *MerchantLoginRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 商户登录
// @Summary 商户登录
// @Description merchant login
// @Tags 商户
// @Accept json
// @Produce json
// @Param args body model.MerchantLoginVO true "参数"
// @Success 200 {object} server.MerchantLoginResponse "{"status":true}"
// @Router /merchants/login [post]
func (r *MerchantLoginRequest) Exec(ctx context.Context) interface{} {
	resp := MerchantLoginResponse{}

	data, code, err := svc.MerchantLogin(ctx, &model.MerchantLoginVO{
		ContactPhone: r.ContactPhone,
		Code:         r.Code,
	})
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
