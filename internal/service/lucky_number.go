package service

import (
	"context"

	"github.com/pkg/errors"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// CanPartLuckyNumberActivity 用户是否可以参与猜数字活动
func (s *Service) CanPartLuckyNumberActivity(ctx context.Context, customerID uint64) (bool, wsgin.APICode, error) {
	// 用户是否在指定时间段内领取过福利
	logs, _ := s.dao.IsReceiveBenefits(ctx, customerID)
	if len(logs) == 0 {
		return false, apicode.ErrNoParticipation, errors.New("您还未在时间范围内完成过签到，快去签到领福利吧")
	}
	return true, wsgin.APICodeSuccess, nil
}

// GetLuckyNumberDetail 获取用户猜的幸运数字
func (s *Service) GetLuckyNumberDetail(ctx context.Context, customerID uint64) (*model.LuckyNumberRecord, wsgin.APICode, error) {
	luckyNumberRecord, err := s.dao.GetLuckyNumberRecord(ctx, customerID)
	if err != nil {
		return nil, apicode.ErrDetail, err
	}
	return luckyNumberRecord, wsgin.APICodeSuccess, nil
}

// AddLuckyNumber 添加用户猜的幸运数字
func (s *Service) AddLuckyNumber(ctx context.Context, customerID uint64, num int64) ([]int64, wsgin.APICode, error) {
	nums, err := s.dao.StoreLuckyNumberRecord(ctx, customerID, num)
	if err != nil {
		return nums, apicode.ErrSave, err
	}
	return nil, wsgin.APICodeSuccess, nil
}

// GetLuckyNumberBefore 获取用户上期猜的数字
func (s *Service) GetLuckyNumberBefore(ctx context.Context, customerID uint64) (*model.LuckyNumberRecord, wsgin.APICode, error) {
	luckyNumberRecord, err := s.dao.GetLuckyNumberRecordBefore(ctx, customerID)
	if err != nil {
		return nil, apicode.ErrDetail, err
	}
	return luckyNumberRecord, wsgin.APICodeSuccess, nil
}

// GetLuckyPeopleBefore 获取上期幸运用户
func (s *Service) GetLuckyPeopleBefore(ctx context.Context) (*model.Customer, wsgin.APICode, error) {
	customer, err := s.dao.GetLuckyPeopleBefore(ctx)
	if err != nil {
		return nil, apicode.ErrLuckyPeople, err
	}
	return customer, wsgin.APICodeSuccess, nil
}
