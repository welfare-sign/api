package server

import (
	"context"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// ActivityAddPrizeNumberRequest .
type ActivityAddPrizeNumberRequest struct {
	wsgin.MustAuthRequest

	ActivityID  uint64 `json:"activity_id"`                     // 活动id
	PrizeNumber uint64 `json:"prize_number" binding:"required"` // 中奖号码
}

// ActivityAddPrizeNumberResponse .
type ActivityAddPrizeNumberResponse struct {
	wsgin.BaseResponse

	Data *model.Activity `json:"data"`
}

// New .
func (r *ActivityAddPrizeNumberRequest) New() wsgin.Process {
	return &ActivityAddPrizeNumberRequest{}
}

// Extract .
func (r *ActivityAddPrizeNumberRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec .
// @Summary 开奖
// @Description Draw
// @Tags 活动
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param args body server.ActivityAddPrizeNumberRequest true "参数"
// @Success 200 {object} server.ActivityAddPrizeNumberResponse "{"status":true}"
// @Router /activitys/draw [post]
func (r *ActivityAddPrizeNumberRequest) Exec(ctx context.Context) interface{} {
	resp := ActivityAddPrizeNumberResponse{}
	data, code, err := svc.DrawActivity(ctx, r.ActivityID, r.PrizeNumber)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
