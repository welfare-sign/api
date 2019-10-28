package service

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/jwt"
	"welfare-sign/internal/pkg/log"
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

// MerchantLogin 商家登录
func (s *Service) MerchantLogin(ctx context.Context, vo model.MerchantLoginVO) (string, error) {
	merchant, err := s.dao.FindMerchant(ctx, map[string]interface{}{"contact_phone": vo.ContactPhone})
	if err != nil {
		log.Info(ctx, "MerchantLogin.FindMerchant() error", zap.Error(err))
		return "", err
	}
	if err := s.ValidateCode(ctx, vo.ContactPhone, vo.Code); err != nil {
		log.Info(ctx, "MerchantLogin.ValidateCode() error", zap.Error(err))
		return "", err
	}
	return jwt.CreateToken(merchant.ContactName, merchant.ContactPhone)
}

// GetMerchantDetailBySelfAccessToken 获取商户详情
func (s *Service) GetMerchantDetailBySelfAccessToken(ctx context.Context, token string) (*model.Merchant, error) {
	tokenParames, err := jwt.ParseToken(token)
	if err != nil {
		return nil, err
	}
	merchant, err := s.dao.FindMerchant(ctx, map[string]interface{}{"contact_phone": tokenParames.Mobile})
	if err != nil {
		return nil, err
	}
	return merchant, nil
}

// GetWriteOff 获取核销页面数据
func (s *Service) GetWriteOff(ctx context.Context, vo model.MerchantWriteOffVO) (*model.MerchantWriteOffRespVO, error) {
	var resp model.MerchantWriteOffRespVO

	tokenParames, err := jwt.ParseToken(vo.AccessToken)
	if err != nil {
		return nil, err
	}
	merchant, err := s.dao.FindMerchant(ctx, map[string]interface{}{"contact_phone": tokenParames.Mobile})
	if err != nil {
		log.Info(ctx, "GetWriteOff.FindMerchant() error", zap.Error(err))
		return nil, err
	}
	customer, err := s.dao.FindCustomer(ctx, map[string]interface{}{"id": vo.CustomerID})
	if err != nil {
		log.Info(ctx, "GetWriteOff.FindCustomer() error", zap.Error(err))
		return nil, err
	}
	issueRecord, err := s.dao.FindIssueRecord(ctx, map[string]interface{}{
		"merchant_id": merchant.ID,
		"customer_id": customer.ID,
	})
	if err != nil {
		log.Info(ctx, "GetWriteOff.FindIssueRecord() error", zap.Error(err))
		return nil, nil
	}
	resp.Merchant = merchant
	resp.Customer = customer
	resp.IssueRecord = issueRecord
	return &resp, nil
}

// ExecWriteOff 执行核销
func (s *Service) ExecWriteOff(ctx context.Context, vo model.MerchantExecWriteOffVO) (*model.MerchantWriteOffRespVO, error) {
	resp, err := s.GetWriteOff(ctx, model.MerchantWriteOffVO{AccessToken: vo.AccessToken, CustomerID: vo.CustomerID})
	if err != nil {
		return nil, err
	}
	if (resp.IssueRecord.TotalReceive - resp.IssueRecord.Received) < vo.Num {
		return nil, errors.New("核销数目不正确")
	}
	hasRece := resp.IssueRecord.Received + vo.Num
	totalRece := resp.Merchant.Received + vo.Num
	resp.IssueRecord.Received = resp.IssueRecord.TotalReceive - hasRece
	return resp, s.dao.EcecWriteOff(ctx, resp.Merchant.ID, resp.Customer.ID, hasRece, totalRece)
}
