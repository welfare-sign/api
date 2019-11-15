package dao

import (
	"context"
	"time"

	"go.uber.org/zap"

	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/log"
)

// ListCheckinRecord 获取签到列表
func (d *dao) ListCheckinRecord(ctx context.Context, query interface{}, args ...interface{}) ([]*model.CheckinRecord, error) {
	var res []*model.CheckinRecord
	err := checkErr(d.db.Where(query, args...).Find(&res).Order("CreatedAt ASC").Error)
	return res, err
}

// InitCheckinRecords 用户第一次访问签到页面时，初始化签到信息并返回
func (d *dao) InitCheckinRecords(ctx context.Context, customerID uint64) ([]*model.CheckinRecord, error) {
	tx := d.db.Begin()

	var res []*model.CheckinRecord
	// 目前限制签到5天，只创建5条记录
	for i := 0; i < 5; i++ {
		cr := &model.CheckinRecord{}
		cr.Status = global.InactiveStatus
		cr.UpdatedAt = time.Now()
		cr.CreatedAt = time.Now()
		cr.CustomerID = customerID
		cr.Day = uint64(i) + 1
		res = append(res, cr)
		if err := tx.Create(cr).Error; err != nil {
			log.Warn(ctx, "dao.InitCheckinRecords() error", zap.Error(err))
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return res, nil
}

// FindCheckinRecord 查询签到记录
func (d *dao) FindCheckinRecord(ctx context.Context, query interface{}, args ...interface{}) (*model.CheckinRecord, error) {
	var checkinRecord model.CheckinRecord
	err := checkErr(d.db.Where(query, args...).First(&checkinRecord).Error)
	return &checkinRecord, err
}

// ExecCheckin 记录用户签到
func (d *dao) ExecCheckin(ctx context.Context, customerID, day uint64) error {
	return d.db.Model(&model.CheckinRecord{}).Where(map[string]interface{}{
		"status":      global.InactiveStatus,
		"customer_id": customerID,
		"day":         day,
	}).Update("status", global.ActiveStatus).Error
}

// InvalidCheckin 作废用户签到记录
func (d *dao) InvalidCheckin(ctx context.Context, customerID uint64) error {
	return checkErr(d.db.Model(&model.CheckinRecord{}).Where(map[string]interface{}{
		"status":      global.ActiveStatus,
		"customer_id": customerID,
	}).Update("status", global.InactiveStatus).Error)
}

// HelpCheckin 帮助他人补签
func (d *dao) HelpCheckin(ctx context.Context, customerID, helpCustomerID, day uint64) error {
	return checkErr(d.db.Model(&model.CheckinRecord{}).Where(map[string]interface{}{
		"status":      global.InactiveStatus,
		"customer_id": customerID,
		"day":         day,
	}).Updates(map[string]interface{}{
		"status":                   global.ActiveStatus,
		"help_checkin_customer_id": helpCustomerID,
	}).Error)
}
