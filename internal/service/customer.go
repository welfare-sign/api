package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/config"
	"welfare-sign/internal/pkg/jwt"
	"welfare-sign/internal/pkg/log"
	"welfare-sign/internal/pkg/util"
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
	if vo.Status != "" {
		query["status"] = vo.Status
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
func (s *Service) GetCustomerDetail(ctx context.Context, customerID uint64) (*model.Customer, wsgin.APICode, error) {
	customer, err := s.dao.FindCustomer(ctx, map[string]interface{}{
		"id":     customerID,
		"status": global.ActiveStatus,
	})
	if err != nil {
		return nil, apicode.ErrDetail, err
	}
	if customer.ID == 0 {
		return nil, wsgin.APICodeSuccess, nil
	}
	return customer, wsgin.APICodeSuccess, nil
}

// GetCustomerCheckinRecord 获取用户签到记录
func (s *Service) GetCustomerCheckinRecord(ctx context.Context, customerID, inCustomerID uint64) ([]*model.CheckinRecord, wsgin.APICode, error) {
	if inCustomerID != 0 {
		records, _ := s.dao.ListCheckinRecord(ctx, "status <> ? AND customer_id = ?", global.DeleteStatus, inCustomerID)
		return records, wsgin.APICodeSuccess, nil
	} else {
		records, err := s.dao.ListCheckinRecord(ctx, "status <> ? AND customer_id = ?", global.DeleteStatus, customerID)
		if err != nil {
			return nil, apicode.ErrGetCheckinRecord, err
		}
		// 用户无签到记录时，自动创建5条信息并返回
		if len(records) == 0 {
			records, err = s.dao.InitCheckinRecords(ctx, customerID)
			if err != nil {
				return nil, apicode.ErrGetCheckinRecord, err
			}
		}
		return records, wsgin.APICodeSuccess, nil
	}
}

// CustomerLogin 客户登录
func (s *Service) CustomerLogin(ctx context.Context, c string) (string, wsgin.APICode, error) {
	var (
		successResp  model.WxSuccessResp
		errResp      model.WxErrResp
		userinfoResp model.WxUserResp
	)
	// 使用code获取access_token
	resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", viper.GetString(config.KeyWxAppID), viper.GetString(config.KeyWxAppSecret), c))
	if err != nil {
		return "", apicode.ErrLogin, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", apicode.ErrLogin, err
	}
	if err := json.Unmarshal(bytes, &errResp); err != nil {
		return "", apicode.ErrLogin, err
	}
	if errResp.Errcode != 0 {
		log.Warn(ctx, "(CustomerLogin)get access_token error", zap.Error(err))
		return "", apicode.ErrLogin, err
	}
	if err := json.Unmarshal(bytes, &successResp); err != nil {
		return "", apicode.ErrLogin, err
	}

	// 查看该用户是否被禁用
	disableCustomer, err := s.dao.FindCustomer(ctx, map[string]interface{}{
		"open_id": successResp.OpenID,
		"status":  global.DeleteStatus,
	})
	if err != nil {
		return "", apicode.ErrLogin, err
	}
	if disableCustomer.ID != 0 {
		return "", apicode.ErrLogin, errors.New("用户已被禁用")
	}

	// 使用access_token + openid获取用户信息
	resp, err = http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", successResp.AccessToken, successResp.OpenID))
	if err != nil {
		return "", apicode.ErrLogin, err
	}
	defer resp.Body.Close()
	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", apicode.ErrLogin, err
	}
	if err := json.Unmarshal(bytes, &errResp); err != nil {
		return "", apicode.ErrLogin, err
	}
	if errResp.Errcode != 0 {
		log.Warn(ctx, "(CustomerLogin)get userinfo error", zap.Error(err))
		return "", apicode.ErrLogin, err
	}
	if err := json.Unmarshal(bytes, &userinfoResp); err != nil {
		return "", apicode.ErrLogin, err
	}
	customer := &model.Customer{}
	if err := util.StructCopy(customer, &userinfoResp); err != nil {
		return "", apicode.ErrLogin, err
	}

	customer, err = s.dao.UpsertCustomer(ctx, &userinfoResp)
	if err != nil {
		return "", apicode.ErrLogin, err
	}
	token, err := jwt.CreateToken(customer.ID, customer.Name, customer.Mobile)
	if err != nil {
		log.Info(ctx, "CustomerLogin.CreateToken() error", zap.Error(err))
		return "", apicode.ErrLogin, err
	}
	return token, wsgin.APICodeSuccess, nil
}

