package server

import (
	"context"

	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// CompositeIndexAddRequest .
type CompositeIndexAddRequest struct {
	wsgin.MustAuthRequest

	CompositeDate string  `json:"composite_date" binding:"required"` // 上证指数日期
	Points        float64 `json:"points" binding:"required"`         // 指数
}

// CompositeIndexAddResponse .
type CompositeIndexAddResponse struct {
	wsgin.BaseResponse
}

// New .
func (r *CompositeIndexAddRequest) New() wsgin.Process {
	return &CompositeIndexAddRequest{}
}

// Extract .
func (r *CompositeIndexAddRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec .
// @Summary 添加或者更新上证指数
// @Description post composite index
// @Tags 上证指数
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param args body server.CompositeIndexAddRequest true "参数"
// @Success 200 {object} server.CompositeIndexAddResponse "{"status":true}"
// @Router /composite_index [post]
func (r *CompositeIndexAddRequest) Exec(ctx context.Context) interface{} {
	resp := CompositeIndexAddResponse{}
	code, err := svc.AddCompositeIndex(ctx, r.CompositeDate, r.Points)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	return resp
}
