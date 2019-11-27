package server

import (
	"context"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// QRCodeRequest .
type QRCodeRequest struct {
	wsgin.MustAuthRequest

	Response gin.ResponseWriter
}

// New .
func (r *QRCodeRequest) New() wsgin.Process {
	return &QRCodeRequest{}
}

// Extract .
func (r *QRCodeRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	r.Response = c.Writer
	return r.DefaultExtract(r, c)
}

// Exec 获取客户所属二维码
// @Summary 获取客户所属二维码
// @Description get customer qrcode
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Success 200 {string} string "{"status":true}"
// @Router /customers/qrcode [get]
func (r *QRCodeRequest) Exec(ctx context.Context) interface{} {
	data, err := svc.GetQRCode(ctx, r.TokenParames.UID)
	if err == nil {
		r.Response.Write(data)
	}
	return nil
}
