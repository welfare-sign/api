package wsgin

import (
	"github.com/gin-gonic/gin"

	"welfare-sign/internal/pkg/jwt"
)

// AuthPagingRequest 可选登录后的分页请求基类
type AuthPagingRequest struct {
	BasePagingRequest

	TokenParames *jwt.TokenParames
}

// Extract .
func (r *AuthPagingRequest) Extract(c *gin.Context) (code APICode, err error) {
	return r.DefaultExtract(r, c)
}

// DefaultExtract default extract
func (r *AuthPagingRequest) DefaultExtract(data interface{}, c *gin.Context) (code APICode, err error) {
	return r.ExtractWithBindFunc(data, c, c.ShouldBind)
}

// ExtractWithBindFunc default ExtractWithBindFunc
func (r *AuthPagingRequest) ExtractWithBindFunc(data interface{}, c *gin.Context, bindFunc BindFunc) (code APICode, err error) {
	code, err = r.BaseRequest.ExtractWithBindFunc(data, c, bindFunc)
	if err != nil {
		return
	}
	params, code, err := authFunc(c)
	if err != nil {
		return
	}
	r.TokenParames = params
	return
}
