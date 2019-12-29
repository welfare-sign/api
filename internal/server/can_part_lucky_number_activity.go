package server

import (
	"context"

	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// CanPartLuckyNumberActivityRequest .
type CanPartLuckyNumberActivityRequest struct {
	wsgin.MustAuthRequest

	ActivityID uint64 `form:"activity_id" json:"activity_id"` // 活动id
}

// CanPartLuckyNumberActivityResponse .
type CanPartLuckyNumberActivityResponse struct {
	wsgin.BaseResponse

	Data bool `json:"data"`
}

// New .
func (r *CanPartLuckyNumberActivityRequest) New() wsgin.Process {
	return &CanPartLuckyNumberActivityRequest{}
}

// Extract .
func (r *CanPartLuckyNumberActivityRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec .
// @Summary 用户是否可以参与猜数字活动
// @Description can participate in the lucky number activity
// @Tags 客户
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id query string true "活动ID"
// @Success 200 {object} server.CanPartLuckyNumberActivityResponse "{"status":true}"
// @Router /customers/can_part_lucky_number_activity [get]
func (r *CanPartLuckyNumberActivityRequest) Exec(ctx context.Context) interface{} {
	resp := CanPartLuckyNumberActivityResponse{}

	data, code, err := svc.CanPartLuckyNumberActivity(ctx, r.TokenParames.UID, r.ActivityID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
