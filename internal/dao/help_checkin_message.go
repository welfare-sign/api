package dao

import (
	"context"

	"welfare-sign/internal/global"
	"welfare-sign/internal/model"
)

// FindHelpCheckinMesage 查找补签消息
func (d *dao) FindHelpCheckinMesage(ctx context.Context, query interface{}, args ...interface{}) (*model.HelpCheckinMessage, error) {
	var helpCheckinMessage model.HelpCheckinMessage
	err := checkErr(d.db.Where(query, args...).First(&helpCheckinMessage).Error)
	return &helpCheckinMessage, err
}

// UpdateHelpCheckinMessage 更新补签消息
func (d *dao) UpdateHelpCheckinMessage(ctx context.Context, customerID uint64) error {
	if err := d.db.Model(&model.HelpCheckinMessage{}).Where(map[string]interface{}{
		"customer_id": customerID,
		"is_read":     global.UnRead,
	}).Update("is_read", global.Readed).Error; err != nil {
		return err
	}
	return nil
}
