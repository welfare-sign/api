package dao

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
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

// GetCompositeIndex 获取每周五的上证指数
func (d *dao) GetCompositeIndex(ctx context.Context) (*model.CompositeIndex, error) {
	var compositeIndex model.CompositeIndex
	err := checkErr(d.db.Raw(getCompositeIndex, global.ActiveStatus).First(&compositeIndex).Error)
	if err != nil || compositeIndex.ID == 0 {
		return &compositeIndex, err
	}
	return &compositeIndex, err
}

// GetCompositeIndexByQuery 根据参数获取上证指数
func (d *dao) GetCompositeIndexByQuery(ctx context.Context, query interface{}) (*model.CompositeIndex, error) {
	var compositeIndex model.CompositeIndex
	err := checkErr(d.db.Where(query).Find(&compositeIndex).Error)
	return &compositeIndex, err
}

// GetCompositeIndexBefore 获取上周五的上证指数
func (d *dao) GetCompositeIndexBefore(ctx context.Context) (*model.CompositeIndex, error) {
	var compositeIndex model.CompositeIndex
	err := checkErr(d.db.Raw(getCompositeIndexBefore, global.ActiveStatus).First(&compositeIndex).Error)
	if err != nil || compositeIndex.ID == 0 {
		return &compositeIndex, err
	}
	return &compositeIndex, err
}

// StoreCompositeIndex 存储或者更新上证指数
func (d *dao) StoreCompositeIndex(ctx context.Context, compositeDate string, points float64) error {
	tx := d.db.Begin()

	var rightDate int
	tx.Raw(isRightDate, compositeDate).Pluck("date", &rightDate)
	if rightDate == 0 {
		tx.Rollback()
		return errors.New("填写的期不是本周五")
	}
	var compositeIndex model.CompositeIndex
	if err := checkErr(tx.Where(map[string]interface{}{
		"composite_date": compositeDate,
		"status":         global.ActiveStatus,
	}).First(&compositeIndex).Error); err != nil {
		tx.Rollback()
		return err
	}
	if compositeIndex.ID != 0 { // 更新
		compositeIndex.Points = points
		if err := tx.Save(compositeIndex).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else { // 添加
		compositeIndex.SetDefaultAttr()
		compositeIndex.CompositeDate = compositeDate
		compositeIndex.Points = points
		if err := tx.Create(&compositeIndex).Error; err != nil {
			tx.Rollback()
			return err
		}
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
	p, err := strconv.ParseInt(strconv.FormatFloat(points*100, 'f', 2, 64)[:6], 10, 64)
	if err != nil {
		tx.Rollback()
		return err
	}
	numberRecords, err := d.GetRoundLuckyNumberRecordSQL(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	recordList := make(util.SortRecordList, 0, len(numberRecords))
	ids := make([]string, 0, len(numberRecords))
	for _, r := range numberRecords {
		recordList = append(recordList, util.SortRecord{
			ID:  r.ID,
			Num: r.LuckyNumber,
		})
		ids = append(ids, strconv.FormatUint(r.ID, 10))
	}
	sort.Sort(recordList)
	useRecordList := make(util.SortRecordList, 0, len(numberRecords))
	for _, r := range recordList {
		r.Num = r.Num - p
		if r.Num < 0 {
			r.Tag = "-"
		} else {
			r.Tag = "+"
		}
		r.Num = util.Abs(r.Num)
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
		return err
	}

	tx.Commit()
	return nil
}
