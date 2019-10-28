package service

import (
	"context"

	"welfare-sign/internal/model"
)

// GetCustomerList 获取客户列表
func (s *Service) GetCustomerList(ctx context.Context, vo model.CustomerListVO) ([]*model.Customer, error) {
	query := make(map[string]interface{})
	if vo.Name != "" {
		query["name"] = vo.Name
	}
	if vo.Mobile != "" {
		query["mobile"] = vo.Mobile
	}
	if vo.PageNo == 0 {
		vo.PageNo = 1
	}
	if vo.PageSize == 0 {
		vo.PageSize = 10
	}

	return s.dao.ListCustomer(ctx, query, vo.PageNo, vo.PageSize)
}

// GetCustomerDetail 获取客户详情
func (s *Service) GetCustomerDetail(ctx context.Context, customerId uint64) (*model.Customer, error) {
	return s.dao.FindCustomer(ctx, map[string]interface{}{"id": customerId})
}
