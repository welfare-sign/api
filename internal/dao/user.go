package dao

import (
	"context"

	"welfare-sign/internal/model"
)

const (
	getTmpCheckinRecordListSQL = `
	SELECT *FROM checkin_record WHERE status <> 'X' AND DATE(need_checkin_time) <= CURRENT_DATE
	`
)

// FindUser find user
func (d *dao) FindUser(ctx context.Context, query interface{}) (*model.User, error) {
	var user model.User
	err := d.db.Where(query).Find(&user).Error
	return &user, err
}

// GetTmpCheckinRecordList 获取全部用户签到列表
// TODO: 临时
func (d *dao) GetTmpCheckinRecordList(ctx context.Context) ([]*model.CheckinRecordListResp, error) {
	var (
		checkinRecords []*model.CheckinRecord
		data           = make([]*model.CheckinRecordListResp, 0, 30)
	)
	err := d.db.Raw(getTmpCheckinRecordListSQL).Find(&checkinRecords).Error
	if err != nil {
		return data, err
	}
	for i := 0; i < len(checkinRecords); i++ {
		var checkinRecordListResp model.CheckinRecordListResp
		checkinRecordListResp.CheckinRecord = checkinRecords[i]
		checkinRecordListResp.Customer, _ = d.FindCustomer(ctx, map[string]interface{}{"id": checkinRecords[i].CustomerID})
		data = append(data, &checkinRecordListResp)
	}
	return data, nil
}

// UpdateCustomerCheckinRecord 更新用户签到记录
// TODO: 临时
func (d *dao) UpdateCustomerCheckinRecord(ctx context.Context, checkinRecord uint64, status string) error {
	if err := d.db.Model(&model.CheckinRecord{}).Where("id = ?", checkinRecord).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
