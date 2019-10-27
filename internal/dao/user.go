package dao

import (
	"context"

	"welfare-sign/internal/model"
)

// FindUser find user
func (d *dao) FindUser(ctx context.Context, data model.User) (*model.User, error) {
	var user model.User
	err := d.db.Where(&data).Find(&user).Error
	return &user, err
}
