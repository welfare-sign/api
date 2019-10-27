package dao

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	"welfare-sign/internal/dao/cache"
	"welfare-sign/internal/dao/mysql"
	"welfare-sign/internal/model"
)

// Dao dao interface
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
	CreateMerchant(ctx context.Context, data model.Merchant) error
	ListMerchant(ctx context.Context, query interface{}, pageNo, pageSize int) ([]*model.Merchant, error)
	FindUser(ctx context.Context, data model.User) (*model.User, error)
	ListCustomer(ctx context.Context, query interface{}, pageNo, pageSize int) ([]*model.Customer, error)
	SaveSMSCode(ctx context.Context, mobile, code string) error
}

// dao dao.
type dao struct {
	db    *gorm.DB
	cache *redis.Client
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// New new a dao and return.
func New() Dao {
	return &dao{
		db:    mysql.New(),
		cache: cache.New(),
	}
}

// Close close the resource.
func (d *dao) Close() {
	d.db.Close()
	d.cache.Close()
}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	if err := d.db.DB().PingContext(ctx); err != nil {
		return err
	}
	if err := d.cache.Ping().Err(); err != nil {
		return err
	}
	return nil
}
