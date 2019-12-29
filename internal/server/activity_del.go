package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/pkg/wsgin"
)

// ActivityDelRequest 删除活动
type ActivityDelRequest struct {
	wsgin.MustAuthRequest

	ActivityID uint64 `json:"activity_id" form:"activity_id"` // 活动id
}

// ActivityDelResponse .
type ActivityDelResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *ActivityDelRequest) New() wsgin.Process {
	return &ActivityDelRequest{}
}

// Extract .
func (r *ActivityDelRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 删除活动
// @Summary 删除活动
// @Description delete activity
// @Security ApiKeyAuth
// @Tags 活动
// @Accept json
// @Produce json
// @Param args body server.ActivityDelRequest true "参数"
// @Success 200 {object} server.ActivityDelResponse "{"status":true}"
// @Router /activitys [delete]
func (r *ActivityDelRequest) Exec(ctx context.Context) interface{} {
	resp := ActivityDelResponse{}

	code, err := svc.DelActivity(ctx, r.ActivityID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
