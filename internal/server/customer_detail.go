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

	CustomerID        uint64 `json:"customer_id" form:"customer_id"`                   // 客户id
	IsHelpCheckinPage bool   `form:"is_help_checkin_page" json:"is_help_checkin_page"` // 是否是帮签页面
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

	data, code, err := svc.GetCustomerDetail(ctx, r.TokenParames.UID, r.CustomerID, r.IsHelpCheckinPage)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
