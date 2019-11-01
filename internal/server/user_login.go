package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// UserLoginRequest .
type UserLoginRequest struct {
	wsgin.BaseRequest

	Name     string `json:"name" form:"name" binding:"required"`         // 用户名
	Password string `json:"password" form:"password" binding:"required"` // 密码
}

// UserLoginResponse .
type UserLoginResponse struct {
	wsgin.BaseResponse

	Data string `json:"data"` // access token
}

// New .
func (r *UserLoginRequest) New() wsgin.Process {
	return &UserLoginRequest{}
}

// Extract .
func (r *UserLoginRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 后台用户登录
// @Summary 后台用户登录
// @Description user login
// @Tags 后台用户
// @Accept json
// @Produce json
// @Param args body server.UserLoginRequest true "参数"
// @Success 200 {object} server.UserLoginResponse	"{"status":true}"
// @Router /users/login [post]
func (r *UserLoginRequest) Exec(ctx context.Context) interface{} {
	resp := UserLoginResponse{}

	data, code, err := svc.UserLogin(ctx, &model.UserVO{
		Name:     r.Name,
		Password: r.Password,
	})
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
