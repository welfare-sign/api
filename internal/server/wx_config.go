package server

import (
	"context"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// WXConfigRequest .
type WXConfigRequest struct {
	wsgin.BaseRequest

	URL string `form:"url" json:"url" binding:"required"` // url必须是调用JS接口页面的完整URL, 不包含#及其后面部分
}

// WXConfigResponse .
type WXConfigResponse struct {
	wsgin.BaseResponse

	Data *model.WXConfigResp `json:"data"`
}

// New .
func (r *WXConfigRequest) New() wsgin.Process {
	return &WXConfigRequest{}
}

// Extract .
func (r *WXConfigRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取微信接口配置
// @Summary 获取微信接口配置
// @Description get wx config
// @Tags 微信
// @Accept json
// @Produce json
// @Param url query string true "当前页URL，不包含#及其后面部分"
// @Success 200 {object} server.WXConfigResponse "{"status":true}"
// @Router /wx/config [get]
func (r *WXConfigRequest) Exec(ctx context.Context) interface{} {
	resp := WXConfigResponse{}

	data, code, err := svc.GetWXConfig(ctx, r.URL)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
