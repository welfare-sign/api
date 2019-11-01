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
func (d *dao) ListCheckinRecord(ctx context.Context, query interface{}) ([]*model.CheckinRecord, error) {
	var res []*model.CheckinRecord
	err := d.db.Where(query).Find(&res).Error
	return res, err
}

// InitCheckinRecords 用户第一次访问签到页面时，初始化签到信息并返回
func (d *dao) InitCheckinRecords(ctx context.Context, customerId uint64) ([]*model.CheckinRecord, error) {
	tx := d.db.Begin()

	var res []*model.CheckinRecord
	// 目前限制签到5天，只创建5条记录
	for i := 0; i < 5; i++ {
		cr := &model.CheckinRecord{}
		cr.Status = global.CheckinRecordInactiveStatus
		cr.UpdatedAt = time.Now()
		cr.CreatedAt = time.Now()
		cr.CustomerId = customerId
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
