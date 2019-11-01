package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// ExecWriteOffRequest .
type ExecWriteOffRequest struct {
	wsgin.AuthRequest

	CustomerID uint64 `form:"customer_id" json:"customer_id" binding:"required" example:"客户ID"`
	Num        uint64 `form:"num" json:"num" binding:"required" example:"核销数目"`
}

// ExecWriteOffResponse .
type ExecWriteOffResponse struct {
	wsgin.BaseResponse

	Data *model.MerchantWriteOffRespVO `json:"data"`
}

// New .
func (r *ExecWriteOffRequest) New() wsgin.Process {
	return &ExecWriteOffRequest{}
}

// Extract .
func (r *ExecWriteOffRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 商户核销
// @Summary 商户核销
// @Description merchant exec write off
// @Security ApiKeyAuth
// @Tags 商户
// @Accept json
// @Produce json
// @Param args body server.ExecWriteOffRequest true "参数"
// @Success 200 {object} server.ExecWriteOffResponse	"{"status":true}"
// @Router /merchants/writeoff [post]
func (r *ExecWriteOffRequest) Exec(ctx context.Context) interface{} {
	resp := ExecWriteOffResponse{}

	data, code, err := svc.ExecWriteOff(ctx, &model.MerchantExecWriteOffVO{
		MerchantID: r.TokenParames.Uid,
		CustomerID: r.CustomerID,
		Num:        r.Num,
	})
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
