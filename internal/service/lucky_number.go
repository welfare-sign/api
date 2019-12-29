package service

import (
	"context"

	"github.com/pkg/errors"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/util"
	"welfare-sign/internal/pkg/wsgin"
)

// CanPartLuckyNumberActivity 用户是否可以参与猜数字活动
func (s *Service) CanPartLuckyNumberActivity(ctx context.Context, customerID, activityID uint64) (bool, wsgin.APICode, error) {
	activity, err := s.dao.FindActivity(ctx, map[string]interface{}{
		"id":     activityID,
		"status": global.ActiveStatus,
	})
	if err != nil {
		return false, apicode.ErrDetail, err
	}
	// 查看当前时间是否在活动截止范围内
	if !util.IsCurrentTimeBetweenD1AndD2(activity.StartTime, activity.EndTime) {
		return false, apicode.ErrActivityClosed, errors.New("活动截止时间已过，请等待下次活动开始")
	}
	if activity.Scope == global.ScopeOpen {
		return true, wsgin.APICodeSuccess, nil
	}
	// 用户是否在指定时间段内领取过福利
	logs, _ := s.dao.IsReceiveBenefitsInD1AndD2(ctx, customerID, activity.StartTime, activity.EndTime)
	if len(logs) == 0 {
		return false, apicode.ErrNoParticipation, errors.New("您还未在时间范围内完成过签到，快去签到领福利吧")
	}
	return true, wsgin.APICodeSuccess, nil
}

// GetLuckyNumberDetail 获取用户猜的幸运数字
func (s *Service) GetLuckyNumberDetail(ctx context.Context, customerID, activityID uint64) (*model.LuckyNumberRecord, wsgin.APICode, error) {
	luckyNumberRecord, err := s.dao.GetLuckyNumberRecord(ctx, customerID, activityID)
	if err != nil {
		return nil, apicode.ErrDetail, err
	}
	return luckyNumberRecord, wsgin.APICodeSuccess, nil
}

// AddLuckyNumber 添加用户猜的幸运数字
func (s *Service) AddLuckyNumber(ctx context.Context, customerID, activityID, num uint64) ([]uint64, wsgin.APICode, error) {
	// 活动时间是否已过
	activity, err := s.dao.FindActivity(ctx, map[string]interface{}{
		"id":     activityID,
		"status": global.ActiveStatus,
	})
	if err != nil {
		return []uint64{}, apicode.ErrSave, err
	}
	if !util.IsCurrentTimeBetweenD1AndD2(activity.StartTime, activity.EndTime) {
		return []uint64{}, apicode.ErrSave, errors.New("活动截止时间已过，请等待下次活动开始")
	}

	nums, err := s.dao.StoreLuckyNumberRecord(ctx, customerID, activityID, num)
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
func (s *Service) GetLuckyPeopleBefore(ctx context.Context) ([]*model.LuckyNumberRecord, wsgin.APICode, error) {
	luckys, err := s.dao.GetLuckyPeopleBefore(ctx)
	if err != nil {
		return nil, apicode.ErrLuckyPeople, err
	}
	return luckys, wsgin.APICodeSuccess, nil
}
