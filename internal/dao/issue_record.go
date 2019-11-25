package dao

import (
	"context"

	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/log"

	"go.uber.org/zap"
)

// FindIssueRecord 获取礼包发放记录详情
func (d *dao) FindIssueRecord(ctx context.Context, query interface{}) (*model.IssueRecord, error) {
	var issueRecord model.IssueRecord
	err := checkErr(d.db.Where(query).First(&issueRecord).Error)
	return &issueRecord, err
}

// ListIssueRecord 获取礼包发放记录列表
func (d *dao) ListIssueRecord(ctx context.Context, query interface{}, args ...interface{}) ([]*model.IssueRecord, error) {
	var res []*model.IssueRecord
	err := checkErr(d.db.Where(query, args...).Find(&res).Order("CreatedAt DESC").Error)
	return res, err
}

// ListIssueRecordDetail 获取礼包发放记录列表，携带商家、客户信息
func (d *dao) ListIssueRecordDetail(ctx context.Context, query interface{}, args ...interface{}) ([]*model.IssueRecord, error) {
	var res []*model.IssueRecord
	if err := checkErr(d.db.Where(query, args...).Find(&res).Order("CreatedAt DESC").Error); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return res, nil
	}

	for i := 0; i < len(res); i++ {
		merchant, _ := d.FindMerchant(ctx, map[string]interface{}{
			"id":     res[i].MerchantID,
			"status": global.ActiveStatus,
		})
		if merchant.ID != 0 {
			res[i].Merchant = merchant
		}
		customer, _ := d.FindCustomer(ctx, map[string]interface{}{
			"id":     res[i].CustomerID,
			"status": global.ActiveStatus,
		})
		if customer.ID != 0 {
			res[i].Customer = customer
		}
	}

	return res, nil
}

// CreateIssueRecord create issue record
func (d *dao) CreateIssueRecord(ctx context.Context, data model.IssueRecord, merchant *model.Merchant, mobile string) error {
	tx := d.db.Begin()

	data.SetDefaultAttr()
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		log.Warn(ctx, "CreateIssueRecord.Create() error", zap.Error(err))
		return err
	}
	if err := tx.Model(&model.CheckinRecord{}).Where(map[string]interface{}{
		"status":      global.ActiveStatus,
		"customer_id": data.CustomerID,
	}).Update("status", global.InactiveStatus).Error; err != nil {
		tx.Rollback()
		log.Warn(ctx, "CreateIssueRecord.Update() error", zap.Error(err))
		return err
	}
	if err := tx.Save(merchant).Error; err != nil {
		log.Warn(ctx, "CreateIssueRecord.SaveMerchant() error", zap.Error(err))
		tx.Rollback()
		return err
	}
	if mobile != "" {
		if err := tx.Model(&model.Customer{}).Where(map[string]interface{}{
			"status": global.ActiveStatus,
			"id":     data.CustomerID,
		}).Update("mobile", mobile).Error; err != nil {
			tx.Rollback()
			log.Warn(ctx, "CreateIssueRecord.UpdateCustomer() error", zap.Error(err))
			return err
		}
	}

	tx.Commit()
	return nil
}
