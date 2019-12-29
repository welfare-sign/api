package dao

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
	"welfare-sign/internal/pkg/log"
	"welfare-sign/internal/pkg/util"
)

const (
	getCompositeIndex = `
SELECT *FROM composite_index WHERE status = ? AND DATE(composite_date) = 
DATE(DATE_FORMAT( DATE_ADD(SUBDATE(CURDATE(), WEEKDAY(CURDATE())), INTERVAL 4 DAY), '%Y-%m-%d 12:00:00'))
`
	getCompositeIndexBefore = `
SELECT *FROM composite_index WHERE status = ? AND DATE(composite_date) = DATE(DATE_FORMAT( SUBDATE(CURDATE(), WEEKDAY(CURDATE()) + 3), '%Y-%m-%d 12:00:00'))
`
	isRightDate = `
SELECT DATE(?) = DATE(DATE_FORMAT( DATE_ADD(SUBDATE(CURDATE(), WEEKDAY(CURDATE())), INTERVAL 4 DAY), '%Y-%m-%d 12:00:00')) as date
	`
)

// FindActivity 获取活动详情
func (d *dao) FindActivity(ctx context.Context, query interface{}, args ...interface{}) (*model.Activity, error) {
	var activity model.Activity
	err := checkErr(d.db.Where(query, args...).First(&activity).Error)
	return &activity, err
}

// IsActivityDateLegal 查看活动日期是否合法
func (d *dao) IsActivityDateLegal(ctx context.Context, data *model.ActivityVO) (bool, error) {
	var activity model.Activity
	err := checkErr(d.db.Where("(start_time >= ?) or (start_time >= ? and end_time >= ?) or (start_time <= ? and end_time >= ?) or (start_time <= ? and end_time <= ?)", data.EndTime, data.StartTime, data.EndTime, data.StartTime, data.EndTime, data.StartTime, data.EndTime).First(&activity).Error)
	if err != nil {
		return false, err
	}
	if activity.ID != 0 {
		return false, nil
	}
	return true, nil
}

// UpsertActivity 添加或更新活动信息
func (d *dao) UpsertActivity(ctx context.Context, data *model.ActivityVO) error {
	activity, err := d.FindActivity(ctx, map[string]interface{}{
		"name":   data.Name,
		"status": global.ActiveStatus,
	})
	if checkErr(err) != nil {
		return err
	}
	if activity.ID == 0 {
		var a model.Activity
		util.StructCopy(&a, data)
		a.SetDefaultAttr()
		if err := d.db.Create(&a).Error; err != nil {
			return err
		}
		return nil
	}
	activity.UpdatedAt = time.Now()
	util.StructCopy(&activity, data)
	return d.db.Save(activity).Error
}

// ListActivity 活动列表
func (d *dao) ListActivity(ctx context.Context, query interface{}, pageNo, pageSize int) ([]*model.Activity, int, error) {
	var activitys []*model.Activity
	total := 0
	if err := checkErr(d.db.Where(query).Limit(pageSize).Offset((pageNo - 1) * pageSize).Order("created_at desc").Find(&activitys).Error); err != nil {
		return activitys, total, err
	}
	if err := checkErr(d.db.Model(&model.Activity{}).Where(query).Count(&total).Error); err != nil {
		return activitys, total, err
	}
	for i := 0; i < len(activitys); i++ {
		if activitys[i].PrizeNumber == 0 {
			continue
		}
		c := 0
		if err := checkErr(d.db.Model(&model.LuckyNumberRecord{}).Where("activity_id = ? and ranking between ? and ?", activitys[i].ID, 1, activitys[i].PrizeAmount).Count(&c).Error); err != nil {
			log.Warn(ctx, "计算活动列表中已发放奖励出现错误", zap.Error(err))
			continue
		}
		activitys[i].PrizeIssued = uint64(c)
	}
	return activitys, total, nil
}

// ListActivityParticipant 活动参与者列表
func (d *dao) ListActivityParticipant(ctx context.Context, pageNo, pageSize int, query interface{}, args ...interface{}) ([]*model.LuckyNumberRecord, int, error) {
	var luckys []*model.LuckyNumberRecord
	total := 0
	if err := checkErr(d.db.Where(query, args...).Limit(pageSize).Offset((pageNo - 1) * pageSize).Order("ranking asc,created_at desc").Find(&luckys).Error); err != nil {
		return luckys, total, err
	}
	if err := checkErr(d.db.Model(&model.LuckyNumberRecord{}).Where(query, args...).Count(&total).Error); err != nil {
		return luckys, total, err
	}

	return luckys, total, nil
}

