package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// ActivityParticipantListRequest .
type ActivityParticipantListRequest struct {
	wsgin.AuthPagingRequest

	ActivityID  uint64 `json:"activity_id" form:"activity_id" binding:"required"` // 活动ID
	Mobile      string `json:"mobile" form:"mobile"`                              // 参与者手机号
	IsSearchWin bool   `json:"is_search_win" form:"is_search_win"`                // 是否查询中奖名单列表
}

// ActivityParticipantListResponse .
type ActivityParticipantListResponse struct {
	wsgin.BasePagingResponse

	Data []*model.LuckyNumberRecord `json:"data"`
}

// New .
func (r *ActivityParticipantListRequest) New() wsgin.Process {
	return &ActivityParticipantListRequest{}
}

// Extract .
func (r *ActivityParticipantListRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取活动参与者信息列表
// @Summary 获取活动参与者信息列表
// @Description get activity participant list
// @Security ApiKeyAuth
// @Tags 活动
// @Accept json
// @Produce json
// @Param activity_id query string true "活动ID"
// @Param mobile query string false "参与者手机号"
// @Param is_search_win query bool false "是否查询中奖名单列表"
// @Param page_no query int false "页码" default(1)
// @Param page_size query int false "页数" default(10)
// @Success 200 {object} server.ActivityParticipantListResponse	"{"status":true}"
// @Router /activitys/participant [get]
func (r *ActivityParticipantListRequest) Exec(ctx context.Context) interface{} {
	resp := ActivityParticipantListResponse{}

	data, total, code, err := svc.ListActivityParticipant(ctx, &model.ActivityParticipantListVO{
		ActivityID:  r.ActivityID,
		Mobile:      r.Mobile,
		IsSearchWin: r.IsSearchWin,
		PageNo:      r.PageNo,
		PageSize:    r.PageSize,
	})
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	resp.Total = total
	resp.Current = r.PageNo
	return resp
}
