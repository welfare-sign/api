package service

import (
	"context"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/jwt"
	"welfare-sign/internal/pkg/util"
)

// UserLogin 后台用户登录
func (s *Service) UserLogin(ctx context.Context, vo model.UserVO) (string, error) {
	var data model.User
	if err := util.StructCopy(&data, &vo); err != nil {
		return "", err
	}
	if _, err := s.dao.FindUser(ctx, data); err != nil {
		return "", err
	}
	return jwt.CreateToken(vo.Name, "")
}
