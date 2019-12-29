package server

import (
	"context"
	"errors"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// ActivityDetailRequest .
type ActivityDetailRequest struct {
	wsgin.MustAuthRequest

	ActivityID uint64 `json:"activity_id" form:"activity_id"` // 活动id
	Name       string `json:"name" form:"name"`               // 活动名称
}

// ActivityDetailResponse .
type ActivityDetailResponse struct {
	wsgin.BaseResponse

	Data *model.Activity `json:"data"`
}

// New .
func (r *ActivityDetailRequest) New() wsgin.Process {
	return &ActivityDetailRequest{}
}

// Extract .
func (r *ActivityDetailRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	code, err = r.DefaultExtract(r, c)
	if err != nil {
		return
	}
	if r.ActivityID == 0 && r.Name == "" {
		return wsgin.APICodeInvalidParame, errors.New("活动DID或活动名称至少传入一个")
	}
	return
}

// Exec 获取活动详情
// @Summary 获取活动详情
// @Description get activity detail
// @Tags 活动
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id query string false "活动ID"
// @Param name query string false "活动名称"
// @Success 200 {object} server.ActivityDetailResponse "{"status":true}"
// @Router /activitys/detail [get]
func (r *ActivityDetailRequest) Exec(ctx context.Context) interface{} {
	resp := ActivityDetailResponse{}
	data, code, err := svc.DetailActivity(ctx, r.ActivityID, r.Name)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
