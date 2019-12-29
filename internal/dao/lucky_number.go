package dao

import (
	"context"
	"time"

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
	notAvailableNumSQL = `
	SELECT DISTINCT lucky_number as nums FROM lucky_number_record WHERE status = ? and activity_id = ?
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

// GetLuckyNumberRecord 获取用户在本轮内猜的数字
func (d *dao) GetLuckyNumberRecord(ctx context.Context, customerID, activityID uint64) (*model.LuckyNumberRecord, error) {
	var lucky model.LuckyNumberRecord
	err := checkErr(d.db.Where(map[string]interface{}{
		"activity_id": activityID,
		"customer_id": customerID,
		"status":      global.ActiveStatus,
	}).First(&lucky).Error)
	if err != nil {
		return &lucky, err
	}

	customer, err := d.FindCustomer(ctx, map[string]interface{}{
		"id":     customerID,
		"status": global.ActiveStatus,
	})
	if err != nil {
		return &lucky, err
	}
	lucky.Customer = customer

	activity, err := d.FindActivity(ctx, map[string]interface{}{
		"id":     activityID,
		"status": global.ActiveStatus,
	})
	if err != nil {
		return &lucky, err
	}
	lucky.Activity = activity
	return &lucky, nil
}

// StoreLuckyNumberRecord 存储用户在本轮内猜的数字
func (d *dao) StoreLuckyNumberRecord(ctx context.Context, customerID, activityID, num uint64) ([]uint64, error) {
	var (
		lucky    model.LuckyNumberRecord
		customer model.Customer
		allNums  []uint64
	)
	availableNum := make([]uint64, 0, 2)

	err := checkErr(d.db.Where(map[string]interface{}{
		"activity_id":  activityID,
		"lucky_number": num,
		"status":       global.ActiveStatus,
	}).First(&lucky).Error)
	if err != nil {
		return availableNum, err
	}
	err = checkErr(d.db.Where(map[string]interface{}{
		"id":     customerID,
		"status": global.ActiveStatus,
	}).First(&customer).Error)
	if err != nil {
		return availableNum, err
	}
	if lucky.ID == 0 {
		lucky.SetDefaultAttr()
		lucky.CustomerID = customerID
		lucky.LuckyNumber = num
		lucky.Mobile = customer.Mobile
		lucky.ActivityID = activityID
		err := d.db.Create(&lucky).Error
		return availableNum, err
	}
	d.db.Raw(notAvailableNumSQL, global.ActiveStatus).Pluck("nums", &allNums)
	return util.GetRecommandNum(num, allNums)
}

// GetLuckyNumberRecordBefore 获取用户在上轮内猜的数字
func (d *dao) GetLuckyNumberRecordBefore(ctx context.Context, customerID uint64) (*model.LuckyNumberRecord, error) {
	var (
		activity model.Activity
		lucky    model.LuckyNumberRecord
	)
	// 根据时间，最近结束的那轮，且是已开奖
	err := checkErr(d.db.Where("status = ? and end_time < ? and prize_number != 0", global.ActiveStatus, time.Now()).Order("end_time desc").First(&activity).Error)
	if err != nil {
		return nil, err
	}

	err = checkErr(d.db.Where(map[string]interface{}{
		"activity_id": activity.ID,
		"customer_id": customerID,
		"status":      global.ActiveStatus,
	}).First(&lucky).Error)
	if err != nil {
		return nil, err
	}
	lucky.Activity = &activity
	return &lucky, nil
}

// GetLuckyPeopleBefore 获取上期幸运用户
func (d *dao) GetLuckyPeopleBefore(ctx context.Context) ([]*model.LuckyNumberRecord, error) {
	var (
		activity model.Activity
		luckys   []*model.LuckyNumberRecord
	)
	// 根据时间，最近结束的那轮，且是已开奖
	err := checkErr(d.db.Where("status = ? and end_time < ? and prize_number != 0", global.ActiveStatus, time.Now()).Order("end_time desc").First(&activity).Error)
	if err != nil {
		return nil, err
	}

	err = checkErr(d.db.Where("activity_id = ? and status = ? and ranking between ? and ?", activity.ID, global.ActiveStatus, 1, activity.PrizeAmount).Order("ranking asc").Find(&luckys).Error)
	if err != nil {
		return nil, err
	}

	if len(luckys) == 0 {
		return luckys, nil
	}
	for i := 0; i < len(luckys); i++ {
		luckys[i].Activity = &activity
		customer, err := d.FindCustomer(ctx, map[string]interface{}{
			"id":     luckys[i].CustomerID,
			"status": global.ActiveStatus,
		})
		if err != nil {
			continue
		}
		luckys[i].Customer = customer
	}

	return luckys, nil
}
