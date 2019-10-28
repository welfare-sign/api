package service

import (
	"context"
	"errors"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"welfare-sign/internal/pkg/config"
	"welfare-sign/internal/pkg/log"
	"welfare-sign/internal/pkg/sms"
	"welfare-sign/internal/pkg/util"
)

// SendVerifyCode 发送验证码
func (s *Service) SendVerifyCode(ctx context.Context, mobile string) error {
	code := util.GenerateCode()

	if err := sms.Send(mobile, viper.GetString(config.KeySMSTemplate), map[string]string{
		"code": code,
	}); err != nil {
		return err
	}
	if err := s.dao.SaveSMSCode(ctx, mobile, code); err != nil {
		log.Warn(ctx, "验证码保存到缓存失败", zap.Error(err))
	}
	return nil
}

// ValidateCode 根据传入的手机号，验证码验证是否正确
func (s *Service) ValidateCode(ctx context.Context, mobile, code string) error {
	res, err := s.dao.GetSMSCode(ctx, mobile)
	if err != nil {
		return err
	}
	if res != code {
		return errors.New("验证码不正确")
	}
	return nil
}
