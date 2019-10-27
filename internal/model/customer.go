package model

import "time"

// Customer 顾客
type Customer struct {
	Base

	Name            string    `json:"name" gorm:"not null"`                    // 称呼
	Mobile          string    `json:"mobile" gorm:"type:varchar(50);not null"` // 手机号
	LastCheckinTime time.Time `json:"last_checkin_time" gorm:"type:datetime"`  // 最后一次签到时间
}

// CustomerListVO 查询顾客列表参数
type CustomerListVO struct {
	Name     string `form:"name" json:"name"`
	Mobile   string `form:"mobile" json:"mobile"`
	PageNo   int    `form:"page_no" json:"page_no"`
	PageSize int    `form:"page_size" json:"page_size"`
}
