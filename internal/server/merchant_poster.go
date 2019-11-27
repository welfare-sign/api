package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"welfare-sign/internal/pkg/wsgin"
)

// MerchantPosterRequest .
type MerchantPosterRequest struct {
	wsgin.MustAuthRequest

	Response gin.ResponseWriter
}

// New .
func (r *MerchantPosterRequest) New() wsgin.Process {
	return &MerchantPosterRequest{}
}

// Extract .
func (r *MerchantPosterRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	r.Response = c.Writer
	return r.DefaultExtract(r, c)
}

// Exec 获取商户随机一张海报
// @Summary 获取商户随机一张海报
// @Description get merchant poster
// @Security ApiKeyAuth
// @Tags 商户
// @Accept json
// @Produce json
// @Success 200 {string} string "{"status":true}"
// @Router /merchants/poster [get]
func (r *MerchantPosterRequest) Exec(ctx context.Context) interface{} {
	data, err := svc.GetRoundMerchantPoster(ctx)
	if err == nil {
		r.Response.Write(data)
	}
	return nil
}
