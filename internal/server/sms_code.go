package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/pkg/wsgin"
)

// CodeRequest .
type CodeRequest struct {
	wsgin.BaseRequest

	Mobile string `json:"mobile" form:"mobile" binding:"required,mobile"`
}

// CodeResponse .
type CodeResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *CodeRequest) New() wsgin.Process {
	return &CodeRequest{}
}

// Extract .
func (r *CodeRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 发送验证码
// @Summary 发送验证码
// @Description send sms code
// @Tags 验证码
// @Accept json
// @Produce json
// @Param mobile query string false "手机号"
// @Success 200 {object} server.CodeResponse "{"status":true}"
// @Router /verify_code [get]
func (r *CodeRequest) Exec(ctx context.Context) interface{} {
	resp := CodeResponse{}

	code, err := svc.SendVerifyCode(ctx, r.Mobile)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
