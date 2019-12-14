package dao

import (
	"context"

	"welfare-sign/internal/model"
)

const (
	getRegisterStatSQL = `
SELECT DATE_FORMAT(created_at, "%Y-%m-%d") AS date, count(id) as num FROM customer
WHERE DATE_FORMAT(created_at, "%Y-%m-%d") BETWEEN ? AND ?
GROUP BY DATE_FORMAT(created_at, "%Y-%m-%d")
`
	getCheckinStat = `
	SELECT DATE_FORMAT(created_at, "%Y-%m-%d") AS date, count(id) as num FROM checkin_record_log
WHERE DATE_FORMAT(created_at, "%Y-%m-%d") BETWEEN ? AND ?
GROUP BY DATE_FORMAT(created_at, "%Y-%m-%d")
	`
)

// GetRegisterStat 统计用户注册数目
func (d *dao) GetRegisterStat(ctx context.Context, beginDate, endDate string) ([]*model.RegisterStat, error) {
	var stat []*model.RegisterStat
	err := checkErr(d.db.Raw(getRegisterStatSQL, beginDate, endDate).Find(&stat).Error)
	return stat, err
}

// GetCheckinStat 统计用户签到数目
func (d *dao) GetCheckinStat(ctx context.Context, beginDate, endDate string) ([]*model.CheckinStat, error) {
	var stat []*model.CheckinStat
	err := checkErr(d.db.Raw(getCheckinStat, beginDate, endDate).Find(&stat).Error)
	return stat, err
}
