package service

import (
	"context"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/util"
)

// AddMerchant 新增商户
func (s *Service) AddMerchant(ctx context.Context, vo model.MerchantVO) error {
	var data model.Merchant
	if err := util.StructCopy(&data, &vo); err != nil {
		return err
	}
	return s.dao.CreateMerchant(ctx, data)
}

// GetMerchantList 获取商户列表
func (s *Service) GetMerchantList(ctx context.Context, vo model.MerchantListVO) ([]*model.Merchant, error) {
	query := make(map[string]interface{})
	if vo.StoreName != "" {
		query["store_name"] = vo.StoreName
	}
	if vo.ContactName != "" {
		query["contact_name"] = vo.ContactName
	}
	if vo.ContactPhone != "" {
		query["contact_phone"] = vo.ContactPhone
	}
	if vo.PageNo == 0 {
		vo.PageNo = 1
	}
	if vo.PageSize == 0 {
		vo.PageSize = 10
	}

	return s.dao.ListMerchant(ctx, query, vo.PageNo, vo.PageSize)
}
