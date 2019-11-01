package wsgin

import "github.com/gin-gonic/gin"

// BasePagingRequest 带有分页请求参数
type BasePagingRequest struct {
	BaseRequest

	PageNo   int `form:"page_no" json:"page_no" example:"1" minimum:"1"`
	PageSize int `form:"page_size" json:"page_size" example:"10" minimum:"1" maximum:"20" binding:"gte=1,lte=20"`
}

// Extract extract
func (r *BasePagingRequest) Extract(c *gin.Context) (code APICode, err error) {
	return r.DefaultExtract(r, c)
}

// DefaultExtract default extract
func (r *BasePagingRequest) DefaultExtract(data interface{}, c *gin.Context) (code APICode, err error) {
	return r.ExtractWithBindFunc(data, c, c.ShouldBind)
}
