package dao

import (
	"context"

	"welfare-sign/internal/model"
)

// FindUser find user
func (d *dao) FindUser(ctx context.Context, query interface{}) (*model.User, error) {
	var user model.User
	err := d.db.Where(query).Find(&user).Error
	return &user, err
}
