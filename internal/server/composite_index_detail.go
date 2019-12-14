package server

import (
	"context"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// CompositeIndexDetailRequest .
type CompositeIndexDetailRequest struct {
	wsgin.MustAuthRequest

	CompositeDate string `form:"composite_date" json:"composite_date" binding:"required"` // 日期
}

// CompositeIndexDetailResponse .
type CompositeIndexDetailResponse struct {
	wsgin.BaseResponse

	Data *model.CompositeIndex `json:"data"`
}

// New .
func (r *CompositeIndexDetailRequest) New() wsgin.Process {
	return &CompositeIndexDetailRequest{}
}

// Extract .
func (r *CompositeIndexDetailRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 获取上证指数
// @Summary 获取上证指数
// @Description get composite index
// @Tags 上证指数
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param composite_date query string true "日期"
// @Success 200 {object} server.CompositeIndexDetailResponse "{"status":true}"
// @Router /composite_index [get]
func (r *CompositeIndexDetailRequest) Exec(ctx context.Context) interface{} {
	resp := CompositeIndexDetailResponse{}
	data, code, err := svc.GetCompositeIndex(ctx, r.CompositeDate)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
