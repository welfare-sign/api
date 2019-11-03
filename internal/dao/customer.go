package dao

import (
	"context"

	"welfare-sign/internal/dao/mysql"
	"welfare-sign/internal/model"
)

// ListCustomer get customer list
func (d *dao) ListCustomer(ctx context.Context, query interface{}, pageNo, pageSize int) ([]*model.Customer, int, error) {
	var customers []*model.Customer
	total := 0
	if err := d.db.Where(query).Limit(pageSize).Offset((pageNo - 1) * pageSize).Order("created_at desc").Find(&customers).Error; mysql.IsError(err) {
		return customers, total, err
	}
	if err := d.db.Where(query).Find(&model.Customer{}).Count(&total).Error; mysql.IsError(err) {
		return customers, total, err
	}
	return customers, total, nil
}

// FindCustomer 获取客户详情
func (d *dao) FindCustomer(ctx context.Context, query interface{}) (*model.Customer, error) {
	var customer model.Customer
	err := d.db.Where(query).First(&customer).Error
	return &customer, err
}
