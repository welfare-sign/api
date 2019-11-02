package service

import (
	"context"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// GetCustomerList 获取客户列表
func (s *Service) GetCustomerList(ctx context.Context, vo *model.CustomerListVO) ([]*model.Customer, int, wsgin.APICode, error) {
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

	customers, total, err := s.dao.ListCustomer(ctx, query, vo.PageNo, vo.PageSize)
	if err != nil {
		return nil, total, apicode.ErrGetListData, err
	}
	return customers, total, wsgin.APICodeSuccess, nil
}

// GetCustomerDetail 获取客户详情
func (s *Service) GetCustomerDetail(ctx context.Context, customerId uint64) (*model.Customer, wsgin.APICode, error) {
	customer, err := s.dao.FindCustomer(ctx, map[string]interface{}{"id": customerId})
	if err != nil {
		return nil, apicode.ErrDetail, err
	}
	return customer, wsgin.APICodeSuccess, nil
}

// GetCustomerCheckinRecord 获取用户签到记录
func (s *Service) GetCustomerCheckinRecord(ctx context.Context, customerId uint64) ([]*model.CheckinRecord, wsgin.APICode, error) {
	records, err := s.dao.ListCheckinRecord(ctx, map[string]interface{}{
		"status":      "A",
		"customer_id": customerId,
	})
	if err != nil {
		return nil, apicode.ErrGetCheckinRecord, err
	}
	// 用户无签到记录时，自动创建5条信息并返回
	if len(records) == 0 {
		records, err = s.dao.InitCheckinRecords(ctx, customerId)
		if err != nil {
			return nil, apicode.ErrGetCheckinRecord, err
		}
	}
	return records, wsgin.APICodeSuccess, nil
}
