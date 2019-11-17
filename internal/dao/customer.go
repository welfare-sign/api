package dao

import (
	"context"
	"time"

	"welfare-sign/internal/dao/mysql"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/util"
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
	err := checkErr(d.db.Where(query).First(&customer).Error)
	return &customer, err
}

// UpsertCustomer update or insert customer
func (d *dao) UpsertCustomer(ctx context.Context, data *model.WxUserResp) (customer *model.Customer, err error) {
	customer, err = d.FindCustomer(ctx, map[string]interface{}{"open_id": data.OpenID})
	if checkErr(err) != nil {
		return
	}
	if customer.ID == 0 {
		var c model.Customer
		util.StructCopy(&c, data)
		c.SetDefaultAttr()
		if err := d.db.Create(&c).Error; err != nil {
			return nil, err
		}
		return d.FindCustomer(ctx, map[string]interface{}{"open_id": data.OpenID})
	}
	customer.UpdatedAt = time.Now()
	util.StructCopy(&customer, data)
	return customer, d.db.Save(customer).Error
}
