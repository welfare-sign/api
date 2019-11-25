package server

import (
	"context"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// ExecIssueRecordRequest .
type ExecIssueRecordRequest struct {
	wsgin.AuthRequest

	MerchantID uint64 `json:"merchant_id" binding:"required"` // 店铺ID
	Mobile     string `json:"mobile"`                         // 手机号
	Code       string `json:"code"`                           // 验证码
}

// ExecIssueRecordResponse .
type ExecIssueRecordResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *ExecIssueRecordRequest) New() wsgin.Process {
	return &ExecIssueRecordRequest{}
}

// Extract .
func (r *ExecIssueRecordRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 用户领取福利
// @Summary 用户领取福利
// @Description customer exec issue record
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Param args body server.ExecIssueRecordRequest true "参数"
// @Success 200 {object} server.ExecIssueRecordResponse	"{"status":true}"
// @Router /customers/issue_records [post]
func (r *ExecIssueRecordRequest) Exec(ctx context.Context) interface{} {
	resp := ExecIssueRecordResponse{}

	code, err := svc.ExecIssueRecords(ctx, r.TokenParames.UID, r.MerchantID, r.Mobile, r.Code)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