// DrawActivity 活动开奖
func (d *dao) DrawActivity(ctx context.Context, activityID, number uint64) (*model.Activity, error) {
	tx := d.db.Begin()

	var activity model.Activity
	if err := checkErr(tx.Where(map[string]interface{}{
		"id":     activityID,
		"status": global.ActiveStatus,
	}).First(&activity).Error); err != nil {
		tx.Rollback()
		return nil, err
	}
	if activity.ID == 0 {
		tx.Rollback()
		return nil, errors.New("活动未找到或活动已无效")
	}
	activity.UpdatedAt = time.Now()
	activity.PrizeNumber = number
	if err := tx.Save(activity).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// 更新
	/*
		UPDATE mytable SET
		    myfield = CASE id
		        WHEN 1 THEN 'value'
		        WHEN 2 THEN 'value'
		        WHEN 3 THEN 'value'
		    END
		WHERE id IN (1,2,3)
	*/
	var luckyNumberRecords []*model.LuckyNumberRecord
	if err := checkErr(tx.Where(map[string]interface{}{
		"activity_id": activityID,
		"status":      global.ActiveStatus,
	}).Find(&luckyNumberRecords).Error); err != nil {
		tx.Rollback()
		return nil, err
	}
	// 没有人参与时，直接返回
	if len(luckyNumberRecords) == 0 {
		tx.Commit()
		return &activity, nil
	}
	// 有人参与
	recordList := make(util.SortRecordList, 0, len(luckyNumberRecords))
	ids := make([]string, 0, len(luckyNumberRecords))
	for _, r := range luckyNumberRecords {
		recordList = append(recordList, util.SortRecord{
			ID:  r.ID,
			Num: r.LuckyNumber,
		})
		ids = append(ids, strconv.FormatUint(r.ID, 10))
	}
	sort.Sort(recordList)
	useRecordList := make(util.SortRecordList, 0, len(luckyNumberRecords))
	for _, r := range recordList {
		diff := int64(r.Num) - int64(number)
		if diff < 0 {
			r.Tag = "-"
		} else {
			r.Tag = "+"
		}
		r.Num = uint64(util.Abs(diff))
		useRecordList = append(useRecordList, r)
	}
	sort.Sort(useRecordList)
	sql := "UPDATE lucky_number_record SET ranking = CASE id "
	for k, r := range useRecordList {
		sql += fmt.Sprintf("WHEN %d THEN %d ", r.ID, k+1)
	}
	sql += fmt.Sprintf("END WHERE id IN (%s)", strings.Join(ids, ","))
	if err := tx.Raw(sql).Error; err != nil {
		tx.Rollback()
		return &activity, err
	}

	tx.Commit()
	return &activity, nil
}

// DelActivity 删除活动
func (d *dao) DelActivity(ctx context.Context, activityID uint64) error {
	return checkErr(d.db.Delete(model.Activity{}, "id = ?", activityID).Error)
}

// CurrentlyAvailableActivity 当前可参与的活动
func (d *dao) CurrentlyAvailableActivity(ctx context.Context) (*model.Activity, error) {
	var activity model.Activity
	err := checkErr(d.db.Where("start_time <= ? and end_time >= ?", time.Now(), time.Now()).First(&activity).Error)
	return &activity, err
}

// ActivityAllPrizeIssued 活动所有已发放的奖品
func (d *dao) ActivityAllPrizeIssued(ctx context.Context) (int, error) {
	total := 0
	var activitys []*model.Activity
	if err := checkErr(d.db.Where("prize_number != ?", 0).Find(&activitys).Error); err != nil {
		return total, err
	}
	if len(activitys) == 0 {
		return 0, nil
	}
	querySlice := make([]string, 0, len(activitys))
	for _, v := range activitys {
		querySlice = append(querySlice, fmt.Sprintf(" (activity_id = %d and ranking between %d and %d) ", v.ID, 1, v.PrizeAmount))
	}
	query := strings.Join(querySlice, "or")
	err := checkErr(d.db.Model(&model.LuckyNumberRecord{}).Where(query).Count(&total).Error)
	return total, err
}
