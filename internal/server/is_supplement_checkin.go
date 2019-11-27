package server

import (
	"context"

	"welfare-sign/internal/pkg/wsgin"

	"github.com/gin-gonic/gin"
)

// IsSupplementCheckinRequest 是否是补签
type IsSupplementCheckinRequest struct {
	wsgin.MustAuthRequest
}

// IsSupplementCheckinResponse .
type IsSupplementCheckinResponse struct {
	wsgin.BaseResponse

	Data bool `json:"data"`
}

// New .
func (r *IsSupplementCheckinRequest) New() wsgin.Process {
	return &IsSupplementCheckinRequest{}
}

// Extract .
func (r *IsSupplementCheckinRequest) Extract(c *gin.Context) (code wsgin.APICode, err error) {
	return r.DefaultExtract(r, c)
}

// Exec 是否是补签
// @Summary 是否是补签
// @Description is supplement checkin
// @Security ApiKeyAuth
// @Tags 客户
// @Accept json
// @Produce json
// @Success 200 {object} server.IsSupplementCheckinResponse "{"status":true}"
// @Router /customers/issue_records/is_supplement [get]
func (r *IsSupplementCheckinRequest) Exec(ctx context.Context) interface{} {
	resp := IsSupplementCheckinResponse{}

	data, code, err := svc.IsSupplement(ctx, r.TokenParames.UID)
	resp.BaseResponse = wsgin.NewResponse(ctx, code, err)
	resp.Data = data
	return resp
}
