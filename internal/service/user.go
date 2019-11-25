package service

import (
	"context"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/jwt"
	"welfare-sign/internal/pkg/wsgin"
)

// UserLogin 后台用户登录
func (s *Service) UserLogin(ctx context.Context, vo *model.UserVO) (string, wsgin.APICode, error) {
	user, err := s.dao.FindUser(ctx, map[string]interface{}{
		"name":     vo.Name,
		"password": vo.Password,
	})
	if err != nil {
		return "", apicode.ErrLogin, err
	}
	token, err := jwt.CreateToken(user.ID, user.Name, "")
	if err != nil {
		return "", apicode.ErrCreateToken, err
	}
	return token, wsgin.APICodeSuccess, nil
}

// GetAllCustomerCheckinRecordList 获取用户有效的签到记录列表
// TODO: 临时
func (s *Service) GetAllCustomerCheckinRecordList(ctx context.Context) ([]*model.CheckinRecordListResp, wsgin.APICode, error) {
	checkinRecordListResps, _ := s.dao.GetTmpCheckinRecordList(ctx)
	return checkinRecordListResps, wsgin.APICodeSuccess, nil
}

// ModifyCustomerCheckinRecord 后台用户更新用户签到状态
// TODO: 临时
func (s *Service) ModifyCustomerCheckinRecord(ctx context.Context, checkinRecordID uint64, status string) (wsgin.APICode, error) {
	if err := s.dao.UpdateCustomerCheckinRecord(ctx, checkinRecordID, status); err != nil {
		return wsgin.APICodeDefault, err
	}
	return wsgin.APICodeSuccess, nil
}
