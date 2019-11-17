package server

import (
	"context"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// NearMerchantRequest 获取用户附近最近的几家店铺
type NearMerchantRequest struct {
	wsgin.AuthRequest

	Lon      float64 `form:"lon" json:"lon" binding:"required"` // 经度
	Lat      float64 `form:"lat" json:"lat" binding:"required"` // 维度
	Distince float64 `form:"distince" json:"distince"`          // 距离多少公里内，默认10
	Num      int     `form:"num" json:"num"`                    // 返回数量，默认4个
}

// NearMerchantResponse .
type NearMerchantResponse struct {
	wsgin.BaseResponse

	Data []*model.Merchant `json:"data"`
}

// New .
func (r *NearMerchantRequest) New() wsgin.Process {
	return &NearMerchantRequest{}
}

// Extract .
func (r *NearMerchantRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取用户附近最近的几家店铺
// @Summary 获取用户附近最近的几家店铺
// @Description get customer near merchant
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Param lon query int true "经度"
// @Param lat query int true "维度"
// @Param distince query int false "距离多少公里内，默认10"
// @Param num query int false "返回数量，默认4个"
// @Success 200 {object} server.NearMerchantResponse	"{"status":true}"
// @Router /customers/near_merchant [get]
func (r *NearMerchantRequest) Exec(ctx context.Context) interface{} {
	resp := NearMerchantResponse{}
	if r.Distince == 0 {
		r.Distince = 10
	}
	if r.Num == 0 {
		r.Num = 4
	}
	data, code, err := svc.CustomerNearMerchant(ctx, &model.NearMerchantVO{
		Lon:      r.Lon,
		Lat:      r.Lat,
		Distince: r.Distince,
		Num:      r.Num,
	})
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
