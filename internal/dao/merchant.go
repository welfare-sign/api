package dao

import (
	"context"

	"welfare-sign/internal/dao/mysql"
	"welfare-sign/internal/model"
)

// CreateMerchant create merchant
func (d *dao) CreateMerchant(ctx context.Context, data model.Merchant) error {
	data.SetDefaultAttr()
	return d.db.Create(&data).Error
}

// ListMerchant get merchant list
// pageNo >= 1
func (d *dao) ListMerchant(ctx context.Context, query interface{}, pageNo, pageSize int) ([]*model.Merchant, int, error) {
	var merchants []*model.Merchant
	total := 0
	err := d.db.Where(query).Limit(pageSize).Offset((pageNo - 1) * pageSize).Order("created_at desc").Find(&merchants).Error
	if mysql.IsError(err) {
		return merchants, total, err
	}
	if err := d.db.Where(query).Find(&model.Merchant{}).Count(&total).Error; mysql.IsError(err) {
		return merchants, total, err
	}
	return merchants, total, nil
}

// FindMerchant 获取商家详情
func (d *dao) FindMerchant(ctx context.Context, query interface{}) (*model.Merchant, error) {
	var merchant model.Merchant
	err := d.db.Where(query).First(&merchant).Error
	return &merchant, err
}

// EcecWriteOff 执行核销
func (d *dao) EcecWriteOff(ctx context.Context, merchantId, customerId, hasRece, totalRece uint64) error {
	tx := d.db.Begin()
	if err := tx.Model(&model.IssueRecord{}).Where(map[string]interface{}{
		"merchant_id": merchantId,
		"customer_id": customerId,
	}).Updates(map[string]interface{}{"received": hasRece}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Merchant{}).Where("id = ?", merchantId).Updates(map[string]interface{}{"received": totalRece}).Error; err != nil {
		tx.Rollback()
		return nil
	}
	tx.Commit()
	return nil
}
