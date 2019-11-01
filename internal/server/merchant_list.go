package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// MerchantListRequest .
type MerchantListRequest struct {
	wsgin.AuthPagingRequest

	StoreName    string `json:"store_name" form:"store_name" example:"商户名"`
	ContactName  string `json:"contact_name" form:"contact_name" example:"联系人"`
	ContactPhone string `json:"contact_phone" form:"contact_phone" binding:"omitempty,mobile" example:"联系电话"`
}

// MerchantListResponse .
type MerchantListResponse struct {
	wsgin.BaseResponse

	Data []*model.Merchant `json:"data"`
}

// New .
func (r *MerchantListRequest) New() wsgin.Process {
	return &MerchantListRequest{}
}

// Extract .
func (r *MerchantListRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取商户列表
// @Summary 获取商户列表
// @Description get merchant list
// @Security ApiKeyAuth
// @Tags 商户
// @Accept json
// @Produce json
// @Param store_name query string false "商户名"
// @Param contact_name query string false "联系人"
// @Param contact_phone query string false "联系电话"
// @Param page_no query int false "页码" default(1)
// @Param page_size query int false "页数" default(10)
// @Success 200 {object} server.MerchantListResponse	"{"status":true}"
// @Router /merchants [get]
func (r *MerchantListRequest) Exec(ctx context.Context) interface{} {
	resp := MerchantListResponse{}

	data, code, err := svc.GetMerchantList(ctx, &model.MerchantListVO{
		StoreName:    r.StoreName,
		ContactName:  r.ContactName,
		ContactPhone: r.ContactPhone,
		PageNo:       r.PageNo,
		PageSize:     r.PageSize,
	})
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
