package server

import (
	"context"
	"errors"
	"time"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// ActivityAddRequest .
type ActivityAddRequest struct {
	wsgin.MustAuthRequest

	Name        string    `json:"name" binding:"required"`         // 活动名称
	StartTime   time.Time `json:"start_time" binding:"required"`   // 开始时间
	EndTime     time.Time `json:"end_time" binding:"required"`     // 结束时间
	Scope       string    `json:"scope"  binding:"required"`       // 参与方式：O：无限制；R：必须在福利签领取过签到奖励
	PrizeAmount uint64    `json:"prize_amount" binding:"required"` // 奖品数量
	Poster      string    `json:"poster" binding:"required"`       // 活动海报
}

// ActivityAddResponse .
type ActivityAddResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *ActivityAddRequest) New() wsgin.Process {
	return &ActivityAddRequest{}
}

// Extract .
func (r *ActivityAddRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	code, err = r.DefaultExtract(r, c)
	if err != nil {
		return
	}
	if r.StartTime.After(r.EndTime) || r.StartTime.Equal(r.EndTime) {
		return wsgin.APICodeInvalidParame, errors.New("开始时间必须小于结束时间")
	}
	return
}

// Exec .
// @Summary 添加或者更新活动信息
// @Description post or update activity
// @Tags 活动
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param args body server.ActivityAddRequest true "参数"
// @Success 200 {object} server.ActivityAddResponse "{"status":true}"
// @Router /activitys [put]
func (r *ActivityAddRequest) Exec(ctx context.Context) interface{} {
	resp := ActivityAddResponse{}
	code, err := svc.UpsertActivity(ctx, &model.ActivityVO{
		Name:        r.Name,
		StartTime:   r.StartTime,
		EndTime:     r.EndTime,
		Scope:       r.Scope,
		PrizeAmount: r.PrizeAmount,
		Poster:      r.Poster,
	})
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
