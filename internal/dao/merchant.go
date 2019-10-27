package dao

import (
	"context"

	"welfare-sign/internal/model"
)

// CreateMerchant create merchant
func (d *dao) CreateMerchant(ctx context.Context, data model.Merchant) error {
	data.SetDefaultAttr()
	return d.db.Create(&data).Error
}

// ListMerchant get merchant list
// pageNo >= 1
func (d *dao) ListMerchant(ctx context.Context, query interface{}, pageNo, pageSize int) ([]*model.Merchant, error) {
	var merchants []*model.Merchant
	err := d.db.Where(query).Limit(pageSize).Offset((pageNo - 1) * pageSize).Order("created_at desc").Find(&merchants).Error
	return merchants, err
}
