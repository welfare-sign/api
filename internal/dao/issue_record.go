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
	err := checkErr(d.db.Where(query, args...).Find(&res).Related(&model.Merchant{}).Related(&model.Customer{}).Order("CreatedAt DESC").Error)
	return res, err
}

// CreateIssueRecord create issue record
func (d *dao) CreateIssueRecord(ctx context.Context, data model.IssueRecord) error {
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
	tx.Commit()
	return nil
}
