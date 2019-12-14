package server

import (
	"context"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// LuckyNumberBeforeRequest .
type LuckyNumberBeforeRequest struct {
	wsgin.MustAuthRequest
}

// LuckyNumberBeforeResponse .
type LuckyNumberBeforeResponse struct {
	wsgin.BaseResponse

	Data *model.LuckyNumberRecord `json:"data"`
}

// New .
func (r *LuckyNumberBeforeRequest) New() wsgin.Process {
	return &LuckyNumberBeforeRequest{}
}

// Extract .
func (r *LuckyNumberBeforeRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 用户上期猜的数字
// @Summary 用户上期猜的数字
// @Description get customer lucky number before
// @Tags 客户
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} server.LuckyNumberBeforeResponse "{"status":true}"
// @Router /customers/lucky_number/before [get]
func (r *LuckyNumberBeforeRequest) Exec(ctx context.Context) interface{} {
	resp := LuckyNumberBeforeResponse{}
	data, code, err := svc.GetLuckyNumberBefore(ctx, r.TokenParames.UID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
