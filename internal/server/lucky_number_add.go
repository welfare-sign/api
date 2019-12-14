package server

import (
	"context"

	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// LuckyNumberAddRequest .
type LuckyNumberAddRequest struct {
	wsgin.MustAuthRequest

	Num int64 `json:"num" binding:"required"` // 猜的数字
}

// LuckyNumberAddResponse .
type LuckyNumberAddResponse struct {
	wsgin.BaseResponse

	Data []int64 `json:"data"`
}

// New .
func (r *LuckyNumberAddRequest) New() wsgin.Process {
	return &LuckyNumberAddRequest{}
}

// Extract .
func (r *LuckyNumberAddRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 写入用户猜的数字
// @Summary 写入用户猜的数字
// @Description post customer lucky number
// @Tags 客户
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param args body server.LuckyNumberAddRequest true "参数"
// @Success 200 {object} server.LuckyNumberAddResponse "{"status":true}"
// @Router /customers/lucky_number [post]
func (r *LuckyNumberAddRequest) Exec(ctx context.Context) interface{} {
	resp := LuckyNumberAddResponse{}
	data, code, err := svc.AddLuckyNumber(ctx, r.TokenParames.UID, r.Num)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
