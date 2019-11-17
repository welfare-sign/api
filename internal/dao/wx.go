package dao

import (
	"context"
	"time"

	"welfare-sign/internal/model"
)

// 缓存中key
const (
	KeyWXAccessToken = "ws:wx:ak"
	KeyWXJSTicket    = "ws:wx:ticket"
)

// StoreWXAccessToken 存储微信access_token
func (d *dao) StoreWXAccessToken(ak string, expire time.Duration) error {
	return checkCacheError(d.cache.Set(KeyWXAccessToken, ak, expire).Err())
}

// StoreWXJSTicket 存储微信JS API ticket
func (d *dao) StoreWXJSTicket(ticket string, expire time.Duration) error {
	return checkCacheError(d.cache.Set(KeyWXJSTicket, ticket, expire).Err())
}

// GetWXAccessToken 获取微信access_token
func (d *dao) GetWXAccessToken() (string, error) {
	res, err := d.cache.Get(KeyWXAccessToken).Result()
	return res, checkCacheError(err)
}

// GetWXJSTicket 获取微信JS API ticket
func (d *dao) GetWXJSTicket() (string, error) {
	res, err := d.cache.Get(KeyWXJSTicket).Result()
	return res, checkCacheError(err)
}

// FindWXPayRecord 查询微信支付流水号
func (d *dao) FindWXPayRecord(ctx context.Context, query map[string]interface{}) (*model.WXPayRecord, error) {
	var record model.WXPayRecord
	err := checkErr(d.db.Where(query).First(&record).Error)
	return &record, err
}
