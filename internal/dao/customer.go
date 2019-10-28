package dao

import (
	"context"

	"welfare-sign/internal/model"
)

// ListCustomer get customer list
func (d *dao) ListCustomer(ctx context.Context, query interface{}, pageNo, pageSize int) ([]*model.Customer, error) {
	var customers []*model.Customer
	err := d.db.Where(query).Limit(pageSize).Offset((pageNo - 1) * pageSize).Order("created_at desc").Find(&customers).Error
	return customers, err
}

// FindCustomer 获取客户详情
func (d *dao) FindCustomer(ctx context.Context, query interface{}) (*model.Customer, error) {
	var customer model.Customer
	err := d.db.Where(query).First(&customer).Error
	return &customer, err
}
