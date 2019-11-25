package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// CheckinRecordListRequest .
type CheckinRecordListRequest struct {
	wsgin.AuthRequest

	CustomerID uint64 `form:"customer_id" json:"customer_id"` // 当该参数有值时以该参数为准
}

// CheckinRecordListResponse .
type CheckinRecordListResponse struct {
	wsgin.BaseResponse

	Data []*model.CheckinRecordListResp `json:"data"`
}

// New .
func (r *CheckinRecordListRequest) New() wsgin.Process {
	return &CheckinRecordListRequest{}
}

// Extract .
func (r *CheckinRecordListRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取签到记录列表(临时)
// @Summary 获取签到记录列表(临时)
// @Description get checkin record list
// @Security ApiKeyAuth
// @Tags 后台用户
// @Accept json
// @Produce json
// @Success 200 {object} server.CheckinRecordListResponse	"{"status":true}"
// @Router /users/checkin_record_list [get]
func (r *CheckinRecordListRequest) Exec(ctx context.Context) interface{} {
	resp := CheckinRecordListResponse{}

	data, code, err := svc.GetAllCustomerCheckinRecordList(ctx)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
