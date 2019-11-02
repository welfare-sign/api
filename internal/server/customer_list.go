package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// CustomerListRequest .
type CustomerListRequest struct {
	wsgin.AuthPagingRequest

	Name   string `json:"name" form:"name" example:"用户名"`
	Mobile string `json:"mobile" form:"mobile" binding:"mobile" example:"联系电话"`
}

// CustomerListResponse .
type CustomerListResponse struct {
	wsgin.BasePagingResponse

	Data []*model.Customer `json:"data"`
}

// New .
func (r *CustomerListRequest) New() wsgin.Process {
	return &CustomerListRequest{}
}

// Extract .
func (r *CustomerListRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取客户列表
// @Summary 获取客户列表
// @Description get customer list
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Param name query string false "用户名"
// @Param mobile query string false "联系电话"
// @Param page_no query int false "页码" default(1)
// @Param page_size query int false "页数" default(10)
// @Success 200 {object} server.CustomerListResponse	"{"status":true}"
// @Router /customers [get]
func (r *CustomerListRequest) Exec(ctx context.Context) interface{} {
	resp := CustomerListResponse{}

	data, total, code, err := svc.GetCustomerList(ctx, &model.CustomerListVO{
		Name:     r.Name,
		Mobile:   r.Mobile,
		PageNo:   r.PageNo,
		PageSize: r.PageSize,
	})
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	resp.Total = total
	resp.Current = r.PageNo
	return resp
}
