package dao

import (
	"context"

	"welfare-sign/internal/model"
)

// FindIssueRecord 获取礼包发放记录详情
func (d *dao) FindIssueRecord(ctx context.Context, query interface{}) (*model.IssueRecord, error) {
	var issueRecord model.IssueRecord
	err := d.db.Where(query).First(&issueRecord).Error
	return &issueRecord, err
}
