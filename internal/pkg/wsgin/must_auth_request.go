package wsgin

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"welfare-sign/internal/pkg/jwt"
)

// MustAuthRequest 用户必须登录才可访问接口
type MustAuthRequest struct {
	BaseRequest

	TokenParames *jwt.TokenParames `json:"-"`
}

// Extract .
func (r *MustAuthRequest) Extract(c *gin.Context) (code APICode, err error) {
	return r.DefaultExtract(r, c)
}

// DefaultExtract default extract
func (r *MustAuthRequest) DefaultExtract(data interface{}, c *gin.Context) (code APICode, err error) {
	return r.ExtractWithBindFunc(data, c, c.ShouldBind)
}

// ExtractWithBindFunc default ExtractWithBindFunc
func (r *MustAuthRequest) ExtractWithBindFunc(data interface{}, c *gin.Context, bindFunc BindFunc) (code APICode, err error) {
	code, err = r.BaseRequest.ExtractWithBindFunc(data, c, bindFunc)
	if err != nil {
		return
	}
	params, code, err := mustAuthFunc(c)
	if err != nil {
		return
	}
	r.TokenParames = params
	return
}

func mustAuthFunc(c *gin.Context) (*jwt.TokenParames, APICode, error) {
	token := strings.TrimSpace(c.Query("access_token"))
	if token == "" {
		token = strings.TrimSpace(c.GetHeader("Authorization"))
	}
	if token == "" {
		return nil, APICodeNoPermission, errors.New("token not exists")
	}
	tokenParams, err := jwt.ParseToken(token)
	if err != nil {
		return nil, APICodeNoPermission, err
	}
	return tokenParams, APICodeSuccess, nil
}
