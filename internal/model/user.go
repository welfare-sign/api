package model

// User 后台管理员
type User struct {
	Base

	Name     string `json:"name" gorm:"not null"`     // 用户名
	Password string `json:"password" gorm:"not null"` // 密码
}

// UserVO 用户登录
type UserVO struct {
	Name     string `json:"name" validate:"required"`     // 用户名
	Password string `json:"password" validate:"required"` // 密码
}
