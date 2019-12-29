package server

import (
	"context"

	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// ActivityAllPrizeIssuedRequest .
type ActivityAllPrizeIssuedRequest struct {
	wsgin.BaseRequest
}

// ActivityAllPrizeIssuedResponse .
type ActivityAllPrizeIssuedResponse struct {
	wsgin.BaseResponse

	Data int `json:"data"`
}

// New .
func (r *ActivityAllPrizeIssuedRequest) New() wsgin.Process {
	return &ActivityAllPrizeIssuedRequest{}
}

// Extract .
func (r *ActivityAllPrizeIssuedRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec .
// @Summary 获取所有活动已发放的奖品总数
// @Description get activity all prize issued
// @Tags 活动
// @Accept json
// @Produce json
// @Success 200 {object} server.ActivityAllPrizeIssuedResponse "{"status":true}"
// @Router /activitys/all_prize_issued [get]
func (r *ActivityAllPrizeIssuedRequest) Exec(ctx context.Context) interface{} {
	resp := ActivityAllPrizeIssuedResponse{}
	data, code, err := svc.ActivityAllPrizeIssued(ctx)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
