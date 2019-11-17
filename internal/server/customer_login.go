package server

import (
	"context"

	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// CustomerLoginRequest .
type CustomerLoginRequest struct {
	wsgin.BaseRequest

	Code string `json:"code" form:"code" binding:"required"`
}

// CustomerLoginResponse .
type CustomerLoginResponse struct {
	wsgin.BaseResponse

	Data string `json:"data"`
}

// New .
func (r *CustomerLoginRequest) New() wsgin.Process {
	return &CustomerLoginRequest{}
}

// Extract .
func (r *CustomerLoginRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 客户登录
// @Summary 客户登录
// @Description customer login
// @Tags 客户
// @Accept json
// @Produce json
// @Param code query string true "微信回调code"
// @Success 200 {object} server.CustomerLoginResponse "{"status":true}"
// @Router /customers/login [post]
func (r *CustomerLoginRequest) Exec(ctx context.Context) interface{} {
	resp := CustomerLoginResponse{}

	data, code, err := svc.CustomerLogin(ctx, r.Code)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
