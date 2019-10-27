package dao

import (
	"context"
	"time"

	"github.com/spf13/viper"

	"welfare-sign/internal/pkg/config"
)

// KeyCacheSMSCodePrefix 验证码在缓存中的键前缀
const KeyCacheSMSCodePrefix = "welfare:sms:code:"

// SaveSMSCode 保存验证码，后续验证使用
func (d *dao) SaveSMSCode(ctx context.Context, mobile, code string) error {
	return d.cache.Set(KeyCacheSMSCodePrefix+mobile, code, viper.GetDuration(config.KeyYuanpianExpire)*time.Minute).Err()
}