// CustomerNearMerchant 获取用户附近最近的几家店铺
func (s *Service) CustomerNearMerchant(ctx context.Context, data *model.NearMerchantVO) ([]*model.Merchant, wsgin.APICode, error) {
	merchants, err := s.dao.NearMerchant(ctx, data)
	if err != nil {
		return nil, apicode.ErrGetNearMerchant, err
	}
	return merchants, wsgin.APICodeSuccess, nil
}

// ExecCheckinRecord 用户签到
func (s *Service) ExecCheckinRecord(ctx context.Context, customerID uint64) (wsgin.APICode, error) {
	hasChecked, err := s.dao.HasChecked(ctx, customerID)
	if err != nil {
		return apicode.ErrExecCheckinRecord, err
	}
	if hasChecked {
		return apicode.ErrHasCheckin, errors.New("has checkin")
	}

	unchecked, err := s.dao.GetUnchecked(ctx, customerID)
	if err != nil {
		return apicode.ErrExecCheckinRecord, err
	}
	if unchecked.ID != 0 {
		return apicode.ErrExecCheckinRecord, errors.New("请先完成补签后再来签到")
	}

	if err = s.dao.ExecCheckin(ctx, customerID); err != nil {
		log.Warn(ctx, "用户签到发生错误: 更新记录时发生错误", zap.Error(err))
		return wsgin.APICodeServerError, errors.New("用户签到发生错误")
	}
	return wsgin.APICodeSuccess, nil
}

// GetQRCode 客户获取二维码
func (s *Service) GetQRCode(ctx context.Context, customerID uint64) (data []byte, err error) {
	data, err = qrcode.Encode(fmt.Sprintf(viper.GetString(config.KeyQRCodeURL), strconv.FormatUint(customerID, 10)), qrcode.Medium, 256)
	if err != nil {
		log.Warn(ctx, "客户获取二维码失败", zap.Error(err))
	}
	return
}

// GetIssueRecords 客户查看我的福利
func (s *Service) GetIssueRecords(ctx context.Context, customerID uint64) ([]*model.IssueRecord, wsgin.APICode, error) {
	issueRecords, err := s.dao.ListIssueRecordDetail(ctx, map[string]interface{}{
		"customer_id": customerID,
		"status":      global.ActiveStatus,
	})
	if err != nil {
		return nil, apicode.ErrIssueRecord, err
	}
	return issueRecords, wsgin.APICodeSuccess, nil
}

// ExecIssueRecords 客户领取福利
func (s *Service) ExecIssueRecords(ctx context.Context, customerID, merchantID uint64) (wsgin.APICode, error) {
	checkinRecords, err := s.dao.ListCheckinRecord(ctx, map[string]interface{}{
		"status":      global.ActiveStatus,
		"customer_id": customerID,
	})
	if err != nil {
		log.Warn(ctx, "ExecIssueRecords.ListCheckinRecord() error", zap.Error(err))
		return apicode.ErrExecIssueRecord, err
	}
	if len(checkinRecords) != 5 {
		return apicode.ErrExecIssueRecord, errors.New("您还未签满5天")
	}

	merchant, err := s.dao.FindMerchant(ctx, map[string]interface{}{
		"merchant_id": merchantID,
		"status":      global.ActiveStatus,
	})
	if err != nil {
		return apicode.ErrExecIssueRecord, err
	}
	if merchant.ID == 0 {
		return apicode.ErrExecIssueRecord, errors.New("该商户已被禁用")
	}
	if merchant.Received >= merchant.TotalReceive {
		return apicode.ErrExecIssueRecord, errors.New("该商家的福利已被领完了")
	}

	var issueRecord model.IssueRecord
	issueRecord.MerchantID = merchantID
	issueRecord.CustomerID = customerID
	issueRecord.TotalReceive = merchant.CheckinNum
	issueRecord.Received = 0
	if err := s.dao.CreateIssueRecord(ctx, issueRecord); err != nil {
		return apicode.ErrExecIssueRecord, err
	}
	return wsgin.APICodeSuccess, nil
}

