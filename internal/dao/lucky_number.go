package dao

import (
	"context"

	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/util"
)

const (
	getLuckyNumberRecordSQL = `
SELECT *FROM lucky_number_record WHERE status = ? AND customer_id = ? and created_at BETWEEN 
DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 2), '%Y-%m-%d 00:00:00') AND
DATE_FORMAT( DATE_ADD(SUBDATE(CURDATE(), WEEKDAY(CURDATE())), INTERVAL 4 DAY), '%Y-%m-%d 12:00:00')
`
	existsNumSQL = `
	SELECT *FROM lucky_number_record WHERE status = ? AND lucky_number = ? and created_at BETWEEN 
DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 2), '%Y-%m-%d 00:00:00') AND
DATE_FORMAT( DATE_ADD(SUBDATE(CURDATE(), WEEKDAY(CURDATE())), INTERVAL 4 DAY), '%Y-%m-%d 12:00:00')
	`
	notAvailableNumSQL = `
	SELECT DISTINCT lucky_number as nums FROM lucky_number_record WHERE status = ? and created_at BETWEEN 
DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 2), '%Y-%m-%d 00:00:00') AND
DATE_FORMAT( DATE_ADD(SUBDATE(CURDATE(), WEEKDAY(CURDATE())), INTERVAL 4 DAY), '%Y-%m-%d 12:00:00')
	`
	getLuckyNumberRecordBefore = `
SELECT *FROM lucky_number_record WHERE status = ? AND customer_id = ? and created_at BETWEEN 
DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 9), '%Y-%m-%d 00:00:00') AND
DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 3), '%Y-%m-%d 12:00:00')
	`
	getLuckyPeopleBefore = `
SELECT *FROM lucky_number_record WHERE status = ? AND ranking = 1 and created_at BETWEEN 
DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 9), '%Y-%m-%d 00:00:00') AND
DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 3), '%Y-%m-%d 12:00:00')
	`
	getRoundLuckyNumberRecordSQL = `
	SELECT * FROM lucky_number_record WHERE status = ? and created_at BETWEEN 
DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 2), '%Y-%m-%d 00:00:00') AND
DATE_FORMAT( DATE_ADD(SUBDATE(CURDATE(), WEEKDAY(CURDATE())), INTERVAL 4 DAY), '%Y-%m-%d 12:00:00')
	`
)

// GetRoundLuckyNumberRecordSQL 获取本轮内所有有效的数字
func (d *dao) GetRoundLuckyNumberRecordSQL(ctx context.Context) ([]*model.LuckyNumberRecord, error) {
	var luckyList []*model.LuckyNumberRecord
	err := checkErr(d.db.Raw(getRoundLuckyNumberRecordSQL, global.ActiveStatus).Find(&luckyList).Error)
	return luckyList, err
}

// GetLuckyNumberRecord 获取用户在本轮内猜的数字
func (d *dao) GetLuckyNumberRecord(ctx context.Context, customerID uint64) (*model.LuckyNumberRecord, error) {
	var lucky model.LuckyNumberRecord
	err := checkErr(d.db.Raw(getLuckyNumberRecordSQL, global.ActiveStatus, customerID).First(&lucky).Error)
	if err != nil || lucky.ID == 0 {
		return &lucky, err
	}
	lucky.CompositeIndex, _ = d.GetCompositeIndex(ctx)
	return &lucky, nil
}

// StoreLuckyNumberRecord 存储用户在本轮内猜的数字
func (d *dao) StoreLuckyNumberRecord(ctx context.Context, customerID uint64, num int64) ([]int64, error) {
	var (
		lucky   model.LuckyNumberRecord
		allNums []int64
	)
	availableNum := make([]int64, 0, 2)

	err := checkErr(d.db.Raw(existsNumSQL, global.ActiveStatus, num).First(&lucky).Error)
	if err != nil {
		return availableNum, err
	}
	if lucky.ID == 0 {
		lucky.SetDefaultAttr()
		lucky.CustomerID = customerID
		lucky.LuckyNumber = num
		err := d.db.Create(&lucky).Error
		return availableNum, err
	}
	d.db.Raw(notAvailableNumSQL, global.ActiveStatus).Pluck("nums", &allNums)
	return util.GetRecommandNum(num, allNums)
}

// GetLuckyNumberRecordBefore 获取用户在上轮内猜的数字
func (d *dao) GetLuckyNumberRecordBefore(ctx context.Context, customerID uint64) (*model.LuckyNumberRecord, error) {
	var lucky model.LuckyNumberRecord
	err := checkErr(d.db.Raw(getLuckyNumberRecordBefore, global.ActiveStatus, customerID).First(&lucky).Error)
	if err != nil || lucky.ID == 0 {
		return &lucky, err
	}
	lucky.CompositeIndex, _ = d.GetCompositeIndexBefore(ctx)
	return &lucky, nil
}

// GetLuckyPeopleBefore 获取上期幸运用户
func (d *dao) GetLuckyPeopleBefore(ctx context.Context) (*model.Customer, error) {
	var lucky model.LuckyNumberRecord
	err := checkErr(d.db.Raw(getLuckyPeopleBefore, global.ActiveStatus).First(&lucky).Error)
	if err != nil || lucky.ID == 0 {
		return nil, err
	}
	customer, err := d.FindCustomer(ctx, map[string]interface{}{
		"id":     lucky.CustomerID,
		"status": global.ActiveStatus,
	})
	if err != nil {
		return nil, err
	}
	return customer, nil
}
