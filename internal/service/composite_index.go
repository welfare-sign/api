package service

import (
	"context"
	"errors"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// AddCompositeIndex 添加上证指数
func (s *Service) AddCompositeIndex(ctx context.Context, compositeDate string, points float64) (wsgin.APICode, error) {
	if points < 1000 {
		return apicode.ErrSave, errors.New("上证指数不正确")
	}
	if err := s.dao.StoreCompositeIndex(ctx, compositeDate, points); err != nil {
		return apicode.ErrSave, err
	}
	return wsgin.APICodeSuccess, nil
}

// GetCompositeIndex 获取上证指数
func (s *Service) GetCompositeIndex(ctx context.Context, compositeDate string) (*model.CompositeIndex, wsgin.APICode, error) {
	data, err := s.dao.GetCompositeIndexByQuery(ctx, map[string]interface{}{
		"status":         global.ActiveStatus,
		"composite_date": compositeDate,
	})
	if err != nil {
		return nil, apicode.ErrDetail, err
	}
	return data, wsgin.APICodeSuccess, nil
}
