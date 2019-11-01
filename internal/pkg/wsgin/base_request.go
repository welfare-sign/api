package wsgin

import "github.com/gin-gonic/gin"

// BaseRequest 基础请求
type BaseRequest struct{}

// Extract extract
func (r *BaseRequest) Extract(c *gin.Context) (code APICode, err error) {
	return r.DefaultExtract(r, c)
}

// DefaultExtract default extract
func (r *BaseRequest) DefaultExtract(data interface{}, c *gin.Context) (code APICode, err error) {
	return r.ExtractWithBindFunc(data, c, c.ShouldBind)
}

// ExtractWithBindFunc extract with bindFunc
func (r *BaseRequest) ExtractWithBindFunc(data interface{}, c *gin.Context, bindFunc BindFunc) (code APICode, err error) {
	err = bindFunc(data)
	return
}
