package server

import (
	"context"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// LuckyNumberDetailRequest .
type LuckyNumberDetailRequest struct {
	wsgin.MustAuthRequest
}

// LuckyNumberDetailResponse .
type LuckyNumberDetailResponse struct {
	wsgin.BaseResponse

	Data *model.LuckyNumberRecord `json:"data"`
}

// New .
func (r *LuckyNumberDetailRequest) New() wsgin.Process {
	return &LuckyNumberDetailRequest{}
}

// Extract .
func (r *LuckyNumberDetailRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取用户猜的数字
// @Summary 获取用户猜的数字
// @Description get customer lucky number
// @Tags 客户
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} server.LuckyNumberDetailResponse "{"status":true}"
// @Router /customers/lucky_number [get]
func (r *LuckyNumberDetailRequest) Exec(ctx context.Context) interface{} {
	resp := LuckyNumberDetailResponse{}
	data, code, err := svc.GetLuckyNumberDetail(ctx, r.TokenParames.UID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
