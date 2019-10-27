package cache

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"welfare-sign/internal/pkg/config"
)

// New new redis client
func New() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString(config.KeyRedisHost),
		Password: viper.GetString(config.KeyRedisPWD),
		DB:       viper.GetInt(config.KeyRedisDB),
	})
	if err := client.Ping().Err(); err != nil {
		panic(errors.WithMessage(err, "cache.New() error"))
	}
	return client
}
