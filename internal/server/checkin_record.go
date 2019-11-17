package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// CheckinRecordRequest .
type CheckinRecordRequest struct {
	wsgin.AuthRequest

	CustomerID uint64 `form:"customer_id" json:"customer_id"` // 当该参数有值时以该参数为准
}

// CheckinRecordResponse .
type CheckinRecordResponse struct {
	wsgin.BaseResponse

	Data []*model.CheckinRecord `json:"data"`
}

// New .
func (r *CheckinRecordRequest) New() wsgin.Process {
	return &CheckinRecordRequest{}
}

// Extract .
func (r *CheckinRecordRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取签到记录
// @Summary 获取签到记录
// @Description get customer checkin record
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Param customer_id query string false "当该参数有值时以该参数为准，否则以token用户为准"
// @Success 200 {object} server.CheckinRecordResponse	"{"status":true}"
// @Router /customers/checkin_record [get]
func (r *CheckinRecordRequest) Exec(ctx context.Context) interface{} {
	resp := CheckinRecordResponse{}

	data, code, err := svc.GetCustomerCheckinRecord(ctx, r.TokenParames.UID, r.CustomerID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
