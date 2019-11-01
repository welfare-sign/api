package service

import (
	"context"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/jwt"
	"welfare-sign/internal/pkg/wsgin"
)

// UserLogin 后台用户登录
func (s *Service) UserLogin(ctx context.Context, vo *model.UserVO) (string, wsgin.APICode, error) {
	user, err := s.dao.FindUser(ctx, map[string]interface{}{
		"name":     vo.Name,
		"password": vo.Password,
	})
	if err != nil {
		return "", apicode.ErrLogin, err
	}
	token, err := jwt.CreateToken(user.ID, user.Name, "")
	if err != nil {
		return "", apicode.ErrCreateToken, err
	}
	return token, wsgin.APICodeSuccess, nil
}
