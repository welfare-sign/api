package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/pkg/wsgin"
)

// CustomerDelRequest 删除客户
type CustomerDelRequest struct {
	wsgin.MustAuthRequest

	CustomerID uint64 `form:"customer_id" json:"customer_id"` // 客户ID
}

// CustomerDelResponse .
type CustomerDelResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *CustomerDelRequest) New() wsgin.Process {
	return &CustomerDelRequest{}
}

// Extract .
func (r *CustomerDelRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 删除客户
// @Summary 删除客户
// @Description delete customer
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Param args body server.CustomerDelRequest true "参数"
// @Success 200 {object} server.CustomerDelResponse "{"status":true}"
// @Router /customers [delete]
func (r *CustomerDelRequest) Exec(ctx context.Context) interface{} {
	resp := CustomerDelResponse{}

	code, err := svc.DeleteCustomer(ctx, r.CustomerID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
