package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"welfare-sign/internal/pkg/config"
)

var sign = []byte(viper.GetString(config.KeyJWTSign))

// TokenParames 生成的token中含有的参数
type TokenParames struct {
	jwt.StandardClaims

	Name   string // 姓名
	Mobile string // 手机号
}

// CreateToken 生成Token
func CreateToken(name, mobile string) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenParames{
		Name:   name,
		Mobile: mobile,
		StandardClaims: jwt.StandardClaims{
			Issuer:    viper.GetString(config.KeyJWTIssuer),
			ExpiresAt: time.Now().Add(viper.GetDuration(config.KeyJWTExpire) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})
	return tokenClaims.SignedString(sign)
}

// ParseToken 解析Token
func ParseToken(token string) (*TokenParames, error) {
	// 基于公钥验证Token合法性
	tokenClaims, err := jwt.ParseWithClaims(token, &TokenParames{}, func(token *jwt.Token) (interface{}, error) {
		// 基于JWT的第一部分中的alg字段值进行一次验证
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("token's encryption type error")
		}
		return sign, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := tokenClaims.Claims.(*TokenParames); ok && tokenClaims.Valid {
		return c, nil
	}
	return nil, errors.New("token invalid")
}
