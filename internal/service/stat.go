package service

import (
	"context"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// GetRegisterStat 统计注册用户
func (s *Service) GetRegisterStat(ctx context.Context, beginDate, endDate string) ([]*model.RegisterStat, wsgin.APICode, error) {
	stats, err := s.dao.GetRegisterStat(ctx, beginDate, endDate)
	if err != nil {
		return stats, apicode.ErrGetListData, err
	}
	return stats, wsgin.APICodeSuccess, nil
}

// GetCheckinStat 统计用户执行签到次数
func (s *Service) GetCheckinStat(ctx context.Context, beginDate, endDate string) ([]*model.CheckinStat, wsgin.APICode, error) {
	stats, err := s.dao.GetCheckinStat(ctx, beginDate, endDate)
	if err != nil {
		return stats, apicode.ErrGetListData, err
	}
	return stats, wsgin.APICodeSuccess, nil
}
