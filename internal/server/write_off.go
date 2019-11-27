package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// WriteOffRequest .
type WriteOffRequest struct {
	wsgin.MustAuthRequest

	CustomerID uint64 `json:"customer_id" form:"customer_id" binding:"required" example:"客户ID"`
}

// WriteOffResponse .
type WriteOffResponse struct {
	wsgin.BaseResponse

	Data *model.MerchantWriteOffRespVO `json:"data"`
}

// New .
func (r *WriteOffRequest) New() wsgin.Process {
	return &WriteOffRequest{}
}

// Extract .
func (r *WriteOffRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 商家获取核销信息
// @Summary 商户获取核销信息
// @Description merchant get customer write off
// @Security ApiKeyAuth
// @Tags 商户
// @Accept json
// @Produce json
// @Param customer_id query int true "客户ID"
// @Success 200 {object} server.WriteOffResponse	"{"status":true}"
// @Router /merchants/writeoff [get]
func (r *WriteOffRequest) Exec(ctx context.Context) interface{} {
	resp := WriteOffResponse{}

	data, code, err := svc.GetWriteOff(ctx, r.TokenParames.UID, r.CustomerID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
