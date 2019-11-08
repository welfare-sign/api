package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/config"
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
	customer, err := s.dao.FindCustomer(ctx, map[string]interface{}{"id": customerID})
	if err != nil {
		return nil, apicode.ErrDetail, err
	}
	return customer, wsgin.APICodeSuccess, nil
}

// GetCustomerCheckinRecord 获取用户签到记录
func (s *Service) GetCustomerCheckinRecord(ctx context.Context, customerID uint64) ([]*model.CheckinRecord, wsgin.APICode, error) {
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

// CustomerLogin 客户登录
func (s *Service) CustomerLogin(ctx context.Context, c string) (*model.CustomerLoginResp, wsgin.APICode, error) {
	var (
		successResp  model.WxSuccessResp
		errResp      model.WxErrResp
		userinfoResp model.WxUserResp
	)
	// 使用code获取access_token
	resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", viper.GetString(config.KeyWxAppID), viper.GetString(config.KeyWxAppSecret), c))
	if err != nil {
		return nil, apicode.ErrLogin, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, apicode.ErrLogin, err
	}
	if err := json.Unmarshal(bytes, &errResp); err != nil {
		return nil, apicode.ErrLogin, err
	}
	if errResp.Errcode != 0 {
		log.Warn(ctx, "(CustomerLogin)get access_token error", zap.Error(err))
		return nil, apicode.ErrLogin, err
	}
	if err := json.Unmarshal(bytes, &successResp); err != nil {
		return nil, apicode.ErrLogin, err
	}

	// 使用access_token + openid获取用户信息
	resp, err = http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", successResp.AaccessToken, successResp.OpenID))
	if err != nil {
		return nil, apicode.ErrLogin, err
	}
	defer resp.Body.Close()
	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, apicode.ErrLogin, err
	}
	if err := json.Unmarshal(bytes, &errResp); err != nil {
		return nil, apicode.ErrLogin, err
	}
	if errResp.Errcode != 0 {
		log.Warn(ctx, "(CustomerLogin)get userinfo error", zap.Error(err))
		return nil, apicode.ErrLogin, err
	}
	if err := json.Unmarshal(bytes, &userinfoResp); err != nil {
		return nil, apicode.ErrLogin, err
	}
	customer := &model.Customer{}
	if err := util.StructCopy(customer, &userinfoResp); err != nil {
		return nil, apicode.ErrLogin, err
	}

	customer, err = s.dao.UpsertCustomer(ctx, &userinfoResp)
	if err != nil {
		return nil, apicode.ErrLogin, err
	}
	records, _, err := s.GetCustomerCheckinRecord(ctx, customer.ID)
	if err != nil {
		return nil, apicode.ErrLogin, err
	}
	var data model.CustomerLoginResp
	data.Customer = customer
	data.CheckinRecordList = records
	return &data, wsgin.APICodeSuccess, nil
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
func (s *Service) ExecCheckinRecord(ctx context.Context, customerID, day uint64) (wsgin.APICode, error) {
	checkinRecord, err := s.dao.FindCheckinRecord(ctx, map[string]interface{}{
		"customer_id": customerID,
		"day":         day,
		"status":      "A",
	})
	if err != nil {
		return apicode.ErrExecCheckinRecord, err
	}
	if checkinRecord.CustomerID != 0 {
		return apicode.ErrHasCheckin, errors.New("has checkin")
	}
	checkinRecords, err := s.dao.ListCheckinRecord(ctx, "status <> ? AND customer_id = ?", global.DeleteStatus, customerID)
	if err != nil {
		return apicode.ErrExecCheckinRecord, err
	}
	if len(checkinRecords) == 0 || len(checkinRecords) != 5 {
		log.Warn(ctx, "用户签到发生错误: 签到记录为0或不等于5", zap.Error(errors.New("用户签到发生错误")))
		return wsgin.APICodeServerError, errors.New("用户签到发生错误")
	}
	d1 := checkinRecords[0].CreatedAt.AddDate(0, 0, int(day)-1).Format("2006-01-02")
	now := time.Now().Format("2006-01-02")
	if d1 != now {
		return apicode.ErrExecCheckinRecord, errors.New("只可完成当天的签到")
	}
	if err = s.dao.ExecCheckin(ctx, customerID, day); err != nil {
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
