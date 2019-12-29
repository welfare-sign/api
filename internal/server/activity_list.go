package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// ActivityListRequest .
type ActivityListRequest struct {
	wsgin.AuthPagingRequest

	Name string `json:"name" form:"name"` // 活动名
}

// ActivityListResponse .
type ActivityListResponse struct {
	wsgin.BasePagingResponse

	Data []*model.Activity `json:"data"`
}

// New .
func (r *ActivityListRequest) New() wsgin.Process {
	return &ActivityListRequest{}
}

// Extract .
func (r *ActivityListRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取活动列表
// @Summary 获取活动列表
// @Description get activity list
// @Security ApiKeyAuth
// @Tags 活动
// @Accept json
// @Produce json
// @Param name query string false "活动名"
// @Param page_no query int false "页码" default(1)
// @Param page_size query int false "页数" default(10)
// @Success 200 {object} server.ActivityListResponse	"{"status":true}"
// @Router /activitys [get]
func (r *ActivityListRequest) Exec(ctx context.Context) interface{} {
	resp := ActivityListResponse{}

	data, total, code, err := svc.ListActivity(ctx, &model.ActivityListVO{
		Name:     r.Name,
		PageNo:   r.PageNo,
		PageSize: r.PageSize,
	})
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	resp.Total = total
	resp.Current = r.PageNo
	return resp
}
