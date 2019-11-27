package dao

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/log"
)

const (
	hasCheckedSQL = `
SELECT * from checkin_record WHERE CURRENT_DATE = DATE(need_checkin_time) AND customer_id = ? 
AND status = ?
`
	getUncheckedSQL = `
SELECT * from checkin_record WHERE DATE(need_checkin_time) < CURRENT_DATE AND customer_id = ? 
AND status = ?
	`
	execCheckinSQL = `
	UPDATE checkin_record SET status = ?, updated_at = ? WHERE DATE(need_checkin_time) = CURRENT_DATE AND status = ? AND customer_id = ?
	`
	helpCheckinSQL = `
	UPDATE checkin_record SET status = ?, updated_at = ?, help_checkin_customer_id = ? WHERE id = ?
	`
	payCheckinSQL = `
	UPDATE checkin_record SET status = ?, updated_at = ? WHERE id = ?
	`
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
		cr.NeedCheckinTime = time.Now().Add(time.Duration(i) * 24 * time.Hour)
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
func (d *dao) ExecCheckin(ctx context.Context, customerID uint64) error {
	tx := d.db.Begin()

	if err := tx.Exec(execCheckinSQL, global.ActiveStatus, time.Now(), global.InactiveStatus, customerID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&model.Customer{}).Where(map[string]interface{}{
		"id":     customerID,
		"status": global.ActiveStatus,
	}).Update("last_checkin_time", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// InvalidCheckin 作废用户签到记录
func (d *dao) InvalidCheckin(ctx context.Context, customerID uint64) error {
	return checkErr(d.db.Model(&model.CheckinRecord{}).Where("status <> ? AND customer_id = ?", global.DeleteStatus, customerID).Update("status", global.DeleteStatus).Error)
}

// HelpCheckin 帮助他人补签
func (d *dao) HelpCheckin(ctx context.Context, checkRecordID, customerID, helpCustomerID uint64) error {
	tx := d.db.Begin()

	if err := tx.Exec(helpCheckinSQL, global.ActiveStatus, time.Now(), helpCustomerID, checkRecordID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&model.Customer{}).Where(map[string]interface{}{
		"id":     customerID,
		"status": global.ActiveStatus,
	}).Update("last_checkin_time", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}
	msg := model.HelpCheckinMessage{}
	msg.SetDefaultAttr()
	msg.CheckinRecordID = checkRecordID
	msg.CustomerID = customerID
	msg.IsRead = global.UnRead
	if err := tx.Create(&msg).Error; err != nil {
		log.Warn(ctx, "分享补签时创建补签消息失败", zap.Error(err))
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// HasChecked 用户当前是否已签到
func (d *dao) HasChecked(ctx context.Context, customerID uint64) (bool, error) {
	var checkinRecord model.CheckinRecord
	if err := d.db.Raw(hasCheckedSQL, customerID, global.ActiveStatus).First(&checkinRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	if checkinRecord.ID != 0 {
		return true, nil
	}
	return false, nil
}

// GetUnchecked 获取当天之前未签到的记录
func (d *dao) GetUnchecked(ctx context.Context, customerID uint64) (*model.CheckinRecord, error) {
	var checkinRecord model.CheckinRecord
	err := checkErr(d.db.Raw(getUncheckedSQL, customerID, global.InactiveStatus).First(&checkinRecord).Error)
	return &checkinRecord, err
}

// PayCheckin 用户支付后补签
func (d *dao) PayCheckin(ctx context.Context, checkRecordID, customerID uint64, payRecord *model.WXPayRecord) error {
	tx := d.db.Begin()

	if err := tx.Exec(payCheckinSQL, global.ActiveStatus, time.Now(), checkRecordID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&model.Customer{}).Where(map[string]interface{}{
		"id":     customerID,
		"status": global.ActiveStatus,
	}).Update("last_checkin_time", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}

	payRecord.SetDefaultAttr()
	payRecord.CheckinRecordID = checkRecordID
	if err := tx.Create(payRecord).Error; err != nil {
		log.Warn(ctx, "dao.PayCheckin.Create.WXPayRecord error", zap.Error(err))
		tx.Rollback()
		return err
	}

	msg := model.HelpCheckinMessage{}
	msg.SetDefaultAttr()
	msg.CheckinRecordID = checkRecordID
	msg.CustomerID = customerID
	msg.IsRead = global.UnRead
	if err := tx.Create(&msg).Error; err != nil {
		log.Warn(ctx, "支付回调时创建补签消息失败", zap.Error(err))
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
