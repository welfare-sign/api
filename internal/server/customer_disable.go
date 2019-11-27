package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/pkg/wsgin"
)

// CustomerDisableRequest 禁用客户
type CustomerDisableRequest struct {
	wsgin.MustAuthRequest

	CustomerID uint64 `form:"customer_id" json:"customer_id"` // 客户ID
}

// CustomerDisableResponse .
type CustomerDisableResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *CustomerDisableRequest) New() wsgin.Process {
	return &CustomerDisableRequest{}
}

// Extract .
func (r *CustomerDisableRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 禁用客户
// @Summary 禁用客户
// @Description disable customer
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Param args body server.CustomerDisableRequest true "参数"
// @Success 200 {object} server.CustomerDisableResponse "{"status":true}"
// @Router /customers/disable [post]
func (r *CustomerDisableRequest) Exec(ctx context.Context) interface{} {
	resp := CustomerDisableResponse{}

	code, err := svc.DisableCustomer(ctx, r.CustomerID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