// RefreshCheckinRecord 客户重新签到
func (s *Service) RefreshCheckinRecord(ctx context.Context, customerID uint64) (wsgin.APICode, error) {
	if err := s.dao.InvalidCheckin(ctx, customerID); err != nil {
		return apicode.ErrRefreshCheckinRecord, err
	}
	return wsgin.APICodeSuccess, nil
}

// HelpCheckinRecord 帮助他人签到
func (s *Service) HelpCheckinRecord(ctx context.Context, helpCustomerID, customerID uint64) (wsgin.APICode, error) {
	customer, err := s.dao.FindCustomer(ctx, map[string]interface{}{
		"id":     helpCustomerID,
		"status": global.ActiveStatus,
	})
	if err != nil {
		return apicode.ErrHelpCheckin, err
	}
	if customer.ID == 0 {
		return apicode.ErrHelpCheckin, errors.New("帮签用户不存在")
	}

	hasHelpChecked, err := s.dao.FindCheckinRecord(ctx, map[string]interface{}{
		"customer_id":              customerID,
		"help_checkin_customer_id": helpCustomerID,
		"status":                   global.ActiveStatus,
	})
	if err != nil {
		return apicode.ErrHelpCheckin, err
	}
	if hasHelpChecked.ID != 0 {
		return apicode.ErrHasHelpCheckin, errors.New("has help checkin")
	}

	unChecked, err := s.dao.GetUnchecked(ctx, customerID)
	if err != nil {
		return apicode.ErrHelpCheckin, err
	}
	if unChecked.ID == 0 {
		return apicode.ErrHelpCheckin, errors.New("该用户没有需要补签的记录")
	}

	if err := s.dao.HelpCheckin(ctx, unChecked.ID, customerID, helpCustomerID); err != nil {
		log.Warn(ctx, "帮签发生错误", zap.Error(err))
		return apicode.ErrHelpCheckin, err
	}
	return wsgin.APICodeSuccess, nil
}

// IsSupplement 是否是补签
func (s *Service) IsSupplement(ctx context.Context, customerID uint64) (bool, wsgin.APICode, error) {
	records, err := s.dao.ListCheckinRecord(ctx, map[string]interface{}{
		"customer_id": customerID,
		"status":      global.ActiveStatus,
	})
	if err != nil {
		return false, apicode.ErrGetIsSupplement, err
	}
	if len(records) == 0 {
		return false, apicode.ErrGetIsSupplement, err
	}
	lastRecord := records[len(records)-1]
	if lastRecord.HelpCheckinCustomerID != 0 {
		return true, wsgin.APICodeSuccess, nil
	}
	orderRecord, err := s.dao.FindWXPayRecord(ctx, map[string]interface{}{
		"checkin_record_id": lastRecord.ID,
		"status":            global.ActiveStatus,
	})
	if err != nil {
		return false, apicode.ErrGetIsSupplement, err
	}
	if orderRecord.ID != 0 {
		return true, wsgin.APICodeSuccess, nil
	}
	return false, wsgin.APICodeSuccess, nil
}

// DisableCustomer 禁用客户
func (s *Service) DisableCustomer(ctx context.Context, customerID uint64) (wsgin.APICode, error) {
	customer, err := s.dao.FindCustomer(ctx, map[string]interface{}{
		"id": customerID,
	})
	if err != nil || customer.ID == 0 {
		return apicode.ErrDisable, err
	}
	if customer.Status == global.DeleteStatus {
		return apicode.ErrHasDisable, errors.New("用户已经被禁用")
	}
	customer.Status = global.DeleteStatus
	if err := s.dao.UpdateCustomer(ctx, customer); err != nil {
		return apicode.ErrDisable, err
	}
	return wsgin.APICodeSuccess, nil
}

// DeleteCustomer 删除客户
func (s *Service) DeleteCustomer(ctx context.Context, customerID uint64) (wsgin.APICode, error) {
	s.dao.DeleteCustomer(ctx, customerID)
	return wsgin.APICodeSuccess, nil
}
