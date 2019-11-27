package server

import (
	"context"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// IssueRecordRequest .
type IssueRecordRequest struct {
	wsgin.MustAuthRequest
}

// IssueRecordResponse .
type IssueRecordResponse struct {
	wsgin.BaseResponse

	Data []*model.IssueRecord `json:"data"`
}

// New .
func (r *IssueRecordRequest) New() wsgin.Process {
	return &IssueRecordRequest{}
}

// Extract .
func (r *IssueRecordRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 客户查看我的福利
// @Summary 客户查看我的福利
// @Description customer get issue record
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Success 200 {object} server.IssueRecordResponse	"{"status":true}"
// @Router /customers/issue_records [get]
func (r *IssueRecordRequest) Exec(ctx context.Context) interface{} {
	resp := IssueRecordResponse{}

	data, code, err := svc.GetIssueRecords(ctx, r.TokenParames.UID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
