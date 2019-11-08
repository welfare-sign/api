package dao

import (
	"context"

	"welfare-sign/internal/dao/mysql"
	"welfare-sign/internal/model"
)

const nearMerchantSQL = `SELECT
id, (
	6371 * acos (
	cos ( radians(?) )
	* cos( radians( lat ) )
	* cos( radians( lon ) - radians(?) )
	+ sin ( radians(?) )
	* sin( radians( lat ) )
  )
) AS distance
FROM merchant
HAVING distance <= ?
ORDER BY distance ASC
LIMIT ?;`

// CreateMerchant create merchant
func (d *dao) CreateMerchant(ctx context.Context, data model.Merchant) error {
	data.SetDefaultAttr()
	return d.db.Create(&data).Error
}

// ListMerchant get merchant list
// pageNo >= 1
func (d *dao) ListMerchant(ctx context.Context, query interface{}, pageNo, pageSize int) ([]*model.Merchant, int, error) {
	var merchants []*model.Merchant
	total := 0
	err := d.db.Where(query).Limit(pageSize).Offset((pageNo - 1) * pageSize).Order("created_at desc").Find(&merchants).Error
	if mysql.IsError(err) {
		return merchants, total, err
	}
	if err := d.db.Where(query).Find(&model.Merchant{}).Count(&total).Error; mysql.IsError(err) {
		return merchants, total, err
	}
	return merchants, total, nil
}

// FindMerchant 获取商家详情
func (d *dao) FindMerchant(ctx context.Context, query interface{}) (*model.Merchant, error) {
	var merchant model.Merchant
	err := checkErr(d.db.Where(query).First(&merchant).Error)
	return &merchant, err
}

// EcecWriteOff 执行核销
func (d *dao) EcecWriteOff(ctx context.Context, merchantID, customerID, hasRece, totalRece uint64) error {
	tx := d.db.Begin()
	if err := tx.Model(&model.IssueRecord{}).Where(map[string]interface{}{
		"merchant_id": merchantID,
		"customer_id": customerID,
	}).Updates(map[string]interface{}{"received": hasRece}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Merchant{}).Where("id = ?", merchantID).Updates(map[string]interface{}{"received": totalRece}).Error; err != nil {
		tx.Rollback()
		return nil
	}
	tx.Commit()
	return nil
}

type nears struct {
	ID       uint64  `json:"id"`
	Distance float64 `json:"distance"`
}

// NearMerchant 附近的商家
func (d *dao) NearMerchant(ctx context.Context, data *model.NearMerchantVO) ([]*model.Merchant, error) {
	var (
		merchants []*model.Merchant
		ns        []*nears
	)

	err := d.db.Raw(nearMerchantSQL, data.Lat, data.Lon, data.Lat, data.Distince, data.Num).Find(&ns).Error
	if checkErr(err) != nil {
		return nil, err
	}
	if len(ns) == 0 {
		return merchants, nil
	}
	ids := make([]uint64, 0, len(ns))
	for i := 0; i < len(ns); i++ {
		ids = append(ids, ns[i].ID)
	}
	err = d.db.Where("id in (?)", ids).Find(&merchants).Error
	if checkErr(err) != nil {
		return nil, err
	}
	return merchants, nil
}
