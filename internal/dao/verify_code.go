package dao

import (
	"context"
	"time"

	"github.com/spf13/viper"

	"welfare-sign/internal/pkg/config"
)

// KeyCacheSMSCodePrefix 验证码在缓存中的键前缀
const KeyCacheSMSCodePrefix = "welfare:sms:"

// SaveSMSCode 保存验证码，后续验证使用
func (d *dao) SaveSMSCode(ctx context.Context, mobile, code string) error {
	return d.cache.Set(KeyCacheSMSCodePrefix+mobile, code, viper.GetDuration(config.KeyYuanpianExpire)*time.Minute).Err()
}

// GetSMSCode 验证传入的手机，获取验证码
func (d *dao) GetSMSCode(ctx context.Context, mobile string) (string, error) {
	return d.cache.Get(KeyCacheSMSCodePrefix + mobile).Result()
}
