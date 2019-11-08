package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// CustomerDetailRequest .
type CustomerDetailRequest struct {
	wsgin.AuthRequest

	CustomerID uint64 `json:"customer_id" form:"customer_id" example:"客户id"`
}

// CustomerDetailResponse .
type CustomerDetailResponse struct {
	wsgin.BaseResponse

	Data *model.Customer `json:"data"`
}

// New .
func (r *CustomerDetailRequest) New() wsgin.Process {
	return &CustomerDetailRequest{}
}

// Extract .
func (r *CustomerDetailRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取客户详情
// @Summary 获取客户详情
// @Description get customer detail
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Param customer_id query string false "客户ID,客户访问时可不传"
// @Success 200 {object} server.CustomerDetailResponse "{"status":true}"
// @Router /customers/detail [get]
func (r *CustomerDetailRequest) Exec(ctx context.Context) interface{} {
	resp := CustomerDetailResponse{}

	customerID := r.TokenParames.UID
	if r.CustomerID != 0 {
		customerID = r.CustomerID
	}
	data, code, err := svc.GetCustomerDetail(ctx, customerID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
