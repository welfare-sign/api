package service

import (
	"context"
	"errors"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/wsgin"
)

// UpsertActivity 添加或更新活动信息
func (s *Service) UpsertActivity(ctx context.Context, vo *model.ActivityVO) (wsgin.APICode, error) {
	isLegal, err := s.dao.IsActivityDateLegal(ctx, vo)
	if err != nil {
		return apicode.ErrSave, err
	}
	if !isLegal {
		return apicode.ErrActivityDate, errors.New("所选日期不能和已存在的活动日期重叠")
	}
	if err := s.dao.UpsertActivity(ctx, vo); err != nil {
		return apicode.ErrSave, err
	}
	return wsgin.APICodeSuccess, nil
}

// DetailActivity 获取活动详情
func (s *Service) DetailActivity(ctx context.Context, activityID uint64, name string) (*model.Activity, wsgin.APICode, error) {
	var (
		activity *model.Activity
		err      error
	)
	if activityID != 0 {
		activity, err = s.dao.FindActivity(ctx, map[string]interface{}{
			"id":     activityID,
			"status": global.ActiveStatus,
		})
	} else {
		activity, err = s.dao.FindActivity(ctx, map[string]interface{}{
			"name":   name,
			"status": global.ActiveStatus,
		})
	}

	if err != nil {
		return nil, apicode.ErrDetail, err
	}
	if activity.ID == 0 {
		return nil, wsgin.APICodeSuccess, nil
	}
	return activity, wsgin.APICodeSuccess, nil
}

// ListActivity 获取活动列表
func (s *Service) ListActivity(ctx context.Context, vo *model.ActivityListVO) ([]*model.Activity, int, wsgin.APICode, error) {
	query := make(map[string]interface{})
	if vo.Name != "" {
		query["name"] = vo.Name
	}
	if vo.PageNo == 0 {
		vo.PageNo = 1
	}
	if vo.PageSize == 0 {
		vo.PageSize = 10
	}

	activitys, total, err := s.dao.ListActivity(ctx, query, vo.PageNo, vo.PageSize)
	if err != nil {
		return nil, total, apicode.ErrGetListData, err
	}
	return activitys, total, wsgin.APICodeSuccess, nil
}

// DrawActivity 活动开奖
func (s *Service) DrawActivity(ctx context.Context, activityID, number uint64) (*model.Activity, wsgin.APICode, error) {
	activity, err := s.dao.DrawActivity(ctx, activityID, number)
	if err != nil {
		return nil, apicode.ErrActivityDraw, err
	}
	return activity, wsgin.APICodeSuccess, nil
}

// ListActivityParticipant 活动参与者信息列表
func (s *Service) ListActivityParticipant(ctx context.Context, vo *model.ActivityParticipantListVO) ([]*model.LuckyNumberRecord, int, wsgin.APICode, error) {
	query := "activity_id = ?"
	queryValue := make([]interface{}, 0)
	queryValue = append(queryValue, vo.ActivityID)
	if vo.Mobile != "" {
		query += " and mobile = ?"
		queryValue = append(queryValue, vo.Mobile)
	}
	if vo.IsSearchWin {
		activity, err := s.dao.FindActivity(ctx, map[string]interface{}{
			"id":     vo.ActivityID,
			"status": global.ActiveStatus,
		})
		if err != nil {
			return nil, 0, apicode.ErrGetListData, err
		}
		query += " and ranking BETWEEN ? AND ?"
		queryValue = append(queryValue, 1)
		queryValue = append(queryValue, activity.PrizeAmount)
	}

	if vo.PageNo == 0 {
		vo.PageNo = 1
	}
	if vo.PageSize == 0 {
		vo.PageSize = 10
	}

	luckys, total, err := s.dao.ListActivityParticipant(ctx, vo.PageNo, vo.PageSize, query, queryValue)
	if err != nil {
		return nil, total, apicode.ErrGetListData, err
	}
	return luckys, total, wsgin.APICodeSuccess, nil
}

// DelActivity 删除活动
func (s *Service) DelActivity(ctx context.Context, activityID uint64) (wsgin.APICode, error) {
	activity, err := s.dao.FindActivity(ctx, map[string]interface{}{"id": activityID})
	if err != nil {
		return apicode.ErrDelete, err
	}
	if activity.PrizeNumber != 0 {
		return apicode.ErrDelete, errors.New("活动已开奖，无法删除")
	}
	if err := s.dao.DelActivity(ctx, activityID); err != nil {
		return apicode.ErrDelete, err
	}
	return wsgin.APICodeSuccess, nil
}

// CurrentlyAvailableActivity 当前可参与的活动
func (s *Service) CurrentlyAvailableActivity(ctx context.Context) (*model.Activity, wsgin.APICode, error) {
	activity, err := s.dao.CurrentlyAvailableActivity(ctx)
	if err != nil {
		return nil, apicode.ErrDetail, err
	}
	return activity, wsgin.APICodeSuccess, nil
}

// ActivityAllPrizeIssued 活动所有已发放的奖品
func (s *Service) ActivityAllPrizeIssued(ctx context.Context) (int, wsgin.APICode, error) {
	total, err := s.dao.ActivityAllPrizeIssued(ctx)
	if err != nil {
		return 0, wsgin.APICodeDefault, err
	}
	num := total*97 + total*5 + 5029
	return num, wsgin.APICodeSuccess, nil
}
