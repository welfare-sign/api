package server

import (
	"context"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// LuckyPeopleBeforeRequest .
type LuckyPeopleBeforeRequest struct {
	wsgin.MustAuthRequest
}

// LuckyPeopleBeforeResponse .
type LuckyPeopleBeforeResponse struct {
	wsgin.BaseResponse

	Data *model.Customer `json:"data"`
}

// New .
func (r *LuckyPeopleBeforeRequest) New() wsgin.Process {
	return &LuckyPeopleBeforeRequest{}
}

// Extract .
func (r *LuckyPeopleBeforeRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 上期幸运儿
// @Summary 上期幸运儿
// @Description get lucky customer before
// @Tags 客户
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} server.LuckyPeopleBeforeResponse "{"status":true}"
// @Router /customers/lucky/before [get]
func (r *LuckyPeopleBeforeRequest) Exec(ctx context.Context) interface{} {
	resp := LuckyPeopleBeforeResponse{}
	data, code, err := svc.GetLuckyPeopleBefore(ctx)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
