package server

import (
	"context"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// ActivityCurrentlyAvailableRequest 当前用户可参与的活动
type ActivityCurrentlyAvailableRequest struct {
	wsgin.BaseRequest
}

// ActivityCurrentlyAvailableResponse .
type ActivityCurrentlyAvailableResponse struct {
	wsgin.BaseResponse

	Data *model.Activity `json:"data"`
}

// New .
func (r *ActivityCurrentlyAvailableRequest) New() wsgin.Process {
	return &ActivityCurrentlyAvailableRequest{}
}

// Extract .
func (r *ActivityCurrentlyAvailableRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取当前可参与的活动
// @Summary 获取当前可参与的活动
// @Description get currently available activity detail
// @Tags 活动
// @Accept json
// @Produce json
// @Success 200 {object} server.ActivityCurrentlyAvailableResponse "{"status":true}"
// @Router /activitys/currently_available [get]
func (r *ActivityCurrentlyAvailableRequest) Exec(ctx context.Context) interface{} {
	resp := ActivityCurrentlyAvailableResponse{}
	data, code, err := svc.CurrentlyAvailableActivity(ctx)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
