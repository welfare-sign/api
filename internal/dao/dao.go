package dao

import (
	"context"
	"time"

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
	EcecWriteOff(ctx context.Context, merchantID, customerID, hasRece, writeOffNum uint64) error
	ListCheckinRecord(ctx context.Context, query interface{}, args ...interface{}) ([]*model.CheckinRecord, error)
	InitCheckinRecords(ctx context.Context, customerID uint64) ([]*model.CheckinRecord, error)
	UpsertCustomer(ctx context.Context, data *model.WxUserResp) (*model.Customer, error)
	NearMerchant(ctx context.Context, data *model.NearMerchantVO) ([]*model.Merchant, error)
	FindCheckinRecord(ctx context.Context, query interface{}, args ...interface{}) (*model.CheckinRecord, error)
	ExecCheckin(ctx context.Context, customerID uint64) error
	ListIssueRecord(ctx context.Context, query interface{}, args ...interface{}) ([]*model.IssueRecord, error)
	ListIssueRecordDetail(ctx context.Context, query interface{}, args ...interface{}) ([]*model.IssueRecord, error)
	CreateIssueRecord(ctx context.Context, data model.IssueRecord, merchant *model.Merchant, mobile string) error
	InvalidCheckin(ctx context.Context, customerID uint64) error
	HelpCheckin(ctx context.Context, checkRecordID, customerID, helpCustomerID uint64) error
	StoreWXAccessToken(ak string, expire time.Duration) error
	StoreWXJSTicket(ticket string, expire time.Duration) error
	GetWXAccessToken() (string, error)
	GetWXJSTicket() (string, error)
	HasChecked(ctx context.Context, customerID uint64) (bool, error)
	GetUnchecked(ctx context.Context, customerID uint64) (*model.CheckinRecord, error)
	GetAllUnchecked(ctx context.Context, customerID uint64) ([]*model.CheckinRecord, error)
	PayCheckin(ctx context.Context, checkRecordIds []uint64, customerID uint64, payRecord *model.WXPayRecord) error
	FindWXPayRecord(ctx context.Context, query map[string]interface{}) (*model.WXPayRecord, error)
	UpdateMerchant(ctx context.Context, data *model.Merchant) error
	DeleteMerchant(ctx context.Context, merchantID uint64)
	UpdateCustomer(ctx context.Context, data *model.Customer) error
	DeleteCustomer(ctx context.Context, customerID uint64)
	GetRoundMerchantPoster() (*model.Merchant, error)
	DelSMSCode(ctx context.Context, mobile string) error
	GetTmpCheckinRecordList(ctx context.Context) ([]*model.CheckinRecordListResp, error)
	UpdateCustomerCheckinRecord(ctx context.Context, checkinRecord uint64, status string) error
	FindHelpCheckinMesage(ctx context.Context, query interface{}, args ...interface{}) (*model.HelpCheckinMessage, error)
	UpdateHelpCheckinMessage(ctx context.Context, customerID uint64) error
	GetNeedClearIssueRecords(ctx context.Context) ([]*model.IssueRecord, error)
	FailureIssueRecord(ctx context.Context, issueRecord *model.IssueRecord) error
	IsReceiveBenefitsInD1AndD2(ctx context.Context, customerID uint64, d1, d2 time.Time) ([]*model.IssueRecordLog, error)
	GetLuckyNumberRecord(ctx context.Context, customerID, activityID uint64) (*model.LuckyNumberRecord, error)
	StoreLuckyNumberRecord(ctx context.Context, customerID, activityID, num uint64) ([]uint64, error)
	GetLuckyNumberRecordBefore(ctx context.Context, customerID uint64) (*model.LuckyNumberRecord, error)
	GetLuckyPeopleBefore(ctx context.Context) ([]*model.LuckyNumberRecord, error)
	GetRegisterStat(ctx context.Context, beginDate, endDate string) ([]*model.RegisterStat, error)
	GetCheckinStat(ctx context.Context, beginDate, endDate string) ([]*model.CheckinStat, error)
	UpsertActivity(ctx context.Context, data *model.ActivityVO) error
	FindActivity(ctx context.Context, query interface{}, args ...interface{}) (*model.Activity, error)
	ListActivity(ctx context.Context, query interface{}, pageNo, pageSize int) ([]*model.Activity, int, error)
	ListActivityParticipant(ctx context.Context, pageNo, pageSize int, query interface{}, args ...interface{}) ([]*model.LuckyNumberRecord, int, error)
	DrawActivity(ctx context.Context, activityID, number uint64) (*model.Activity, error)
	IsActivityDateLegal(ctx context.Context, data *model.ActivityVO) (bool, error)
	DelActivity(ctx context.Context, activityID uint64) error
	CurrentlyAvailableActivity(ctx context.Context) (*model.Activity, error)
	ActivityAllPrizeIssued(ctx context.Context) (int, error)
}

// dao dao.
type dao struct {
	db    *gorm.DB
	cache *redis.Client
}

func checkErr(err error) error {
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

func checkCacheError(err error) error {
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
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
