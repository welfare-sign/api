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
	ListMerchant(ctx context.Context, query interface{}, pageNo, pageSize int) ([]*model.Merchant, int, error)
	FindUser(ctx context.Context, query interface{}) (*model.User, error)
	ListCustomer(ctx context.Context, query interface{}, pageNo, pageSize int) ([]*model.Customer, int, error)
	SaveSMSCode(ctx context.Context, mobile, code string) error
	GetSMSCode(ctx context.Context, mobile string) (string, error)
	FindMerchant(ctx context.Context, query interface{}) (*model.Merchant, error)
	FindCustomer(ctx context.Context, query interface{}) (*model.Customer, error)
	FindIssueRecord(ctx context.Context, query interface{}) (*model.IssueRecord, error)
	EcecWriteOff(ctx context.Context, merchantId, customerId, hasRece, totalRece uint64) error
	ListCheckinRecord(ctx context.Context, query interface{}) ([]*model.CheckinRecord, error)
	InitCheckinRecords(ctx context.Context, customerId uint64) ([]*model.CheckinRecord, error)
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
